package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func get_torrent_size(t string)float64{
	mi,err := metainfo.LoadFromFile(t)
	if err != nil{
		fmt.Println(err)
		return -1
	}
	

	info,err := mi.UnmarshalInfo()
	if err != nil{
		fmt.Println(err)
		return -1
	}



	var size float64 =0

	for _,fi := range info.Files{
		size+=float64(fi.Length)
	}
	return(float64(size/1073741824))
}



type Tor struct  {
name string
size float64
}

type file_details struct{
	path string
	size float64
	hash string
}

func print_torrent_details(p string){
	var fd file_details
	fd.path = p
	fd.size = get_torrent_size(p)
	mi,_ := metainfo.LoadFromFile(p)
	fd.hash = mi.HashInfoBytes().String()
	fmt.Printf("%s,%f,%s\n",fd.path,fd.size*1073741824,fd.hash)
}

func print_torrent_files_details(p string){
	c, _ := torrent.NewClient(nil)
	// defer c.Close()
	t,_ :=c.AddTorrentFromFile(p)
	t.GotInfo()
	for _,file := range t.Files(){
		fmt.Println(file.Path(),file.Length())
	}
}

func prepare_files_list(input_path string)[]Tor{
	var files_list []Tor
	f, err := os.Open(input_path)
    if err != nil {
        log.Fatal(err)
    }
    files, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
		files_list = append(files_list,Tor{name:file.Name(),size:get_torrent_size(input_path+"//"+file.Name())})
    }
	return files_list
}

func split_files(files []Tor,divisor,limit float64)[][]string{
	var size float64
	var temp []string
	var folders_list [][]string
	for _,file := range files{
		if((size+file.size)/divisor>limit){
			folders_list = append(folders_list, temp)
			size = 0.00
			temp = nil
		}else {
			size+=file.size
			temp = append(temp, file.name)
		}
	}
	folders_list = append(folders_list, temp)
	return folders_list
}

func Copy(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }
    return out.Close()
}

func copy_files(folders_list [][]string,input_path,output_path,limit string){
	for i,folder := range folders_list{
		var current_dir  = output_path+"/00"+strconv.Itoa(i)+"_"+limit
		err := os.Mkdir(current_dir, 0755)
			if err != nil {
				log.Fatal(err)
			}
			for _,file := range folder{
				Copy(input_path+"/"+file,current_dir+"/"+file)
			}
			
		}
}


func main(){
// File Path example  C:\\Users\\eee\\Desktop\\skycoin-services\\torrent-tool-go\\torrents\\r_2477000.torrent
// Directory path example C:\\Users\\eee\\Desktop\\torrent-sorting-go\\output 
var input_path string
var output_path string
var folder_limit string
var sort_type string
var help string
var sub_directory string
var files_info_directory string
flag.StringVar(&help, "help", "","Show help ")
flag.StringVar(&help, "h","", "Show help-shorthand")
flag.StringVar(&sub_directory,"sub-directory","","Path to torrent file to get it's details ")
flag.StringVar(&sub_directory,"sb","","Path to torrent file to get it's details-shorthand ")
flag.StringVar(&files_info_directory,"files-info","","Path to single torrent file to get each of it's files details")
flag.StringVar(&files_info_directory,"fi","","Path to single torrent file to get each of it's files details-shorthand")
flag.StringVar(&input_path, "directory", "","Path to the input directory")
flag.StringVar(&input_path, "d","", "Path to the input directory-shorthand")
flag.StringVar(&output_path, "output", "","Path to the output directory")
flag.StringVar(&output_path, "o","", "Path to the output directory-shorthand")
flag.StringVar(&folder_limit, "size", "1TB","Size limit in NGB or NTB ")
flag.StringVar(&folder_limit, "s","1TB", "Size limit in NGB or NTB-shorthand")
flag.StringVar(&sort_type, "type", "1","Sort type, 1 sorts by size 2 sorts by name")
flag.StringVar(&sort_type, "t","1", "Sort type, 1 sorts by size 2 sorts by name-shorthand")

flag.Parse()

if (sub_directory != ""){
	print_torrent_details(sub_directory)
	os.Exit(3)
}

if(files_info_directory != ""){
	print_torrent_files_details(files_info_directory)
	os.Exit(3)
}

var divisor = 1
if (strings.Contains(strings.ToLower(folder_limit),"tb")){
	divisor = 1000
}
re := regexp.MustCompile("[0-9]+")
var limit = re.FindAllString(folder_limit,-1)
file_limit,_ := strconv.ParseFloat(limit[0],64)

fmt.Println("Reading torrent files...")
var files []Tor = prepare_files_list(input_path)

fmt.Println("Sorting torrent files...")
if (sort_type == "1"){
	sort.Slice(files,func(i,j int)bool{
		return files[i].size > files[j].size
	})
} else {
sort.Slice(files,func(i,j int)bool{
	return files[i].name < files[j].name
})
}

fmt.Println("Copying torrent files into separate new directories...")
var folders_list = split_files(files,float64(divisor),file_limit)



copy_files(folders_list,input_path,output_path,folder_limit)


}