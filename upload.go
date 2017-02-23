package main

import (
	"log"
	"io"
	// "io/ioutil"
	"errors"
	"path"
	"path/filepath"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	MaxInmemUploadSize = 0 // Force file uploads to server disk
)

var (
	CopyFailed = errors.New("error copying file")
)

func uploadRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(MaxInmemUploadSize)

		file, _, err := r.FormFile("file")
		defer func() {
			if file != nil {
				file.Close()
			}
		}()
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Grab the file and copy over to our cache location
		// todo: verify it a bit
		err = cp(file, path.Join(LocalFSCachePath, "/foo-random"))
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
	}

	//
}

func cp(src multipart.File, d string) error {
	if strings.HasPrefix(d, "..") {
		log.Println("A")
		return CopyFailed
	}

	// abs dst path
	dst, err := filepath.Abs(d)
	if err != nil {
		log.Println("B")
		return CopyFailed
	}

	// Check the out file doesn't exist
	_, err = os.Stat(dst)
	if err == nil {
		log.Println("C")
		return CopyFailed
	}

	// Open out file
	out, err := os.Create(dst)
	if err != nil {
		log.Println(dst)
		log.Println("D")
		return CopyFailed
	}
	defer out.Close()

	// Copy the file contents
	_, err = io.Copy(out, src)
	cerr := out.Close()
	if err != nil {
		log.Println("E")
		return CopyFailed
	}
	return cerr
}
