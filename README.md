UDPServer
=========
UDPServer it's a asynchronous udp server & client     
(c) 2014 Cergoo    
under terms of ISC license    

<pre>
srv := UDPServerNew(...)    - create and run server   
pkg = &lt;-srv.ChRead       - get package from client
srv.ChWrite&lt;-pkg         - send package to client
srv = nil                   - stop and destroy server
</pre>

<pre>
client := Connect(...) 		  - create and run client
pkg = <-client.ChRead       - get package from server
client.ChWrite<-pkg         - send package to server
client = nil                - stop and destroy client
<pre>