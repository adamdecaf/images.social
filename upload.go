package main

// todo:
// - can we disable the default /debug/vars handler?

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"expvar"
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
	maxInmemUploadSize = 0 // Force file uploads to server disk
)

var (
	blacklist     Blacklist
	errCopyFailed = errors.New("error copying file")

	// Metrics
	blacklistHits = expvar.NewInt("blacklist-hits")
	filesUploaded = expvar.NewInt("files-uploaded")
	failedUploads = expvar.NewInt("failed-uploads")
)

func init() {
	bl, err := NewBlacklist("./blacklist")
	if err != nil {
		panic(err)
	}
	blacklist = bl
}

func uploadRoute(w http.ResponseWriter, r *http.Request) {
	if blacklist.Blocked(*r) {
		markBlacklistHits()
		fail(w)
	}

	if r.Method == "POST" {
		r.ParseMultipartForm(maxInmemUploadSize)

		file, _, err := r.FormFile("file")
		defer func() {
			if file != nil {
				file.Close()
			}
		}()
		if err != nil {
			fail(w)
			return
		}

		// Grab the file and copy over to our cache location
		out := temp()
		tmp := path.Join(LocalFSCachePath, "/"+out)
		hash, err := cpAndHash(file, tmp)
		if err != nil {
			fail(w)
			return
		}

		// Get the file extension from reading the file
		f, err := os.Open(tmp)
		if err != nil {
			fail(w)
			return
		}
		defer f.Close()
		tpe, err := Detect(f)
		if err != nil {
			fail(w)
			return
		}
		ext := tpe.Ext()

		// Move from temp to name under hash
		final := path.Join(LocalFSCachePath, hash+ext)
		if _, err := os.Stat(final); err != nil {
			os.Rename(tmp, final)
			markFileUpload()
		}
		http.Redirect(w, r, path.Join("/i/", hash+ext), http.StatusFound)
		return
	}
	fail(w)
}

// cpAndHash copies the multipart file to a non-temp file on disk as well as
// returns the hash (Sha256) of the file contents
func cpAndHash(src multipart.File, d string) (string, error) {
	if strings.HasPrefix(d, "..") {
		log.Println("A")
		return "", errCopyFailed
	}

	// abs dst path
	dst, err := filepath.Abs(d)
	if err != nil {
		log.Println("B")
		return "", errCopyFailed
	}

	// Check the out file doesn't exist
	_, err = os.Stat(dst)
	if err == nil {
		log.Println("C")
		return "", errCopyFailed
	}

	// Open out file
	out, err := os.Create(dst)
	if err != nil {
		log.Println(dst)
		log.Println("D")
		return "", errCopyFailed
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
		return "", errCopyFailed
	}
	return hex.EncodeToString(h.Sum(nil)), cerr
}

func temp() string {
	h := sha256.New()
	h.Write([]byte(time.Now().String())) // todo: something that's not time based?
	return hex.EncodeToString(h.Sum(nil))
}

func fail(w http.ResponseWriter) {
	markUploadFailed()
	http.Error(w, "bad request", http.StatusBadRequest)
}

// Metrics collecting
func markFileUpload()    { filesUploaded.Add(1) }
func markUploadFailed()  { failedUploads.Add(1) }
func markBlacklistHits() { blacklistHits.Add(1) }
