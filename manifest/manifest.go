package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"time"
	"github.com/urfave/cli/v2"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)


func initCLI() *cli.App{
	filesList = processDirAndGenerateMeta()

	app := cli.NewApp()
	app.Name = "manifest"
	app.Usage = "create manifest files in current directory"
	app.Version = "1.0.0"
	addCLICommands(app)

	app.Action = func(cnx *cli.Context) error {
		cli.ShowAppHelpAndExit(cnx, 0)
		return nil
	}
	 
	sort.Sort(cli.FlagsByName(app.Flags))
	cli.VersionFlag = &cli.BoolFlag {
		Name: "print-version",
		Usage: "print version",
	}

	return app
}

func addCLICommands(app *cli.App){
	app.Commands = []*cli.Command{
		{
			Name: "init",
			Usage: "initialize tool environment by create the .cxo folder",
			UsageText: "create the manifest foler .cxo in current directory",
			Action: func(cnx *cli.Context) error {
				createCXOdir()
				return nil
			},
		},
		{
			Name: "commit",
			Usage: "commit all the files' metadatum into the .cxo file",
			UsageText: "commit all the metadata files into the .cxo folder",
			Flags: []cli.Flag{
				&cli.BoolFlag {
					Name: "print-json",
					Value: false,
					Usage: "print files in the directory in json ",
				},
				&cli.BoolFlag {
					Name: "meta",
					Value: false,
					Usage: "add metadata section to the file list",
				},
			},
			Action: func(cnx *cli.Context) error {
				// metaFlag := false
				// if cnx.Bool("meta"){
				// 	metaFlag = true
				// }
				if !isCXOFolderExist(){
					panic("folder ./cxo does not exist")
				}

				err := os.MkdirAll("./.cxo/checkpoints/",os.ModePerm)
				if err != nil {
					panic(err)
				}
				cxoFileName = generateCXOFilename()
				
				cxoFile, err := os.OpenFile(cxoFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					panic(err)
				}
				defer cxoFile.Close()

				mapList := generateFilesMaps(filesList)
				serializedkeyValueList := serializeMaps(mapList)
				serializedListResult := encoder.Serialize(*serializedkeyValueList)

				_, err = cxoFile.Write(serializedListResult)
				if err != nil {
					panic(err)
				}
				if cnx.Bool("print-json"){
					printFilesInJson(filesList)
				}
				// header := getManifestHeaderMetaData()

				// body := getManifestBody(filesList)
				// directoryHeader := getManifestDirectoryHeader(serializedkeyValueList, body)
				// OuputBody := ManifestOuputBody{
				// 	ManifestHeader: *directoryHeader,
				// 	ManifestBody: *body,
				// }

				// serializedBodyResult := encoder.Serialize(OuputBody)
				// _, err = cxoFile.Write(serializedBodyResult)
				// if err != nil {
				// 	panic(err)
				// }
				return nil
			},
		},
	}
}

func main(){
	defer func() {  
		if err := recover(); err != nil {
		   fmt.Println(err)
		   fmt.Println("please use 'manifest init' command before 'manifest commit'") 
		   os.Exit(1)
		}
	}()

	app := initCLI()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	} 

}

func processDirAndGenerateMeta() *FilesMetaList {
	var FilesAndDirectories FilesMetaList
	var directories []string
	var directoriesSize []int
	var files []string
	var filesSize []int
	var filesHash [][]byte

	err := filepath.Walk(".",
    func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directories = append(directories, path)
			dirSize,err := getDirectorySize(path)
			if err != nil {
				return err
			}
			directoriesSize = append(directoriesSize, dirSize)
		}else{
			files = append(files, path)
			filesSize = append(filesSize, int(info.Size()))
			filesHash = append(filesHash, []byte(hashFileAndEncoding(path)))
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	FilesAndDirectories.directoryNames = directories
	FilesAndDirectories.fileNames = files
	FilesAndDirectories.fileSizes = filesSize
	FilesAndDirectories.fileHashes = filesHash
	FilesAndDirectories.diretorySizes = directoriesSize
	return &FilesAndDirectories
}

 
func getDirectorySize(directory string) (int,error){
	totalSize := 0
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				totalSize =+ int(info.Size())
				
			} 
			return nil
		})
	return totalSize,err
}

func printFilesInJson(fList *FilesMetaList){
	var dirmeta DirectoryMetaList
	var filemeta FileMetaList
	 
	for indx,fn := range (*fList).fileNames{
		fh := (*fList).fileHashes[indx]
		fs := (*fList).fileSizes[indx]
		fileInfo := FileMeta{ fn, fs, fh }
		filemeta = append(filemeta, fileInfo)
		
	}

	for indx,dn := range (*fList).directoryNames{
		ds := (*fList).diretorySizes[indx]
		dirInfo := DirectoryMeta{ dn, ds }
		dirmeta = append(dirmeta, dirInfo)
	}
	
	sort.Sort(filemeta)
	sort.Sort(dirmeta)
	metadata := struct { 
		Files  []FileMeta   `json:"files"`
		Directories []DirectoryMeta `json:"directories"`
	}{ filemeta,dirmeta }

	jsons, err := json.MarshalIndent(metadata, "", "   ")
	if err != nil {
		panic(err)
	}	
	fmt.Println(string(jsons))
}

