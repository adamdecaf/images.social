package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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

		// Get the file extension from reading the file
		ext := Detect(file).Ext()
		log.Printf("ext = '%s'\n", ext)
		if ext == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Grab the file and copy over to our cache location
		out := temp()
		tmp := path.Join(LocalFSCachePath, "/"+out)
		hash, err := cpAndHash(file, tmp)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Move from temp to name under hash
		final := path.Join(LocalFSCachePath, hash+ext)
		log.Println(final)
		if _, err := os.Stat(final); err != nil {
			os.Rename(tmp, final)
		}
	}
	http.Error(w, "bad request", http.StatusBadRequest)
}

// cpAndHash copies the multipart file to a non-temp file on disk as well as
// returns the hash (Sha256) of the file contents
func cpAndHash(src multipart.File, d string) (string, error) {
	if strings.HasPrefix(d, "..") {
		log.Println("A")
		return "", CopyFailed
	}

	// abs dst path
	dst, err := filepath.Abs(d)
	if err != nil {
		log.Println("B")
		return "", CopyFailed
	}

	// Check the out file doesn't exist
	_, err = os.Stat(dst)
	if err == nil {
		log.Println("C")
		return "", CopyFailed
	}

	// Open out file
	out, err := os.Create(dst)
	if err != nil {
		log.Println(dst)
		log.Println("D")
		return "", CopyFailed
	}
	defer out.Close()

	// Prep hash
	h := sha256.New()
	t := io.TeeReader(src, h)

	// Copy the file contents
	_, err = io.Copy(out, t)
	cerr := out.Close()
	if err != nil {
		log.Println("E")
		return "", CopyFailed
	}
	return hex.EncodeToString(h.Sum(nil)), cerr
}

func temp() string {
	h := sha256.New()
	h.Write([]byte(time.Now().String())) // todo: something that's not time based?
	return hex.EncodeToString(h.Sum(nil))
}
