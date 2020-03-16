package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jtguibas/cinema"
	"github.com/rylio/ytdl"
)

func download(url string, start, end int) (path string, err error) {
	log.Println("url", url)
	toBeDownloaded, err := ytdl.GetVideoInfo(url) // "https://www.youtube.com/watch?v=A02s8omM_hI"
	if err != nil {
		log.Println("Failed to get video info")
		return "", err
	}
	name := toBeDownloaded.Title + ".mp4"
	file, _ := os.Create(name)
	defer file.Close()

	toBeDownloaded.Download(toBeDownloaded.Formats[0], file)
	video, err := cinema.Load(name)
	if err != nil {
		log.Println(err)
		return "", err
	}

	video.SetStart(time.Duration(start))
	video.SetEnd(time.Duration(end))

	err = video.Render("output.mp4")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "output.mp4", nil

}

func main() {
	http.HandleFunc("/download", func(rw http.ResponseWriter, req *http.Request) {
		url := req.URL.Query().Get("url")
		start := req.URL.Query().Get("start")
		if start == "" {
			start = "0"
		}
		end := req.URL.Query().Get("end")
		if end == "" {
			end = "-1"
		}
		st, err := strconv.Atoi(start)
		if err != nil {
			http.Error(rw, "Invalid start time", http.StatusBadRequest)
		}
		e, err := strconv.Atoi(end)
		if err != nil {
			http.Error(rw, "Invalid end time", http.StatusBadRequest)
		}
		path, err := download(url, st, e)
		log.Println("path", path)
		if err != nil {
			log.Println(err)
			http.Error(rw, "you messed up", http.StatusBadRequest)
		}
		http.ServeFile(rw, req, "Lil Uzi Vert - Wassup feat. Future [Official Audio].mp4")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
