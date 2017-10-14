package main

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for range tests {
		t.Run(t.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_conf(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf()
		})
	}
}

func Test_audioListen(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			audioListen(tt.args.w, tt.args.r)
		})
	}
}

func TestShowStat(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ShowStat(tt.args.w, tt.args.r)
		})
	}
}

func Test_render(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		tmpl string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			render(tt.args.w, tt.args.tmpl)
		})
	}
}

func Test_head(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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

func Test_tableAudio(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tableAudio(tt.args.w, tt.args.r)
		})
	}
}

func Test_tableMonitoring(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tableMonitoring(tt.args.w, tt.args.r)
		})
	}
}

func Test_listFiles(t *testing.T) {
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
		// TODO: Add test cases.
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

func Test_runHTTP(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for range tests {
		t.Run(tt.name, func(t *testing.T) {
			runHTTP()
		})
	}
}

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
		// TODO: Add test cases.
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
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		want  time.Time
		want1 float64
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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

func Test_checkErr(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkErr(tt.args.err)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnZip(tt.args.archive, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("UnZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
