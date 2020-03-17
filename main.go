package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {

	// file server
	//fs := http.FileServer(http.Dir("./videos"))
	fs := dotFileHidingFileSystem{http.Dir("./videos")}

	http.Handle("/videos/", http.StripPrefix(strings.TrimRight("/videos/", "/"), http.FileServer(fs)))

	// download video route handler
	http.HandleFunc("/download", downloadHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
