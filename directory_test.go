package main

import "testing"

var testXML = "<ListDirectoryResponse xmlns='http://www.emc.com/cos/'>" +
	"	<DirectoryList>" +
	"		<DirectoryEntry>" +
	"			<ObjectID>50f59fe8a108a310050f5aa980360a0547dc75d64442</ObjectID>" +
	"			<FileType>regular</FileType>" +
	"			<Filename>someFile.txt</Filename>" +
	"		</DirectoryEntry>" +
	"	</DirectoryList>" +
	"</ListDirectoryResponse>"

func TestDirectoryParse(t *testing.T) {
	directoryList := ParseDirectoryEntry([]byte(testXML))
	directoryEntries := len(directoryList.DirectoryEntry)
	if directoryEntries != 1 {
		t.Error("Expected a directoryList with 1 item, got ", directoryEntries)
	}

	firstDirectoryEntry := directoryList.DirectoryEntry[0]
	if firstDirectoryEntry.IsDirectory() {
		t.Error("Expected false, got", true)
	}
}
