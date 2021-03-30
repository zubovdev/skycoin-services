


# Tasks

Use library "flags" or urfave/cli
- https://gobyexample.com/command-line-flags
- https://github.com/urfave/cli

Put each command in /cmd/ sub-folder
- human readable result is fine
- later we want result filled in to machine readable format (like JSON)

## Run on System Startup

Figure out how to have executable add itself to system startup
- for debian

## Read out Path Variable and check for list of applications

- wget
- golang
- git

## Run "go version" and get string output or failure (if executable exists)

1> Is run "go version"
- get output

2> Check if golang is in path

## Use golang to get information about computer

3> Get CPU architecture (ARM, x64, etc)
- and operating system
- operating version

4> Get computer ram

5> Get list of network interfaces

6> Get list of discs and disc size


##Golang TraceRoute utility

Use one of the following 
- https://github.com/kalbhor/tracesite
- https://pkg.go.dev/github.com/aeden/traceroute?utm_source=godoc (github.com/aeden/traceroute)

Check if traceroute is using
- ICMP, UDP, TCP (verify that it is using UDP)

Get results of traceroute back in JSON or machine readable format
- time in ms
- ip for each host
- the hop count
- etc

The parameters must be fed in as a struct (startin host, destination host, hop count, max latency, protocol=UDP, etc)
- the results should be stored in a "results struct"
- there should be a machine readable (JSON) print out of the results struct
- the results struct should include unix, int64 time in seconds, when test was performed
- there should be field for "source ip"

So we have 

```golang
struct TraceRouteInput {
	SourceHostIP String //the host that made request (can ignore)
	SourceHostPort int
	DestinationIp string
	DesinationPort string
	
	Time int64 //int64 unix time
	Protocol String //UDP
	MaxLatency int //
	MaxHops int
	Retries int
	//etc
}

struct TraceRouteResult {
	//etc, return data goes here
}

use github.com/skycoin/skycoin, src/cipher/encoder for serializing the input and result structs to []byte

Also have a  function for serializing the results to JSON

```

## Use Golang to get network information

7> Get public ip address
8> Get local network ip address
9> Test whether IP address is publicly accessible
10> HTTPS test 

## Advanced networking test


10> Do DHT test to determine if computer can do bitorrent DHT

11> Do Connection test for firewall hole punching

12> Do CX tests



14> Ping tests
- latency test

15> Data transfer test against test servers

16> DMSG Test

17> Transport Test

18> Trace Route

19> Log upload test