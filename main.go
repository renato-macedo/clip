package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jtguibas/cinema"
	"github.com/rylio/ytdl"
)

func download(url string, start, end int) (path string, err error) {

	vid, err := ytdl.GetVideoInfo(url)
	if err != nil {
		log.Println("Failed to get video info")
		return "", err
	}
	name := fmt.Sprintf("downloads/%v.mp4", vid.Title) // vid.Title + ".mp4"

	file, _ := os.Create(name)
	defer file.Close()

	vid.Download(vid.Formats[0], file)
	video, err := cinema.Load(name)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if start != 0 {
		video.SetStart(time.Duration(start) * time.Second)
	}

	if end != 0 {
		video.SetEnd(time.Duration(end) * time.Second)
	}

	//video.Trim(st*time.Second, e*time.Second)
	path = fmt.Sprintf("videos/%v-%v-%v.mp4", vid.Title, start, end)
	video.Render(path)
	// log.Println("FFMPEG Command", video.CommandLine(path))
	return path, nil

}

func main() {

	// file server
	//fs := http.FileServer(http.Dir("./videos"))
	fs := dotFileHidingFileSystem{http.Dir("./videos")}

	http.Handle("/videos/", http.StripPrefix(strings.TrimRight("/videos/", "/"), http.FileServer(fs)))

	// download video route handler
	http.HandleFunc("/download", func(rw http.ResponseWriter, req *http.Request) {
		url := req.URL.Query().Get("url")

		start := req.URL.Query().Get("start")
		if start == "" {
			start = "0"
		}

		end := req.URL.Query().Get("end")
		if end == "" {
			end = "0"
		}

		log.Printf("start %v end %v \n", start, end)
		st, err := strconv.Atoi(start)
		if err != nil {
			http.Error(rw, "Invalid start time", http.StatusBadRequest)
			return
		}
		e, err := strconv.Atoi(end)
		if err != nil {
			http.Error(rw, "Invalid end time", http.StatusBadRequest)
			return
		}

		path, err := download(url, st, e)

		if err != nil {
			log.Println(err)
			http.Error(rw, "you messed up", http.StatusBadRequest)
			return
		}

		jsResp, err := json.Marshal(struct{ URL string }{URL: req.Host + "/" + path})
		if err != nil {
			panic(err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(jsResp)

	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
