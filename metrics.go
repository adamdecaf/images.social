package main

import (
	"expvar"
)

var (
	blacklistHits = expvar.NewInt("blacklist-hits")

	filesUploaded = expvar.NewInt("files-uploaded")
	failedUploads = expvar.NewInt("failed-uploads")
)

func MarkFileUpload() {
	filesUploaded.Add(1)
}

func MarkUploadFailed() {
	failedUploads.Add(1)
}

func MarkBlacklistHits() {
	blacklistHits.Add(1)
}
