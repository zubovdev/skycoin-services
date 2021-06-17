# System survey

System-survey tool an CLI application which allows executing specific tests and collect information about the host
system.

The application contains from two parts: collect information and testing.

There are several command to collect specific type of information.

## Requirements
- _Linux OS for production usage_.
- Golang version `>=1.16`.


## Usage

You can use the functions in `cmd` and `tests` package in your Go code, and also execute the following functions throw
the binary file.

- [Application info](#Application info)
- [Golang version](#Golang version)

### Application info

Information about the applications installed and which indicated in the `PATH` environment variable can be retrieved
by `cmd.NewAppList(filter)` function, where filter is type of `string` and could be empty, or a comma separated values
string.

The `cmd.AppList` type is a `map[string][]string` with following methods available:

- `JSON()` returns `[]byte` represented marshalled `AppList` to the **JSON** format.
- `String()` returns human-readable and formatted string.

```go
package main

import (
	"fmt"
	"github.com/skycoin/skycoin-services/system-survey/cmd"
)

func main() {
	// Filter is empty, means return all available applications.
	apps := cmd.GetAppList("")
	// Print the result in a human-readable format.
	fmt.Println(apps.String())
}
```

Example output from my personal workstation (human-readable):

```
/usr/local/bin:
    - cat
    - curl
    - wget
...
```

Example output from my personal workstation (machine-readable):

```
{"/usr/local/bin":["cat","curl","wget"]}
```

#### CLI usage

You just need to run the `./system-survey apps` and it will print human-readable result to the console.  
This CLI command have two flags additionally

- `--json` - for representing the output in machine-readable format.
- `--filter` - comma-separated values for filtering the output.

Full example: `./system-survey apps --json --filter=ls,cat,get,curl`.  
Refer to `./system-survey apps --help` for more information about the CLI usage.

### Golang version

Information about the Golang version installed can be retrieved by `cmd.NewGolangVersion()`,
which returns pointer to a`golangVersion` or `error`.

The `cmd.golangVersion` struct has two methods:

- `JSON()` returns `[]byte` represented marshalled `AppList` to the **JSON** format.
- `String()` returns human-readable and formatted string.

```go
package main

import (
	"fmt"
	"github.com/skycoin/skycoin-services/system-survey/cmd"
	"log"
)

func main() {
	gv, err := cmd.GetGolangVersion()
	if err != nil {
		log.Fatalf("Failed to get golang version: %v", err)
	}
	// Print the result in a human-readable format.
	fmt.Println(gv.String())
}
```

Example output from my personal workstation (human-readable):

```
version=1.16.4, os=windows, arch=amd64
```

Example output from my personal workstation (machine-readable):

```
{"version":"1.16.4","os":"windows","arch":"amd64"}
```

#### CLI Usage

You just need to run the `./system-survey golang` and it will print human-readable result to the console.  
This CLI command have two flags additionally

- `--json` - for representing the output in machine-readable format.

Full example: `./system-survey golang`.  
Refer to `./system-survey golang --help` for more information about the CLI usage.