package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/urfave/cli/v2"
)

func initCLI() *cli.App {
	filesList = processDirAndGenerateMeta(".")

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
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "print-version",
		Usage: "print version",
	}

	return app
}

func addCLICommands(app *cli.App) {
	app.Commands = []*cli.Command{
		{
			Name:      "init",
			Usage:     "initialize tool environment by create the .cxo folder",
			UsageText: "create the manifest folder .cxo in current directory",
			Action: func(cnx *cli.Context) error {
				err := createFolder(".cxo")
				if err == nil {
					fmt.Println("Create .cxo foler in current directory: ")
				}
				return err
			},
		},
		{
			Name:      "commit",
			Usage:     "commit all the files' metadatum into the .cxo file",
			UsageText: "commit all the metadata files into the .cxo folder",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "print-json",
					Value: false,
					Usage: "print files in the directory in json ",
				},
				&cli.BoolFlag{
					Name:  "meta",
					Value: false,
					Usage: "add metadata section in json",
				},
			},
			Action: func(cnx *cli.Context) error {
				metaFlag := false
				if cnx.Bool("meta") {
					if !cnx.Bool("print-json") {
						cli.ShowAppHelpAndExit(cnx, 0)
					}
					metaFlag = true
				}
				cxoPath := currentDir + "/.cxo/"
				if !isFolderExist(cxoPath) {
					fmt.Println("please use 'manifest init' command before 'manifest commit'")
					os.Exit(1)
				}

				err := os.MkdirAll("./.cxo/checkpoints/", os.ModePerm)
				if err != nil {
					panic(err)
				}
				cxoFileName := currentDir + "/.cxo/checkpoints/" + strconv.FormatInt(time.Now().Unix(), 10) + ".cxo"

				cxoFile, err := os.OpenFile(cxoFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					panic(err)
				}
				defer cxoFile.Close()

				var manifestOuputBody ManifestOuputBody
				manifestBody := getManifestBody(filesList)
				manifestOuputBody.ManifestBody = *manifestBody
				manifestOuputBody.ManifestHeader = *getManifestDirectoryHeader(manifestBody)
				manifestOuputBody.ChunkHashList = (*filesList).fileschunkslist

				serializedOuputBody := encoder.Serialize(manifestOuputBody)

				_, err = cxoFile.Write(serializedOuputBody)
				if err != nil {
					panic(err)
				}
				if cnx.Bool("print-json") {
					printFilesInJson(filesList, &manifestOuputBody.ManifestHeader, metaFlag)
				}

				return nil
			},
		},
	}
}

func main() {

	currentDir = getCurrentDir()
	app := initCLI()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func processDirAndGenerateMeta(dir string) *FilesInfoList {
	var FilesAndDirectories FilesInfoList
	var directories []string
	var directoriesSize []int
	var files []string
	var filesSize []int
	var filesHash [][]byte
	var fileschunksList []FileChunkHashList
	var filesMetaList ManifestDirectMetaList
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if info.Name() == ".cxo" {
					return filepath.SkipDir
				}
				directories = append(directories, path)
				dirSize, err := getDirectorySize(path)
				if err != nil {
					return err
				}
				directoriesSize = append(directoriesSize, dirSize)
			} else if info.Name() != appName {
				files = append(files, path)
				filesSize = append(filesSize, int(info.Size()))
				filesHash = append(filesHash, []byte(hashFileAndEncoding(path)))
				chunkshashes, err := getChunkHashes(path)
				if err != nil {
					return err
				}
				fileschunksList = append(fileschunksList, *chunkshashes)

				filesMetaList = append(filesMetaList, getFileMeta(path))
			}

			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	FilesAndDirectories.directoryNames = directories
	FilesAndDirectories.fileNames = files
	FilesAndDirectories.fileSizes = filesSize
	FilesAndDirectories.fileHashes = filesHash
	FilesAndDirectories.diretorySizes = directoriesSize
	FilesAndDirectories.fileschunkslist = fileschunksList
	FilesAndDirectories.filesMetaList = filesMetaList
	return &FilesAndDirectories
}

func getChunkHashes(filepath string) (*FileChunkHashList, error) {

	var filechunks FileChunkHashList

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bf := make([]byte, chunkSize)
	hs := sha256.New()

	for {
		readTotal, err := file.Read(bf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if readTotal == 0 {
			break
		}

		for readTotal < chunkSize {
			readTotal = readTotal + copy(bf[readTotal:], []byte{0x0000})
		}

		hs.Write(bf)
		filechunks.ChunksHashes = append(filechunks.ChunksHashes, hs.Sum(nil))
		hs.Reset()

	}

	return &filechunks, nil
}

func getDirectorySize(directory string) (int, error) {
	totalSize := 0
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				totalSize = +int(info.Size())

			}
			return nil
		})
	return totalSize, err
}

func printFilesInJson(fList *FilesInfoList, dirHeader *ManifestDirectoryHeader, metaflag bool) {
	var dirmeta DirectoryMetaList
	var filemeta FileDataList

	for indx, fn := range (*fList).fileNames {
		fh := (*fList).fileHashes[indx]
		fs := (*fList).fileSizes[indx]
		meta := (*fList).filesMetaList[indx]
		fileInfo := FileData{fn, fs, fh, &meta}
		if !metaflag {
			fileInfo = FileData{fn, fs, fh, nil}
		}

		filemeta = append(filemeta, fileInfo)

	}

	for indx, dn := range (*fList).directoryNames {
		ds := (*fList).diretorySizes[indx]
		dirInfo := DirectoryMeta{dn, ds}
		dirmeta = append(dirmeta, dirInfo)
	}

	sort.Sort(filemeta)
	sort.Sort(dirmeta)
	metadata := struct {
		DirectoryHeader ManifestDirectoryHeader `json:"directory header"`
		Directories     []DirectoryMeta         `json:"directories"`
		Files           []FileData              `json:"files"`
	}{*dirHeader, dirmeta, filemeta}

	jsons, err := json.MarshalIndent(metadata, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsons))
}

