UDPServer
=========
UDPServer asynchronous udp server    
(c) 2014 Cergoo    
under terms of ISC license    

<pre>
srv := UDPServerNew(...) - create and run server   
pkg = <-srv.ChRead       - get package from client
srv.ChWrite<-pkg         - send package to client
srv = nil                - stop and destroy server
</pre>
