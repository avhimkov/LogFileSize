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
	"strings"
	"io/ioutil"
)

func main() {
	//load config
	conf()
	//clear temp folder
	var temp = viper.GetString("dir.temp")
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
func audioListen(w http.ResponseWriter, r *http.Request) {
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
	checkErr(err)
}

//render table, calendar and button for table monit
func head(w http.ResponseWriter, r *http.Request) string {
	r.FormValue("name")
	r.ParseForm()
	date := r.Form.Get("calendar")
	t, err := template.ParseFiles("templates/headtpl.html")
	checkErr(err)
	t.Execute(w, date)
	return date
}

//render table, calendar and button for table audio
func htmlRang(w http.ResponseWriter, r *http.Request) (string, string, string) {

	window := viper.GetStringMap("windows")
	timeform := viper.GetStringMap("time")

	r.FormValue("name")
	r.ParseForm()

	date := r.Form.Get("calendar")
	windowform := r.Form.Get("window")
	timemodif := r.Form.Get("time")
	t, err := template.ParseFiles("templates/range.html")
	checkErr(err)
	//if err != nil {
	//	fmt.Fprint(w, "<p>page not found 404</p>")
	//	panic(err)
	//}
	t.ExecuteTemplate(w, "window", window)
	t.ExecuteTemplate(w, "time", timeform)
	return date, windowform, timemodif
}

//render tableaudio
func tableAudio(w http.ResponseWriter, r *http.Request) {
	typefiles := viper.GetString("filetype.archivefile")

	dir := viper.GetString("dir.works")
	temp := viper.GetString("dir.temp")
	date, windowform, timemodif := htmlRang(w, r)

	windowS := dir + windowform
	fmt.Fprint(w, "<tr class=\"warning\"><td colspan=\"6\" >" + windowform+"  Время: " + timemodif + "</td></tr>")

	listDirArchive, _ , _:= listFiles(windowS, typefiles, date, timemodif) //2017-03-29

	if typefiles == ".zip" {
		for j := range listDirArchive{
			UnZip(listDirArchive[j], temp)
			//UnZip(listDirArchive[j], temp + listDirArchive[j])
		}

		for i := range listDirArchive {
			audiofile := viper.GetString("filetype.audiofile")
			dir := listDirArchive[i]

			dcreat, _ := dateCreate(listDirArchive[i])
			dcreatf := dcreat.Format("2006-01-02")

			_, daysAgo := dateCreate(listDirArchive[i])
			//daysAgo := daysAgo(listDirArchive[i], day)

			dhoursAgo, _ := dateCreate(listDirArchive[i])
			dhoursAgof := dhoursAgo.Hour()

			size, _ := sizeFile(listDirArchive[i])

			_, _, listDirTemp := listFiles(temp, audiofile, "","")
			dirtemp := listDirTemp[i]

			fmt.Fprintf(w, "<tr>" +
				"<td align=\"left\" \">%s</td>" +
				"<td align=\"center\" >%s</td>" +
				"<td align=\"center\" >%.2f дней</td>" +
				"<td align=\"center\" >%d часов</td>" +
				"<td align=\"center\">%s</td>" +
				"<td align=\"center\">" +
				"<form action=\"%s\"><input type=\"submit\" class=\"btn btn-primary\" value=\"Прослушать\"/></form>" +
				"</td></tr>",
				//"<td align=\"center\" style=\"width: 100px;\"><audio controls><source src=%s type=\"audio/wav\"></audio></td></tr>",
				dir, dcreatf, daysAgo, dhoursAgof, size, dirtemp)
		}
	} else {
		for j := range listDirArchive{
			copyFile(listDirArchive[j], temp + listDirArchive[j])
		}
		for i := range listDirArchive {
			audiofile := viper.GetString("filetype.audiofile")

			dir := listDirArchive[i]

			dcreat, _ := dateCreate(listDirArchive[i])
			dcreatf := dcreat.Format("2006-01-02")

			_, daysAgo := dateCreate(listDirArchive[i])
			//daysAgo := daysAgo(listDirArchive[i], day)

			dhoursAgo, _ := dateCreate(listDirArchive[i])
			dhoursAgof := dhoursAgo.Hour()

			size, _ := sizeFile(listDirArchive[i])

			_, _, listDirTemp := listFiles(temp, audiofile, "","")
			dirtemp := listDirTemp[i]

			fmt.Fprintf(w, "<tr>" +
				"<td align=\"left\" \">%s</td>" +
				"<td align=\"center\" >%s</td>" +
				"<td align=\"center\" >%.2f дней</td>" +
				"<td align=\"center\" >%d часов</td>" +
				"<td align=\"center\">%s</td>" +
				"<td align=\"center\">" +
				"<form action=\"%s\"><input type=\"submit\" class=\"btn btn-primary\" value=\"Прослушать\"/></form>" +
				"</td></tr>",
				//"<td align=\"center\" style=\"width: 100px;\"><audio controls><source src=%s type=\"audio/wav\"></audio></td></tr>",
				dir, dcreatf, daysAgo, dhoursAgof, size, dirtemp)
		}
	}
}

//render table montin
func tableMonitoring(w http.ResponseWriter, r *http.Request) {
	date := head(w, r)
	archive := viper.GetString("filetype.archivefile")
	dir := viper.GetString("dir.works")

	_, listDirArchive, _ := listFiles(dir, archive, date, "") //2017-03-29
	//listDirArchive := listFilesDate(dir, archive, date) //2017-03-29

	for i := range listDirArchive {
		smallfile := viper.GetInt64("size.file")
		_, daysAgo := dateCreate(listDirArchive[i])
		//daysAgo := daysAgo(listDirArchive[i], day)
		dcreat, _ := dateCreate(listDirArchive[i])
		dcreatf := dcreat.Format("2006-01-02")
		dir := listDirArchive[i]
		size, _ := sizeFile(listDirArchive[i])
		_, sizeint := sizeFile(listDirArchive[i])
		if sizeint > smallfile {
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

func listFiles(rootpath string, typefile string, data string, time string) ([]string, []string, []string) {

	list := make([]string, 0, 10)
	list1 := make([]string, 0, 10)
	list2 := make([]string, 0, 10)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		modification := info.ModTime().Format("2006-01-02")
		timetempleat := info.ModTime().Format("3")

		if info.IsDir() {
			return nil
		}
		if modification == data {
			if strings.EqualFold(timetempleat,  time) {
				if filepath.Ext(path) == typefile {
					list = append(list, path)
				}
			}
		}
		if modification == data {
			if filepath.Ext(path) == typefile {
				list1 = append(list1, path)
			}
		}

		if filepath.Ext(path) == typefile {
			list2 = append(list2, path)
		}
		return nil
	})
	checkErr(err)
	return list, list1, list2
}

//func return list files in dir appropriate type file
//func listFilesClear(rootpath string, typefile string) []string {
//	list := make([]string, 0, 10)
//	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
//		if info.IsDir() {
//			return nil
//		}
//		if filepath.Ext(path) == typefile {
//			list = append(list, path)
//		}
//		return nil
//	})
//	checkErr(err)
//	return list
//}

//func http server
func runHTTP() {
	dirServer := viper.GetString("dir.server")
	//dirWorks := viper.GetString("dir.works")
	http.HandleFunc("/", ShowStat)
	http.HandleFunc("/audio", audioListen)
	log.Println("http://localhost:8080 Listening...")

	http.HandleFunc(dirServer, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	//http.HandleFunc(dirWorks, func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, r.URL.Path[1:])
	//})

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
func dateCreate(path string) (time.Time, float64) {
	now := time.Now()

	file, err := os.Stat(path)
	checkErr(err)

	modifiedtime := file.ModTime()
	checkErr(err)

	diff := now.Sub(modifiedtime)
	days := float64(diff.Hours() / 24)

	return modifiedtime, days
}

//func return file size string and int64 format
func sizeFile(path string) (string, int64) {
	stat, err := os.Stat(path)
	sizeStr, err := convertSize(stat.Size())
	sizeInt64 := stat.Size()
	checkErr(err)
	return sizeStr, sizeInt64
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//func return list files in dir appropriate type file and date create
func copyFile(src string, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	checkErr(err)
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	checkErr(err)
}

//UnZip file
func UnZip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	checkErr(err)
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
		checkErr(err)
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		checkErr(err)
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}