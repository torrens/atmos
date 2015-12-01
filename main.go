package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Security struct {
	endpoint string
	secret   string
	uid      string
}

var (
	url         string
	secret      string
	uid         string
	atmosDir    string
	storagePath string
	Info        *log.Logger
	Error       *log.Logger
	wg          sync.WaitGroup
	msWait      int = 250
	tokens          = make(chan bool, 5)
)

func main() {
	initApp()

	// Test flags
	if testFlags(url, secret, uid, atmosDir, storagePath) == false || testStoragePath(storagePath) == false {
		return
	}

	// Init Loggers
	initLoggers()

	start := time.Now()
	Info.Println("Started")
	readDirectory(Security{endpoint: url, secret: secret, uid: uid}, atmosDir)
	wg.Wait()
	Info.Println("Complete", time.Since(start).String())
}

func initApp() {
	fmt.Println("Atmos reader 1.0.2")

	// Flags
	flag.StringVar(&url, "url", "", "The URL to the Atmos device in the form of https://some.host.com.")
	flag.StringVar(&secret, "secret", "", "The secret for your Atmos account.")
	flag.StringVar(&uid, "uid", "", "The User ID for your Atmos account.")
	flag.StringVar(&atmosDir, "atmosDir", "", "The Atmos directory you wish to read.")
	flag.StringVar(&storagePath, "storagePath", "", "The local directory to store the Atmos files.")
	flag.Parse()
}

func readDirectory(security Security, resource string) {
	data, err := request(security, "/rest/namespace/"+resource)
	if err != nil {
		Error.Println("Failed to read directory:", resource, err)
	}

	var directoryList DirectoryList = ParseDirectoryEntry(data)
	for _, directoryEntry := range directoryList.DirectoryEntry {
		if directoryEntry.IsDirectory() {
			// Pause before the next directory read to give the target server a rest
			time.Sleep(100 * time.Millisecond)
			readDirectory(security, resource+"/"+directoryEntry.Filename)
		} else {
			wg.Add(1)
			tokens <- true // acquire a token
			go readFile(security, directoryEntry.ObjectID, resource, directoryEntry.Filename)
		}
	}
}

func readFile(security Security, objectId string, resource string, fileName string) {
	defer func() {
		wg.Done()
		<-tokens
	}()
	url := "/rest/objects/" + objectId
	data, err := request(security, url)
	if err != nil {
		Error.Println("Failed to read file:", url, err)
	} else {
		createFile(storagePath+"/"+resource, fileName, data)
	}
}

func request(security Security, resource string) ([]byte, error) {
	now := time.Now().Format(time.RFC1123)
	headers := make(map[string]string)
	headers["x-emc-date"] = now
	headers["date"] = now
	headers["accept"] = "*/*"
	headers["x-emc-uid"] = security.uid

	hashString := hashString("GET", resource, headers)
	signature := ComputeHmac(hashString, security.secret)
	headers["x-emc-signature"] = signature

	client := &http.Client{}
	req, _ := http.NewRequest("GET", security.endpoint+resource, nil)
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Status code equals: " + strconv.Itoa(res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func hashString(httpRequestMethod string, resource string, headers map[string]string) string {
	hashString := httpRequestMethod +
		"\n\n\n" +
		headers["date"] +
		"\n" +
		strings.ToLower(resource) +
		"\n" +
		"x-emc-date:" + headers["x-emc-date"] +
		"\n" +
		"x-emc-uid:" + headers["x-emc-uid"]
	return hashString
}

func ComputeHmac(message string, secret string) string {
	key, _ := base64.StdEncoding.DecodeString(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func createFile(path string, fileName string, fileData []byte) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		Error.Println("Failed to create directory:", path, err)
		return err
	}

	filePath := path + "/" + fileName
	err = ioutil.WriteFile(filePath, fileData, 0644)
	if err != nil {
		Error.Println("Failed to write file:", filePath)
		return err
	}

	Info.Println("File created:", filePath)
	return nil
}
