package main

import (
	"fmt"
	"math"
	"net/http"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"io"
	"time"
	"archive/zip"
	"github.com/spf13/viper"
)

var day = time.Now().Local()

//type data struct {
//	Data string
//}

func main() {
	runHTTP("/file/")
}

func conf()  {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}
	//viper.SetDefault("msg", "Hello World (default)")
}

func ShowStat(w http.ResponseWriter, r *http.Request) {
	conf()
	temp := viper.GetString("temp.temp")
	Okno1 := viper.GetString("windows.okno1")
	Okno2 := viper.GetString("windows.okno2")

	render(w, "header.html")
	//tableOkno(w, Okno1, temp)
	//tableOkno(w, Okno2, temp)

	table(w, Okno1, temp)
	table(w, Okno2, temp)

	render(w, "footer.html")
	os.RemoveAll("file/temp/file")
}

func render(w http.ResponseWriter, tmpl string) {
	tmpl = fmt.Sprintf("templates/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

//func tableOkno(w http.ResponseWriter, okno string, temp string)  {
//	str := fmt.Sprintf("<tr><td colspan=\"5\" align=\"center\" style=\"width: 500px;\">%s</td></tr>", okno)
//	io.WriteString(w, str)
//	table(w, okno, temp)
//}

func removeFile(target string)  {
	err := os.Remove(target)
	if err != nil {
		log.Fatal(err)
	}
}

func table(w http.ResponseWriter, dirZip string, dirTemp string) {
	listDirZip := listfiles(dirZip, ".zip")
	//listDirTemp := listfiles(dirTemp, ".wav")
	listDir := dirlist(dirZip)
	fmt.Fprintf(w, "<tr><td colspan=\"5\" align=\"center\" style=\"width: 500px;\">%s</td></tr>", dirZip)


	for i := range listDirZip {
		unzip(listDirZip[i], dirTemp + listDirZip[i])
	}

	for i := range listDir {
		fmt.Fprintf(w,"<tr><td align=\"left\" style=\"width: 100px;\">%s</td></tr>", listDir[i])
	}

	for i := range listDirZip {

		daysAgo := daysAgo(listDirZip[i], day)
		dcreat := dataCreate(listDirZip[i])
		dir := listDirZip[i]
		size := sizeFile(listDirZip[i])


		listDirTemp := listfiles(dirTemp, ".wav")
		dirtemp := listDirTemp[i]

		fmt.Fprintf(w, "<tr><td align=\"left\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%d дней</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%s</td>" +
			//"<td align=\"center\" style=\"width: 100px;\"><button class=\"play\">play</button></td>"+
			"<td align=\"center\" style=\"width: 100px;\"><audio controls><source src=%s type=\"audio/wav\"></audio></td></tr>",
		 	dir, dcreat, daysAgo, size, dirtemp)
	}
}

func runHTTP(dir string) {
	http.HandleFunc("/", ShowStat)
	log.Println("localhost:8080 Listening...")
	http.HandleFunc(dir, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.ListenAndServe(":8080", nil)
}

func convertSize(size int64) (string, error) {
	if size == 0 {
		return "0B", nil
	}
	sizeName := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	i := int(math.Log2(float64(size)) / 10)
	humanSize := fmt.Sprintf("%d%s", size/int64(math.Pow(1024, float64(i))), sizeName[i])
	return humanSize, nil
}
func dataCreate(path string) string {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}
	modifiedtime := file.ModTime()
	if err != nil {
		fmt.Println(err)
	}
	modifiedtimef := modifiedtime.Format("2006-01-02 15:04:05")

	return modifiedtimef
}

func daysAgo(path string, now time.Time) int {
	dataCreate(path)
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}
	modifiedtime := file.ModTime()
	if err != nil {
		fmt.Println(err)
	}
	diff := now.Sub(modifiedtime)
	days := int(diff.Hours() / 24)
	return days
}

func sizeFile(path string) string {

	stat, err := os.Stat(path)
	sizeStr, err := convertSize(stat.Size())
	if err != nil {
		fmt.Println(err)
	}
	return sizeStr
}

func dirlist(rootpath string) []string {

	list := make([]string, 0, 10)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			list = append(list, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return list
}

func listfiles(rootpath string, typefile string) []string {

	list := make([]string, 0, 10)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == typefile {
			list = append(list, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return list
}

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}