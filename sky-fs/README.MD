
#Manifest Internal Representation

---

Three Data Structures
- list of StorageBlocks
- list of hashes (kv data; chunks)
- list of files

Storage Blocks
- are disc drive segments, memory mapped files for large storage
- are immutable
- EX. 1 GB, 10 TB = 10,000 Storage Blocks
- EX. 32 GB, 10 TB = 300 Storage Blocks

A storage block has
- a hash for the segment (merkle)
- a merkle tree map (can be modified and update hash between writes
- a list of keys (a map) at the start

---

Chunks allow indexing without knowing file list
- encryption at rest

Merkle Hashes of size 256 kb, means
- 


---

Manifest will read the .cxo file in and create an "internal representation"

There is a list of Manifest files (.cxo) files that are loaded into memory at any time for a cxo daemon.

"cxo run server"
- starts CXO server
- cmd/cxo

"cxo add [manifest file .cxo"
- add a .cxo manifest file to the internal representation

"cxo list"
- lists manifest files that are loaded by the server

---

A CXO Server uses the following data structures

```golang

struct FileIndexList {
	 FileIndexList []FileIndex
}

struct FileIndex {
	FileInternalId int64 //sequential internal, integer id
	FileInfoIndex int64 //index into FileInfo Array
	FileSize int64 //size of file in bytes

	LocalDataAmount //How many Chunks we have of file locally

	FileHashListIndex int64 //starting index for file hashes
	//FileHashListLen int64 //number of hashes or chunks for this file
}
```

```golang
struct FileInfo {
	[]byte filepath
	[]byte filename
	//increment Reference count for each manifest file, that contains this filename/filepath
	//decrement reference count for each file in each manifest unloaded from the manifest list
	ReferenceCount int32 // how many manifest files use this
}

struct FileInfoStorage {
	//index for int64 to FileInfo Object
	FileInfoArray []FileInfo
	//maps "filename+filepath" to an integer index into FileInfoArray
	FileIndexMap map[[]byte] int
}
```

```golang
/* Data Storage Structures for Lists of Hashes */

struct FileHashListIndex {
	FileHashListArray[] FileHash
	//slice length
	FileHashListArraySize int64 //number of elements in array
	//slice max before resizing
	FileHashListArrayMax //number of elements until we have to resize array
}

(x *FileHashListIndex ) Get(int64 index) *FileHashList {
	return FileHashListArray[index]
}

struct FileHashList {
	ChunkSize int32 //bytes

}


```