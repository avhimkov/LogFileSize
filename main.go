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

//type Post struct{
//	Title string
//}

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

	render(w, "header.html")

	table(w, "file/dir1")
	table(w, "file/dir2")

	render(w, "footer.html")
}

func table(w http.ResponseWriter, dir string) {
	listDir1 := listfiles(dir)
	for i := range listDir1 {
		dir :=listDir1[i]
		size := sizeFile(listDir1[i])
		str:=fmt.Sprintf("<tr>" +
			"<td align=\"left\" style=\"width: 300px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 300px;\">%s</td>" +
			"<td align=\"center\" style=\"width: 300px;\"><audio controls><source src=%s type=\"audio/mpeg\"></audio></td>" +
			"</tr>", dir, size, "D:\\Music\\1.mp3")
		io.WriteString(w, str)
	}

}

func runHTTP() {
	http.HandleFunc("/", ShowStat)
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