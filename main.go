package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// const baseDir = "e:/test/" // 指定要遍历的目录

var baseDir string

// go run main.go -dir=./files

func init() {
	// 使用 flag 包来处理命令行参数
	flag.StringVar(&baseDir, "dir", ".", "指定要遍历的目录，默认为当前目录")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", serveFiles)
	log.Println("Starting server on :26666")
	if err := http.ListenAndServe(":26666", nil); err != nil {
		log.Fatal(err)
	}
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// if path == "/" {
	// 	path = "/index.html"
	// }
	filePath := filepath.Join(baseDir, strings.TrimPrefix(path, "/"))

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if fileInfo.IsDir() {
		// 如果是目录，则列出目录内容
		listDirectory(w, filePath, path)
	} else {
		// 根据文件类型提供服务
		contentType := http.DetectContentType([]byte(fileInfo.Name()))
		switch {
		case strings.HasPrefix(contentType, "video/"):
			serveVideo(w, r, filePath)
		case strings.HasPrefix(contentType, "image/") || strings.HasPrefix(contentType, "text/"):
			http.ServeFile(w, r, filePath)
		default:
			serveDownload(w, r, filePath, fileInfo.Name())
		}
	}
}

func listDirectory(w http.ResponseWriter, dir, requestPath string) {
	files, _ := os.ReadDir(dir)
	fmt.Fprintf(w, "<h1>Directory: %s</h1>", requestPath)
	for _, file := range files {
		link := filepath.Join(requestPath, file.Name())
		if file.IsDir() {
			fmt.Fprintf(w, "<a href='%s'>%s/</a><br>", link, file.Name())
		} else {
			fmt.Fprintf(w, "<a href='%s'>%s</a><br>", link, file.Name())
		}
	}
}

func serveVideo(w http.ResponseWriter, r *http.Request, videoPath string) {
	http.ServeFile(w, r, videoPath)
	// 这里可以添加视频播放器的支持，例如使用 HTML5 视频标签
	// <video controls src="path_to_video"></video>
}

func serveDownload(w http.ResponseWriter, r *http.Request, filePath, fileName string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	http.ServeFile(w, r, filePath)
}
