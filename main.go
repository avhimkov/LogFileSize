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
)

func main() {
	runHTTP()
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

func ShowStat(w http.ResponseWriter, r *http.Request) {

	dirOkno1:= "file/dir1"
	//dirOkno2:= "file/dir2"
	//dirOkno3:= "file/dir3"
	//dirOkno4:= "file/dir4"
	//dirOkno5:= "file/dir5"

	render(w, "header.html")
	table(w, dirOkno1)
	render(w, "footer.html")
}

func table(w http.ResponseWriter, dir string) {

	//fs := http.FileServer(http.Dir("./file"))
	//http.Handle("/dir1/", http.StripPrefix("/dir1/", fs))
	listDir1 := listfiles(dir)

	for i := range listDir1 {
		dir :=listDir1[i]
		size := sizeFile(listDir1[i])
		//strDir := strings.Replace(dir, "\\", "/", -1)
		str:=fmt.Sprintf("<caption>%s</caption><tr>" +
			"<td align=\"left\" style=\"width: 300px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 300px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 300px;\"><audio controls><source src=%s type=\"audio/mpeg\"></audio></td>" +
			"</tr>", listDir1, dir, size, dir)
		io.WriteString(w, str)
	}
}

func runHTTP() {
	//fs := http.FileServer(http.Dir("./dir1/"))
	//http.Handle("/dir1/", http.StripPrefix("/dir1/", fs))
	http.HandleFunc("/", ShowStat)
	log.Println("localhost:8080 Listening...")
	http.ListenAndServe(":8080", nil)
	//http.FileServer(http.Dir("./file"))
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
func sizeFile(path string) string {

	stat, err := os.Stat(path)
	sizeStr, err := convertSize(stat.Size())
	if err != nil {
		sizeStr = "-"
	}
	return sizeStr
}

func listfiles(rootpath string) []string {

	list := make([]string, 0, 10)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".mp3" {
			list = append(list, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
	return list
}