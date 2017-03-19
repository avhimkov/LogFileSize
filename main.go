package main

import (
	"fmt"
	"math"
	"net/http"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func main() {
	runHTTP()
}

type PageVariables struct{
	Data string
	size string
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

	//files, _ := filepath.Glob("*")
	//fmt.Println(files) // contains a list of all files in the current directory
	//now := strftime.Format("%Y-%m-%d", time.Now())
	//получаем статистику за дату
	//val, _ := tosserstat.Dates[now]
	//сортируем папки по имени
	//var keys []string
	//for k := range val {
	//	keys = append(keys, k)
	//}
	//sort.Strings(keys)
	//render(w, "hello")
	//for _, dir := range keys {
	//	//статистика для каталога-источника
	//	dirStat, ok := val[dir]
	//	LastProcessingDateStr := "-"
	//	if ok {
	//		LastProcessingDateStr = strftime.Format("%H:%M:%S", time.Unix(dirStat.LastProcessingDate, 0))
	//	}
	//}

	render(w, "hello.html")
	size := sizeFile("file/1.mp3")
	trs := fmt.Sprint("<tr><td>%s</td><td>\n",  size)

	fmt.Fprintf(w, trs)
	//io.WriteString(w, trs)
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

	//file, err := os.Open("file/1.txt")
	//if err != nil {
	//	return
	//}
	//defer file.Close()
	//stat, err := os.Stat("file/1.mp3")

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