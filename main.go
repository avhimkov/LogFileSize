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
	conf()
	dir := viper.GetString("dir.dirFile")

	runHTTP(dir)
}

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

func ShowStat(w http.ResponseWriter, r *http.Request) {
	//temp := viper.GetString("temp.temp")
	//Okno1 := viper.GetString("windows.okno1")

	render(w, "header.html")
	//table(w, r, temp)
	htmlRang(w, r)
	render(w, "footer.html")
	//os.RemoveAll("file/temp/file")
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

func removeFile(target string) {
	err := os.Remove(target)
	if err != nil {
		log.Fatal(err)
	}
}

type config struct {
	Name string
	PathMap string `mapstructure:"windows"`
}

//var window = map[string]interface{}{
//	"Окно№1": "file/Окно№1/",
//	"Окно№2": "file/Окно№2/",
//	"Окно№3": "file/Окно№3/",
//	"Окно№4": "file/Окно№4/",
//	"Окно№5": "file/Окно№5/",
//	"Окно№6": "file/Окно№6/",
//	"Окно№7": "file/Окно№7/",
//}

func htmlRang(w http.ResponseWriter, r *http.Request)  {

	window:=viper.GetStringMap("windows")

	html := `<select> // for loop in html template example
		   {{range $key, $value := .}}
		     <option value="{{ $value }}">{{ $key }}</option>
		   {{end}}
		 </select>`

	selectTemplate, err := template.New("select").Parse(string(html))
	if err != nil {
		panic(err)
	}

	// populate dropdown with fruits
	selectTemplate.Execute(w, window)

	// no need for this...
	//w.Write([]byte(html))
}

func dataSend(w http.ResponseWriter, r *http.Request)  (string, string) {

	r.FormValue("name")
	r.ParseForm()
	date1 := r.Form.Get("calendar")
	okno1 := r.Form.Get("okno")

	fmt.Fprint(w, "<form action=\"\" method=\"get\">" +
		"<p><input type=\"date\" name=\"calendar\"/>" +
		"<select id=\"okno\" name=\"okno\">" +
		"<option enable>Выбрать окно</option>" +
		"<option value=\"Окно№1\">Окно№1</option>" +
		"<option value=\"Окно№2\">Окно№2</option>" +
		"<option value=\"Окно№3\">Окно№3</option>" +
		"<option value=\"Окно№4\">Окно№4</option>" +
		"<option value=\"Окно№5\">Окно№5</option>" +
		"<option value=\"Окно№6\">Окно№6</option>" +
		"<option value=\"Окно№7\">Окно№7</option>" +
		"</select>" +
		"<p><input type=\"submit\" value=\"Показать\"></p></p></form>")
	//fmt.Println(r.Form["calendar"])
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	return date1, okno1
}

func table(w http.ResponseWriter, r *http.Request, dirTemp string) {

	data, okno := dataSend(w, r)
	//fmt.Printf("okno: %v \n", okno)
	oknoS:= "file/" + okno + "/"

	listDirZip := listfiles(oknoS, ".zip", data) //2017-03-29

	fmt.Fprintf(w, "<tr><td colspan=\"5\" align=\"center\" style=\"width: 500px;\">%s</td></tr>", oknoS)//dirZip

	for i := range listDirZip {
		unzip(listDirZip[i], dirTemp+listDirZip[i])

		daysAgo := daysAgo(listDirZip[i], day)
		dcreat := dataCreate(listDirZip[i])
		dcreatf := dcreat.Format("2006-01-02")
		fmt.Printf("dcreat: %v \n", dcreatf)
		dir := listDirZip[i]
		size := sizeFile(listDirZip[i])

		listDirTemp := listfilescler(dirTemp, ".wav")
		dirtemp := listDirTemp[i]

		fmt.Printf("dirtemp1: %v \n", dirtemp)

		fmt.Fprintf(w, "<tr><td align=\"left\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%.2f дней</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">" +
			"<form action=\"%s\"><input type=\"submit\" value=\"Прослушать\"/></form>" +
			"</td></tr>",

			//"<td align=\"center\" style=\"width: 100px;\"><audio controls><source src=%s type=\"audio/wav\"></audio></td></tr>",
			dir, dcreatf, daysAgo, size, dirtemp)
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

func sizeFile(path string) string {

	stat, err := os.Stat(path)
	sizeStr, err := convertSize(stat.Size())
	if err != nil {
		fmt.Println(err)
	}
	return sizeStr
}

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
		//fmt.Printf("list3: %v \n", list)
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return list
}

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
