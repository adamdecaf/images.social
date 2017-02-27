package cache

import (
	"os"
)

const (
	localFSCachePath = "./files"
)

// Dir returns the on-disk location for storage. This is exposed so
// the main http handler can redirect requests off to the local fs
func Dir() string {
	return localFSCachePath
}

// Init sets up the cache.
// This right now just sets up the on-disk folder for upload storage
func Init() error {
	_, err := os.Stat(localFSCachePath)
	if err == nil {
		return nil
	}
	err = os.Mkdir(localFSCachePath, 0744)
	if err != nil {
		return err
	}
	return nil
}
