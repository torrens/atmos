package main

import (
	"encoding/xml"
)

type ListDirectoryResponse struct {
	DirectoryList DirectoryList
}

type DirectoryList struct {
	DirectoryEntry []DirectoryEntry
}

type DirectoryEntry struct {
	ObjectID string
	FileType string
	Filename string
}

func (this DirectoryEntry) String() string {
	return this.Filename + " " + this.FileType + " " + this.ObjectID
}

func (this DirectoryEntry) IsDirectory() bool {
	return this.FileType == "directory"
}

func ParseDirectoryEntry(data []byte) DirectoryList {
	var listDirectoryResponse ListDirectoryResponse
	xml.Unmarshal(data, &listDirectoryResponse)
	return listDirectoryResponse.DirectoryList
}
