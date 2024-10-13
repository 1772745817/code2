package main

import (
	"net/http"
	"os"
)

func filehandel(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("hello.txt")
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Unable to get file info", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}
func main() {
	http.HandleFunc("/", filehandel)
	http.ListenAndServe(":8080", nil)
}