func getManifestBody(fList *FilesInfoList) *ManifestDirectoryBody {

	var result ManifestDirectoryBody

	for indx, fname := range (*fList).fileNames {
		fsize := (*fList).fileSizes[indx]
		fhash := (*fList).fileHashes[indx]
		fullname := currentDir + "/" + fname
		paths, fileName := filepath.Split(fullname)
		manifestFile := ManifestFile{
			Path:       []byte(paths),
			FileName:   []byte(fileName),
			Size:       int64(fsize),
			HashList:   HashValue{[]byte("base64,sha256"), fhash},
			MetaString: []byte{},
		}
		result.FileList = append(result.FileList, manifestFile)
	}

	for indx, dirname := range (*fList).directoryNames {
		dirsize := (*fList).diretorySizes[indx]
		fullDirname := currentDir + "/" + dirname
		manifestFile := ManifestFile{
			Path:       []byte(fullDirname),
			FileName:   nil,
			Size:       int64(dirsize),
			HashList:   HashValue{[]byte("base64,sha256"), nil},
			MetaString: []byte{},
		}
		result.FileList = append(result.FileList, manifestFile)
	}

	return &result
}

func getManifestDirectoryHeader(body *ManifestDirectoryBody) *ManifestDirectoryHeader {
	var result ManifestDirectoryHeader
	dataSize := 0

	for _, manifile := range (*body).FileList {
		if manifile.FileName != nil {
			dataSize += int(manifile.Size)
		}
	}
	segLenth := unsafe.Sizeof((*body).FileList)
	version := []byte(versionNo)
	sequenceid := uint64(getSequenceId())
	createat := uint64(time.Now().Unix())
	bodySegmentLength := uint64(segLenth)
	bodyDataFileSize := uint64(dataSize)

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	result = ManifestDirectoryHeader{
		VersionString:     version,
		SequenceId:        sequenceid,
		Creator:           user.Name,
		CreatedAt:         createat,
		BodySegmentLength: bodySegmentLength,
		BodyDataFileSize:  bodyDataFileSize,
		MetaDataTags:      KeysValuesList{},
		ChunkSize:         chunkSize,
	}

	return &result
}

func getManifestHeaderMetaData(header *ManifestDirectoryHeader) *ManifestHeaderMetaData {
	var result ManifestHeaderMetaData

	creationTime := (*header).CreatedAt

	filename, err := getPreviousManifest((*header).SequenceId)
	if err != nil {
		panic(err)
	}
	previousManifest := filename

	serializedheader := encoder.Serialize(*header)
	h := sha256.New()
	h.Write(serializedheader)
	id := base64.StdEncoding.EncodeToString(h.Sum(nil))

	result = ManifestHeaderMetaData{
		CreationTime:     creationTime,
		Creator:          (*header).Creator,
		PreviousManifest: previousManifest,
		SequenceId:       (*header).SequenceId,
		UniqueId:         id,
	}
	return &result
}

func createFolder(folderName string) error {
	var err error
	if isFolderExist(folderName) {
		return nil
	}
	err = os.Mkdir(folderName, 0777)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(folderName, 0777)
	if err != nil {
		log.Fatal(err)
	}
	return err
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

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func isFolderExist(path string) bool {

	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func getPreviousManifest(currentSequenctId uint64) (string, error) {
	var manifestOuputBody ManifestOuputBody
	cxoFolderName := currentDir + "/.cxo/checkpoints/"
	var filename string
	files, _ := ioutil.ReadDir(cxoFolderName)
	for _, file := range files {
		filename = file.Name()
		if strings.HasSuffix(filename, ".cxo") {
			cxoFile, err := os.Open(cxoFolderName + filename)
			if err != nil {
				return "", err
			}
			defer cxoFile.Close()
			fileBytes, err := ioutil.ReadAll(cxoFile)
			if err != nil {
				return "", err
			}
			_, err = encoder.DeserializeRaw(fileBytes, &manifestOuputBody)
			if err != nil {
				return "", err
			}
			if currentSequenctId == manifestOuputBody.ManifestHeader.SequenceId+1 {
				break
			}
		}
	}
	return filename, nil
}

func getSequenceId() int {
	cxoFolderName := currentDir + "/.cxo/checkpoints/"
	files, _ := ioutil.ReadDir(cxoFolderName)
	count := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".cxo") {
			count++
		}
	}
	return count
}

func getCurrentDir() string {
	dir, _ := os.Getwd()
	return dir
}

func getFileMeta(filename string) FileMeta {

	var result FileMeta

	fileInfo, _ := os.Stat(filename)
	fStat := fileInfo.Sys().(*syscall.Stat_t)
	result.LastModified = uint64(fileInfo.ModTime().Unix())
	result.UnixPermission = fileInfo.Mode().String()
	sec, _ := fStat.Ctim.Unix()
	result.CreateAt = uint64(sec)

	return result
}
