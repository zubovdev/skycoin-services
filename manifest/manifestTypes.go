package main

import "time"

type ManifestOuputBody struct{
	ManifestHeader 				ManifestDirectoryHeader
	ManifestBody 				ManifestDirectoryBody
}

type ManifestFile struct{
	Path     	[]byte  
	FileName 	[]byte  
	Size     	int64  
    HashList 	[]HashType
    MetaString	[]byte
}

type ManifestDirectoryHeader struct{
	VersionString       		[]byte
	SequenceId           		uint64
	CreatedAt					time.Time
	BodySegmentLength	        uint64
	BodyDataFileSize	        uint64
	SerializedMapList			SerializedKvList
}

type  ManifestDirectoryBody struct{ 
	FileList []ManifestFile 
}

type ManifestHeaderMetaData struct{
	CreationTime				int64
	Creator						string
	PreviousManifest			string
	SequenceId           		uint64
	UniqueId					string
}

type HashType struct{
    HashType []byte 
    Hash     []byte 
}

type FilesMetaList struct{
	directoryNames		[]string
	fileNames       	[]string
	diretorySizes       []int
	fileSizes    		[]int
	fileHashes   		[][]byte
}

type KeyValueByte struct{ 
    Key []byte 
    Value []byte 
}

type KeyValueString struct{ 
    Key string 
    Value string 
}

type KeyValueList []KeyValueString

type SerializedKvList struct{
	Keys [][]byte
	Values [][]byte
}



type FileMeta struct {
	FileName  string        `json:"name"`
	FileSize  int 			`json:"size"`
	FileHash  []byte      	`json:"hash"`
}

type FileMetaList []FileMeta

type DirectoryMeta struct {
	DirectoryName  string       `json:"name"`
	DirectorySize  int 			`json:"size"`
}
type DirectoryMetaList []DirectoryMeta

var (
	cxoFileName string
	filesList *FilesMetaList
)

 

func (s *SerializedKvList) Add(pair KeyValueByte) {
	s.Keys = append(s.Keys, pair.Key)
	s.Values = append(s.Values, pair.Value)
}

func (s *SerializedKvList) KVRange() <-chan KeyValueByte {
	chnl := make(chan KeyValueByte)
	limit := len(s.Keys)
	go func() {
		for i := 0; i < limit; i++ {
			chnl <- KeyValueByte{s.Keys[i], s.Values[i]}
		}

		close(chnl)
	}()
	return chnl
}

func (s KeyValueList) Len() int {
	return len(s)
}

func (s KeyValueList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s KeyValueList) Less(i, j int) bool {
	if s[i].Key != s[j].Key {
		return s[i].Key < s[j].Key
	}else{
		return s[i].Value < s[j].Value
	}
	
}

func (s FileMetaList) Len() int {
	return len(s)
}

func (s FileMetaList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FileMetaList) Less(i, j int) bool {
	return s[i].FileName < s[j].FileName
}

func (s DirectoryMetaList) Len() int {
	return len(s)
}

func (s DirectoryMetaList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DirectoryMetaList) Less(i, j int) bool {
	return s[i].DirectoryName < s[j].DirectoryName
}