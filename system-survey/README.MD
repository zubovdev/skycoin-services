# Network info
You may use function `netinfo.Get()` to get network information.

# HTTP test server
It's located in `cmd/http`. You may pass `--addr=0.0.0.0:8888` 
flag to run the server on specific address, it's **0.0.0.0:8888** by default.
After server successfully started, then just on method `POST /ping` will be available.
It takes request body:
```
{
    "id": 432214 # Request ID
    "rand": 59189274899 # Random number
}
```

After accepting the request, it will increment `rand` parameter and send this
body back to the requester.