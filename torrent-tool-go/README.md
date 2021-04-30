# Torrent-tool

This is a simple script that sorts and splits torrent files in go, It separates torrents into size based directories

## Installation

Using this package requires a working Go environment. [See the install instructions for Go](http://golang.org/doc/install.html).

Go Modules are required when using this package. [See the go blog guide on using Go Modules](https://blog.golang.org/using-go-modules).

## Flags

To run the script you need to provide 4 flags or the help flag

### [-h,--help]

Shows help and all flags

### [-d,--directory]

Provide the full path to the input directory where the torrents are located

### [-o,--output]

Provide the full path to the output directory where the subdirectories and torrents will be copied

### [-s,--size]

The size limit the torrents will be split into in NGB or NTB format

### [-t,--type]

Sorting type of the torrents, 1 to sort by size and 2 to sort by name

### [-sb,--sub-directory]

Takes a path to a single torrent file and prints file details Path/Name, size is bytes and hash

## Usage

Run the script using go

```
$ go run script.go -d<value> -o<value> -s<value> -t <value>
or
$ go run script.go -directory<value> -output<value> -size<value> -type <value>
```

Or to use the single file details utility

```
$ go run script.go -sb<value>
or
$ go run script.go --sub-directory<value>
```

For example

```
$ go run script.go -d C:\\Users\\user\\Desktop\\torrents\\input -o C:\\Users\\user\\Desktop\\torrents\\output  -s 200GB -t 1
```

Takes all torrents in directory input, sorts them by size and splits them into 200GB folders in directory output

```
$ go run script.go -sb C:\\Users\\user\\Desktop\\torrents\\input\\tor1.torrent
```

Gets the torrent called tor1 and prints it's Path/Name , size in bytes and hash as a single line comma sparated values
