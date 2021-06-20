package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"images_upload/dbutils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
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
		log.Println("Error uploading data")
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
	img := Image{Title: title, Url: "/" + filename, Size: int(header.Size)}

	//insert to the db
	statement, _ := DB.Prepare("INSERT INTO images (title, url, size) VALUES (?, ?, ?)")
	result, err := statement.Exec(img.Title, img.Url, img.Size)
	if err == nil {
		newId, _ := result.LastInsertId()
		img.Id = int(newId)
		w.Header().Set("Content-Type", "application/json")
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getFiles(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowd", 405)
		return
	}

	rows, err := DB.Query("SELECT * FROM images")
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	var images []Image

	//Iterate in order to get all records
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.Id, &img.Title, &img.Url, &img.Size); err != nil {
			log.Println("Error fetching from database")
		}
		images = append(images, img)
	}

	//encode the records slice into json
	resp_JSON, _ := json.Marshal(images)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp_JSON)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {

	//getting the id parameter
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		log.Println("Not found !")
		return
	}

	//deleting file from the directory
	var img Image
	err = DB.QueryRow("SELECT * FROM images WHERE id=?", id).Scan(
		&img.Id, &img.Title, &img.Url, &img.Size,
	)
	if err != nil {
		log.Println("File image not found in the directory")
		w.Write([]byte("File not found"))
		return
	}
	err = os.Remove(img.Url[1:])
	if err != nil {
		log.Println("Error, file not in directory")
		return
	}

	//deleting file from db
	deletestmt, _ := DB.Prepare("DELETE FROM images WHERE id=?")
	_, err = deletestmt.Exec(id)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Image deleted "))

}

func main() {
	addr := flag.String("addr", ":"+os.Getenv("PORT"), "HTTP network address")

	var err error

	DB, err = sql.Open("sqlite3", "./database.db")

	if err != nil {
		log.Println("Driver creation failed")
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: routes(),
	}
	dbutils.Initialize(DB)

	log.Printf("Server started listening on localhost:8000")
	log.Fatal(srv.ListenAndServe())
}
