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

func download(url string, start, end int) (string, error) {

	split := strings.Split(url, "/watch?v=")
	videoID := split[1]

	filepath := fmt.Sprintf("downloads/%v.mp4", videoID)

	exists, err := checkIfFileExistsInDir(videoID+".mp4", "./downloads")
	if err != nil {
		return "", err
	}
	if exists == false {
		vid, err := ytdl.GetVideoInfo(url)
		if err != nil {
			log.Println("Failed to get video info")
			return "", err
		}
		file, _ := os.Create(filepath)
		defer file.Close()
		vid.Download(vid.Formats[0], file)
	}

	toBeRendered := fmt.Sprintf("%v-%v-%v.mp4", videoID, start, end)
	exists, err = checkIfFileExistsInDir(toBeRendered, "./videos")
	if err != nil {
		return "", err
	}
	toBeRendered = "videos/" + toBeRendered
	if exists == true {
		return toBeRendered, nil
	}

	video, err := cinema.Load(filepath)
	if err != nil {
		return "", err
	}

	if start != 0 {
		video.SetStart(time.Duration(start) * time.Second)
	}

	if end != 0 {
		video.SetEnd(time.Duration(end) * time.Second)
	}

	video.Render(toBeRendered)
	// log.Println("FFMPEG Command", video.CommandLine(toBeRendered))
	return toBeRendered, nil

}

func downloadHandler(rw http.ResponseWriter, req *http.Request) {
	url := req.URL.Query().Get("url")

	start := req.URL.Query().Get("start")
	if start == "" {
		start = "0"
	}

	end := req.URL.Query().Get("end")
	if end == "" {
		end = "0"
	}

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

}

func checkIfFileExistsInDir(filename, dirname string) (bool, error) {
	dir, err := os.Open(dirname)
	if err != nil {
		return false, err
	}
	defer dir.Close()
	filenames, err := dir.Readdirnames(0)
	for _, name := range filenames {
		if name == filename {
			return true, err
		}
	}
	return false, err
}
