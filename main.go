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
	"strings"
	"archive/zip"
)

var day = time.Now().Local()
var Okno1 = "file/Окно№1/"
var Okno2 = "file/Окно№2/"
var temp = "file/temp/"
//dirOkno3:= "file/dir3"
//dirOkno4:= "file/dir4"
//dirOkno5:= "file/dir5"

var templ = "<td align=\"center\" style=\"width: 100px;\"><audio controls>" +
	"<source src= " + "" + "type=\"audio/wav\"></audio></td>" +
	"</tr>"

//type vars struct{
//	Dir string
//}

func main() {
	runHTTP("/file/")
	//unzip(Okno1, "/file/temp")
}

func ShowStat(w http.ResponseWriter, r *http.Request) {

	render(w, "header.html")
	tableOkno(w, Okno1, temp)
	//tabletemp(w, temp)
	tableOkno(w, Okno2, temp)
	//tabletemp(w, temp)
	render(w, "footer.html")
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

func tableOkno(w http.ResponseWriter, okno string, temp string)  {
	str := fmt.Sprintf("<tr><td colspan=\"5\" align=\"center\" style=\"width: 500px;\">%s</td></tr>", okno)
	io.WriteString(w, str)
	table(w, okno, temp)
}

func table(w http.ResponseWriter, dirZip string, dirTemp string) {
	listDirZip := listfiles(dirZip, ".zip")
	//listDirTemp := listfiles(dirTemp, ".zip")

	for i := range listDirZip {
		unzip(listDirZip[i], dirTemp + listDirZip[i])
	}

	for i := range listDirZip {
		daysAgo := daysAgo(listDirZip[i], day)
		dcreat := dataCreate(listDirZip[i])
		dir := listDirZip[i]
		size := sizeFile(listDirZip[i])
		str:=fmt.Sprintf("<tr><td align=\"left\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%d дней</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%s</td>" + templ + "", dir, dcreat, daysAgo, size)
		io.WriteString(w, str)
	}
}

//func tabletemp(w http.ResponseWriter, dirTemp string)  {
//	listDirTemp := listfiles(dirTemp,".wav")
//	for i := range listDirTemp{
//		dir := listDirTemp[i]
//		str:=fmt.Sprintf("<td align=\"center\" style=\"width: 100px;\"><audio controls>" +
//			"<source src=%s type=\"audio/wav\"></audio></td>" +
//			"</tr>", dir)
//		io.WriteString(w, str)
//	}
//}

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

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath,string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}