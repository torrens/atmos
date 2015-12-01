package main

import (
	"encoding/xml"
	"sort"
	"strconv"
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
	sort.Sort(ByFileName(listDirectoryResponse.DirectoryList.DirectoryEntry))
	return listDirectoryResponse.DirectoryList
}

type ByFileName []DirectoryEntry

func (f ByFileName) Len() int      { return len(f) }
func (f ByFileName) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f ByFileName) Less(i, j int) bool {
	first, _ := strconv.Atoi(f[i].Filename)
	second, _ := strconv.Atoi(f[j].Filename)
	return  first < second
}
