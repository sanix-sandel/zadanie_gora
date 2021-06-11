package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Image struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Size  int    `json:"size"`
}

var DB *sql.DB

func uploadFile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowd", 405)
		return
	}

	//limiting the file size
	r.ParseMultipartForm(10 << 30)

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("Error")
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := strings.Split(header.Filename, ".")[0]
	filename := path.Join("files", header.Filename)

	//wtrite data to file with permission
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//create an image with data got from file uploaded
	img := Image{Title: title, Url: filename, Size: int(header.Size)}

	fmt.Println(img)
}

func getFiles(w http.ResponseWriter, r *http.Request) {

}

func deleteFile(w http.ResponseWriter, r *http.Request) {

}

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")

	srv := &http.Server{
		Addr:    *addr,
		Handler: routes(),
	}

	log.Printf("Server started listening on localhost:8000")
	log.Fatal(srv.ListenAndServe())
}
