package main

import "net/http"

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getFiles)
	mux.HandleFunc("/delete", deleteFile)
	mux.HandleFunc("/upload", uploadFile)

	fileServer := http.FileServer(http.Dir("./files"))
	mux.Handle("/files/", http.StripPrefix("/files", fileServer))

	return mux
}
