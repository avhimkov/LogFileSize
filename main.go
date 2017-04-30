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

var day = time.Now()

func main() {
	//load config
	conf()
	//clear temp folder
	temp := viper.GetString("dir.temp")
	os.RemoveAll(temp)
	//run http server
	runHTTP()
}

//config JSON file init
func conf() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}
	viper.SetDefault("msg", "Hello World (default)")
}

//render page audio lessen
func AudioListen(w http.ResponseWriter, r *http.Request) {
	render(w, "header.html")
	tableAudio(w, r)
	render(w, "footer.html")
}

// render page stat all file in all folder
func ShowStat(w http.ResponseWriter, r *http.Request) {
	render(w, "header.html")
	tableMonitoring(w, r)
	render(w, "footer.html")
}

//func render template
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

//render table, calendar and button for table monit
func head(w http.ResponseWriter, r *http.Request) string {
	r.FormValue("name")
	r.ParseForm()
	date := r.Form.Get("calendar")
	t, err := template.ParseFiles("templates/headtpl.html")
	if err != nil {
		fmt.Fprint(w, "<p>page not found 404</p>")
		panic(err)
	}
	t.Execute(w, date)
	return date
}

//render table, calendar and button for table audio
func htmlRang(w http.ResponseWriter, r *http.Request) (string, string)  {

	window :=viper.GetStringMap("windows")

	r.FormValue("name")
	r.ParseForm()

	date := r.Form.Get("calendar")
	okno := r.Form.Get("okno")
	t, err := template.ParseFiles("templates/range.html")
	if err != nil {
		fmt.Fprint(w, "<p>page not found 404</p>")
		panic(err)
	}
	t.Execute(w, window)
	return date, okno
}

//render tableaudio
func tableAudio(w http.ResponseWriter, r *http.Request) {

	dir := viper.GetString("dir.AllFiles")
	Temp := viper.GetString("dir.temp")
	fmt.Printf("tableAudio: %v \n", dir)
	data, okno := htmlRang(w, r)
	oknoS := dir + okno
	fmt.Printf("oknoS: %v \n", oknoS)
	fmt.Fprint(w, "<tr class=\"warning\"><td colspan=\"5\" >"+ okno +"</td></tr>")

	listDirZip := listfiles(oknoS, ".zip", data) //2017-03-29
	fmt.Printf("listDirZip: %v \n", listDirZip)

	for i := range listDirZip {
		unzip(listDirZip[i], Temp + listDirZip[i])

		daysAgo := daysAgo(listDirZip[i], day)
		dcreat := dataCreate(listDirZip[i])
		dcreatf := dcreat.Format("2006-01-02")
		dir := listDirZip[i]
		size := sizeFile(listDirZip[i])

		listDirTemp := listfilescler(Temp, ".wav")
		dirtemp := listDirTemp[i]

		fmt.Printf("dirtemp: %v \n", dirtemp)
		fmt.Fprintf(w, "<tr>" +
			"<td align=\"left\" \">%s</td>" +
			"<td align=\"center\" >%s</td>" +
			"<td align=\"center\" >%.2f дней</td>" +
			"<td align=\"center\">%s</td>" +
			"<td align=\"center\" >" +
			"<form action=\"%s\"><input type=\"submit\" class=\"btn btn-primary\" value=\"Прослушать\"/></form>" +
			"</td></tr>",
			//"<td align=\"center\" style=\"width: 100px;\"><audio controls><source src=%s type=\"audio/wav\"></audio></td></tr>",
			dir, dcreatf, daysAgo, size, dirtemp)
	}
}

//render table motin
func tableMonitoring(w http.ResponseWriter, r *http.Request) {
	date := head(w, r)
	dir := viper.GetString("dir.AllFiles")
	fmt.Printf("dirAllMonit: %v \n", dir)
	listDirZip := listfiles(dir, ".zip", date) //2017-03-29

	for i := range listDirZip {
		smallfile := viper.GetInt64("size.file")
		daysAgo := daysAgo(listDirZip[i], day)
		dcreat := dataCreate(listDirZip[i])
		dcreatf := dcreat.Format("2006-01-02")
		dir := listDirZip[i]
		size := sizeFile(listDirZip[i])
		sizeint := sizeFileint(listDirZip[i])
		fmt.Printf("sizeint: %v \n", sizeint)
		if sizeint > smallfile{
			fmt.Fprintf(w, "<tr>" +
				"<td align=\"left\" \">%s</td>" +
				"<td align=\"center\" >%s</td>" +
				"<td align=\"center\" >%.2f дней</td>" +
				"<td align=\"center\">%s</td>" +
				"</tr>",
				dir, dcreatf, daysAgo, size)
		} else {
			fmt.Fprintf(w, "<tr>" +
				"<td bgcolor=\"#ffcc00\" align=\"left\" \">%s</td>" +
				"<td bgcolor=\"#ffcc00\" align=\"center\" >%s</td>" +
				"<td bgcolor=\"#ffcc00\" align=\"center\" >%.2f дней</td>" +
				"<td bgcolor=\"#ffcc00\" align=\"center\">%s</td>" +
				"</tr>",
				dir, dcreatf, daysAgo, size)
		}
	}
}

//func http server
func runHTTP() {
	dir := viper.GetString("dir.Server")
	http.HandleFunc("/", ShowStat)
	http.HandleFunc("/audio", AudioListen)
	log.Println("http://localhost:8080 Listening...")

	http.HandleFunc(dir, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/bootstrap/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	http.ListenAndServe(":8080", nil)
}

//convert size file from int64 to string and convert in humans format
func convertSize(size int64) (string, error) {
	if size == 0 {
		return "0B", nil
	}
	sizeName := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	i := int(math.Log2(float64(size)) / 10)
	humanSize := fmt.Sprintf("%d %s", size/int64(math.Pow(1024, float64(i))), sizeName[i])
	return humanSize, nil
}

//func return time date create files
func dataCreate(path string) time.Time {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}
	modifiedtime := file.ModTime()
	if err != nil {
		fmt.Println(err)
	}
	return modifiedtime
}

//func calculate dais ago
func daysAgo(path string, now time.Time) float64 {
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
	days := float64(diff.Hours() / 24)
	return days
}

//func return file size string format
func sizeFile(path string) string {

	stat, err := os.Stat(path)
	sizeStr, err := convertSize(stat.Size())
	if err != nil {
		fmt.Println(err)
	}
	return sizeStr
}

//func return file size int64 format
func sizeFileint(path string) int64 {

	stat, err := os.Stat(path)
	sizeStr := stat.Size()
	if err != nil {
		fmt.Println(err)
	}
	return sizeStr
}

//func return list files in dir appropriate type file and date create
func listfiles(rootpath string, typefile string, data string) []string {

	list := make([]string, 0, 10)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		modification:=info.ModTime().UTC().Format("2006-01-02")
		if info.IsDir() {
			return nil
		}
		if modification == data {
			if filepath.Ext(path) == typefile {
				list = append(list, path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return list
}

//func return list files in dir appropriate type file
func listfilescler (rootpath string, typefile string) []string {

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

//unzip file
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