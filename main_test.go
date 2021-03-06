package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_render(t *testing.T) {

	rr := httptest.NewRecorder()

	type args struct {
		w    http.ResponseWriter
		tmpl string
	}
	tests := []struct {
		name string
		args args
	}{
		{"tableMonitoring", args{w: rr, tmpl: "header.html"}},
		{"tableMonitoring", args{w: rr, tmpl: "footer.html"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			render(tt.args.w, tt.args.tmpl)
		})
	}
}

func Test_head(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	//handler := http.HandlerFunc(Index)

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "tableMonitoring", args: args{w: rr, r: req}, want: ""},
		//{name: "tableMonitoring",  args: args{w: rr, r: req}, want:"2017-03-03"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := head(tt.args.w, tt.args.r); got != tt.want {
				t.Errorf("head() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_htmlRang(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 string
	}{
		{"tableMonitoring", args{w: rr, r: req}, "", "", ""},
		//{"tableMonitoring", args{w: rr, r: req}, "2017-03-03","Окно №3", "3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := htmlRang(tt.args.w, tt.args.r)
			if got != tt.want {
				t.Errorf("htmlRang() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("htmlRang() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("htmlRang() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

//func Test_tableAudio(t *testing.T) {
//	req, _ := http.NewRequest("GET", "/", nil)
//	rr := httptest.NewRecorder()
//
//	type args struct {
//		w http.ResponseWriter
//		r *http.Request
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{"tableMonitoring", args{w: rr, r: req}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tableAudio(tt.args.w, tt.args.r)
//		})
//	}
//}

//func Test_tableMonitoring(t *testing.T) {
//	req, _ := http.NewRequest("GET", "/", nil)
//	rr := httptest.NewRecorder()
//	type args struct {
//		w http.ResponseWriter
//		r *http.Request
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{"tableMonitoring", args{w: rr, r: req}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tableMonitoring(tt.args.w, tt.args.r)
//		})
//	}
//}

/* func Test_listFiles(t *testing.T) {
	type args struct {
		rootpath string
		typefile string
		data     string
		time     string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 []string
		want2 []string
	}{
		{"ConvertSize", args{"D:/blabla/","zip","",""},
			[]string{},[]string{}, []string{},},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := listFiles(tt.args.rootpath, tt.args.typefile, tt.args.data, tt.args.time)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listFiles() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("listFiles() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("listFiles() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
*/
func TestConvertSize(t *testing.T) {
	type args struct {
		size int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"ConvertSize", args{7467273}, "7 MB", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertSize(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateCreate(t *testing.T) {

	loc, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		fmt.Println(err)
	}

	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		want  time.Time
		want1 float64
	}{
		{"DateCreate", args{"D:/blabla/Окно №3/25_20161102-00139_02-11-2016_19-28.zip"}, time.Date(2016, 11, 02, 19, 28, 17, 392856000, loc), 351.14669147491895},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := DateCreate(tt.args.path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateCreate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DateCreate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSizeFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int64
	}{
		{"Size", args{"D:/blabla/Окно №3/25_20161102-00139_02-11-2016_19-28.zip"}, "7 MB", 7467273},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SizeFile(tt.args.path)
			if got != tt.want {
				t.Errorf("SizeFile() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SizeFile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUnZip(t *testing.T) {
	type args struct {
		archive string
		target  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{"D:/blabla/Окно №3/25_20161102-00139_02-11-2016_19-28.zip", "file/temp"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnZip(tt.args.archive, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("UnZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
