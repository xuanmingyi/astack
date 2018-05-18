package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/xuanmingyi/astack/models"

	"io"
	"log"
	"net/http"
	"os"
)

// 获取大小的借口
type Sizer interface {
	Size() int64
}

// hello world, the web server
func ImageController(w http.ResponseWriter, r *http.Request) {
	if "POST" == r.Method {
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer file.Close()
		f, err := os.Create("filenametosaveas")
		defer f.Close()
		io.Copy(f, file)
		return
	}

	// 上传页面
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	html := `
<form enctype="multipart/form-data" action="/hello" method="POST">
    Send this file: <input name="userfile" type="file" />
    <input type="submit" value="Send File" />
</form>
`
	io.WriteString(w, html)
}

func VolumeController(w http.ResponseWriter, r *http.Request) {
}

func ServerController(w http.ResponseWriter, r *http.Request) {
}

func NetworkController(w http.ResponseWriter, r *http.Request) {
}

func main() {

	db, err := gorm.Open("mysql", "root@/mm?charset=utf8")
	if err != nil {
		panic(fmt.Sprintf("%s: %s", err, "failed to connect database"))
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&models.Image{})

	http.HandleFunc("/image", ImageController)
	http.HandleFunc("/volume", VolumeController)
	http.HandleFunc("/server", ServerController)
	http.HandleFunc("/network", NetworkController)

	err2 := http.ListenAndServe(":9090", nil)
	if err2 != nil {
		log.Fatal("ListenAndServe: ", err2)
	}
}
