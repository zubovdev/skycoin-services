# Network info

You may use function `netinfo.Get(nil)` to get network information.  
You may pass parameter `netinfo.Get(&netinfo.HttpTest{Port: 8888})` - 
it will run additional HTTP tests over IPv4 and IPv6 addresses.