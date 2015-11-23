package main

import (
	"flag"
	"fmt"
	"os"
)

func testFlags(url, secret, uid, atmosDir, storagePath string) bool {
	if len(url) == 0 || len(secret) == 0 || len(uid) == 0 || len(atmosDir) == 0 || len(storagePath) == 0 {
		Error.Println("Missing arguments.")
		flag.Usage()
		return false
	}
	return true
}

func testStoragePath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		msg := fmt.Sprintf("Path: %s doesn't exists.", path)
		Error.Println(msg)
		fmt.Println(msg)
	}
	return err == nil
}