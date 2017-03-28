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

var Okno1 = "file/Окно №1/"
var day = time.Now().Local()
//dirOkno3:= "file/dir3"
//dirOkno4:= "file/dir4"
//dirOkno5:= "file/dir5"

type vars struct{
	Dir string
}

func main() {
	runHTTP("/file/")
}

func ShowStat(w http.ResponseWriter, r *http.Request) {

	render(w, "header.html")
	tableOkno(w, Okno1)
	//tableOkno(w, Okno2)
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

func tableOkno(w http.ResponseWriter, okno string)  {
	str := fmt.Sprintf("<tr><td colspan=\"5\" align=\"center\" style=\"width: 500px;\">%s</td></tr>", okno)
	io.WriteString(w, str)
	table(w, okno)
}

func table(w http.ResponseWriter, dir string) {
	listDir1 := listfiles(dir)

	for i := range listDir1 {
		daysAgo := daysAgo(listDir1[i], day)
		dcreat := dataCreate(listDir1[i])
		dir := listDir1[i]
		size := sizeFile(listDir1[i])
		str:=fmt.Sprintf("<tr><td align=\"left\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%d дней</td>" +
			"<td align=\"center\" style=\"width: 100px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 50px;\"><audio controls><source src=%s type=\"audio/wav\"></audio></td>" +
			"</tr>", dir, dcreat, daysAgo, size, dir)
		io.WriteString(w, str)
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
	//now := time.Now().Local()
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

func listfiles(rootpath string) []string {

	list := make([]string, 0, 10)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".wav" {
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