func generateFilesMaps(fList *FilesMetaList) []*map[string]string {
	var filesHashMap map[string]string = make(map[string]string)
	var fielsSizeMap map[string]string = make(map[string]string)
	var directoriesSizeMap map[string]string = make(map[string]string)
	var result []*map[string]string
	filespath := (*fList).fileNames
	directoriespath := (*fList).directoryNames

	for indx,filehash := range (*fList).fileHashes{
		filesHashMap[filespath[indx]] = string(filehash)
	}

	result = append(result, &fielsSizeMap)

	for indx,filesize := range (*fList).fileSizes{
		fielsSizeMap[filespath[indx]] = strconv.Itoa(filesize)
	}
	result = append(result, &filesHashMap)

	for indx,dirSize := range (*fList).diretorySizes{
		directoriesSizeMap[directoriespath[indx]] = strconv.Itoa(dirSize)
	}
	result = append(result, &directoriesSizeMap)

	return result
}

func getManifestBody(fList *FilesMetaList) *ManifestDirectoryBody{

	var result ManifestDirectoryBody
	currentdir,_ := os.Getwd()

	for indx,fname := range (*fList).fileNames{
		fsize := (*fList).fileSizes[indx]
		fhash := (*fList).fileHashes[indx]
		fullname := currentdir + "/" + fname
		paths, fileName := filepath.Split(fullname)
		manifestFile := ManifestFile {
			Path: []byte(paths) ,
			FileName: []byte(fileName),
			Size: int64(fsize),
			HashList: []HashType {{[]byte("base64,sha256"), fhash}},
			MetaString: []byte{},
		}
		result.FileList = append(result.FileList, manifestFile)
	}

	for indx,dirname := range (*fList).directoryNames{
		dirsize := (*fList).diretorySizes[indx]
		fullDirname := currentdir + "/" + dirname
		manifestFile := ManifestFile {
			Path: []byte(fullDirname) ,
			FileName: nil,
			Size: int64(dirsize),
			HashList: []HashType {{[]byte("base64,sha256"), nil}},
			MetaString: []byte{},
		}
		result.FileList = append(result.FileList, manifestFile)
	}

	return &result
}

func getManifestDirectoryHeader(serializedkvList *SerializedKvList, body *ManifestDirectoryBody) *ManifestDirectoryHeader {
	var result ManifestDirectoryHeader

	segLenth := len((*body).FileList)
	version := []byte("1.0.0")
	sequenceid := uint64(1)
	createat := time.Now() 
	bodySegmentLength := uint64(segLenth)  
	bodyDataFileSize := uint64(3)
	serializedMapList := serializedkvList

	result = ManifestDirectoryHeader{
		VersionString: version,
		SequenceId: sequenceid,
		CreatedAt: createat,
		BodySegmentLength: bodySegmentLength,
		BodyDataFileSize: bodyDataFileSize,
		SerializedMapList: *serializedMapList,
	}

	return &result
}

func getManifestHeaderMetaData(header *ManifestDirectoryHeader) *ManifestHeaderMetaData {
	var result ManifestHeaderMetaData

	creationTime := time.Now().Unix()
	user, err := user.Current()
    if err != nil {
        panic(err)
    }

	previousManifest := getPreviousManifest()
	sequenceid := uint64(1)

	serializedheader := encoder.Serialize(*header)
	h := sha256.New()
	id := base64.StdEncoding.EncodeToString(h.Sum(serializedheader))

	result = ManifestHeaderMetaData {
		CreationTime: creationTime,
		Creator: user.Name,
		PreviousManifest: previousManifest,
		SequenceId: sequenceid,
		UniqueId: id,
	}
	return &result
}
 
func createCXOdir() {
	folderName := ".cxo"
	fmt.Println("Create .cxo foler in current directory: ")
	os.Mkdir(folderName, 0777)
	os.Chmod(folderName, 0777)
}

func serializeMaps(mapList []*map[string]string) *SerializedKvList{
	var serializedkvList SerializedKvList
	var kvList KeyValueList

	for _,mapPointer := range mapList{
		for key,value := range (*mapPointer){
			kvList = append(kvList, KeyValueString{key,value})
		}
	}

	sort.Sort(kvList)
	for _,pair := range kvList{
		serializedkvList.Add(KeyValueByte{[]byte(pair.Key), []byte(pair.Value)})
	}
	return &serializedkvList
}

func hashFileAndEncoding(filePath string) string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        log.Fatal(err)
    }

    return  base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func generateCXOFilename() string {
	dir,_ := os.Getwd()
	cxofilen := dir + "/.cxo/checkpoints/" + strconv.FormatInt(time.Now().Unix(),10)  + ".cxo"
    return cxofilen
}

func isCXOFolderExist() bool {
	dir,_ := os.Getwd()
	path := dir + "/.cxo/"
    _, err := os.Stat(path)
    if err != nil{
        if os.IsExist(err){
            return true
        }
        if os.IsNotExist(err){
            return false
        }
        fmt.Println(err)
        return false
    }
    return true
}

func getPreviousManifest() string {
	return "14387988.cxo"
}