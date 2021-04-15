package main

var (
	// files and directories info list when parsing the current directory
	filesList *FilesInfoList
	// current directory name without "/" tail end
	currentDir string
)

const (
	// size of file chunks, padding 0x0000
	chunkSize = 256000
	versionNo = "1.0.0"
	appName   = "manifest"
)

type ManifestOuputBody struct {
	ManifestHeader ManifestDirectoryHeader
	ManifestBody   ManifestDirectoryBody
	ChunkHashList  []FileChunkHashList
}

type ManifestFile struct {
	Path       []byte
	FileName   []byte
	Size       int64
	HashList   HashValue
	MetaString []byte
}

type ManifestDirectoryHeader struct {
	VersionString     []byte         `json:"version"`
	SequenceId        uint64         `json:"sequence"`
	CreatedAt         uint64         `json:"creation time"`
	Creator           string         `json:"creator"`
	BodySegmentLength uint64         `json:"file list length"`
	BodyDataFileSize  uint64         `json:"files total size"`
	MetaDataTags      KeysValuesList `json:"tags"`
	ChunkSize         int64          `json:"chunk size"`
}

type ManifestDirectoryBody struct {
	FileList []ManifestFile
}

type ManifestHeaderMetaData struct {
	CreationTime     uint64
	Creator          string
	PreviousManifest string
	SequenceId       uint64
	UniqueId         string
}

type FileChunkHashList struct {
	ChunksHashes [][]byte
}

type HashValue struct {
	HashType []byte
	Hash     []byte
}

type FilesInfoList struct {
	directoryNames  []string
	fileNames       []string
	diretorySizes   []int
	fileSizes       []int
	fileHashes      [][]byte
	fileschunkslist []FileChunkHashList
	filesMetaList   ManifestDirectMetaList
}

type KeyValueByte struct {
	Key   []byte
	Value []byte
}

type KeyValueString struct {
	Key   string
	Value string
}

type KeyValueList []KeyValueString

type KeysValuesList struct {
	Keys   [][]byte
	Values [][]byte
}

type FileMeta struct {
	CreateAt       uint64 `json:"creation time"`
	LastModified   uint64 `json:"last modified time"`
	UnixPermission string `json:"permission"`
}

type ManifestDirectMetaList []FileMeta

type FileData struct {
	FileName     string    `json:"name"`
	FileSize     int       `json:"size"`
	FileHash     []byte    `json:"hash"`
	FileMetaData *FileMeta `json:"meta,omitempty"`
}

type FileDataList []FileData

type DirectoryMeta struct {
	DirectoryName string `json:"name"`
	DirectorySize int    `json:"size"`
}
type DirectoryMetaList []DirectoryMeta

func (s *KeysValuesList) Add(pair KeyValueByte) {
	s.Keys = append(s.Keys, pair.Key)
	s.Values = append(s.Values, pair.Value)
}

func (s *KeysValuesList) KVRange() <-chan KeyValueByte {
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
	} else {
		return s[i].Value < s[j].Value
	}

}

func (s FileDataList) Len() int {
	return len(s)
}

func (s FileDataList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FileDataList) Less(i, j int) bool {
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
