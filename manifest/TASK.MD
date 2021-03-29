

# CXO/Manifest Tool Specifications

---

Manifest is a format for CXO for storing information about files.

This is a utility that takes in a directory and outputs a struct describing the files
- file path
- file name
- file size
- file hash

---

Write a command line utility
- under cmd/manifest
- code in manifest

- use "flag" for CLI or https://github.com/urfave/cli

---

- Manifest takes a directory as input

- Manifest recursively enumerates all files and directories for the input directory

- Manifest builds an internal Struct with list
  - list of directories
  - list of files (path and filename)
  - the size of each file

-ManifestDirectoryHeader
  - []byte VersionString
  - SequenceId
  - file must include, unix 64 bit second time (CreatedAt)
  - length of BodySegment (uint64) //size of serialized file list
  - BodyDataFile Size (uint64) // 
  - HashList SerializedMap
    
A "ManifestHeaderMetaData" string containing:
  - CreationTime int64 unix time in base10
  - creator string
  - PreviousManifest
  - "SequenceID" (uint64)
  - UniqueId (string)

|The Metadata will be a string or []byte of a map[string]string that is serialized with https://godoc.org/github.com/skycoin/skycoin/src/cipher/encoder 
- maps are serialized by sorting the key to map in alphabetical order (and if input keys are the same, sort by key, then by value)
- make a "serialize map" function that packs a map with 

``` golang
struct KeyValue { 
    key []byte, 
    value []byte
}
```

And pack these structs into an array, sorting elements by input key. Then serialize the resulting array with encoder.

``` golang
struct SerializedKvList {
   Keys []byte
   Values []byte
}
```

Also note, that serialized maps can have multiple key-value pairs with same value
- have function for add key-value pair to list
- have function for iterating over all key-value pairs
- have function for sorting the key-value list by the byte byte byte order of the concatentation of the serialized key/value pair

---

For Directories
- the "filename" variable is empty (nil)
- the "size" is the recursive sum of all files inside of the directory (computed in second pass)

---

struct Body ManifestDirectoryBody {
   FileList []ManifestFile
}

---

- Manifest has option -print-json to print a JSON printable version of the
  - a list of files (path + filename) in seperate sections as strings
    
``` golang
struct ManifestFile {
  []byte path
  []byte filename
  size int64
  HashList []HashType
  MetaString[] byte
```

```golang
struct HashType {
    []byte HashType
    []byte Hash
}
```
HashType can be
- "base64,sha256"
- "base64,blake3"
- "bin,sha256"
- "bin,blake3"

---

```golang
Struct ManifestOuputBody
- FileList []ManifestFile
- MetaDataList []MetaData (array of meta-data structs)
- []KeyValue
```

---
- Manifest has option -meta to also construct a  new json section which has
  - creation date of file (uint64 unix time)
  - last modified date of file (uint64 unix time)
  - unix file permisions (read, write, execute, etc)
  - store this data in a "ManifestDirectMetaList"    

And to write these in a  SerializedKvMap to MetaString[] for each file

---

``` 
manifest init
```

Creates a .cxo file in current directory
- "manifest commit" will build a manifest file and save to filename [unixtime in seconds].cxo file (serialized ManifestDirectoryHeader) in the .cxo/checkpoints/ folder
- "manifest verify" will check each file in directory against the latest manifest
- "manifest list" will list all manifest files and information such as total size of header, size of files

---

#Testing

- generate a new directory ("testdata")
- put testdata in .gitignore
- populate testdata with random directories and files (size min, size max; number of files, max directory depth, etc)
- compute the manifest file
- then run "manifest verify" function on directory and check that all the file names exist, the file size is correct and that the hashes of the files match


