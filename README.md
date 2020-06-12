 ## TGNUTELLA

 TGnutella is a minimalistic version of [Gnutella](https://en.wikipedia.org/wiki/Gnutella) [servent](https://en.wiktionary.org/wiki/servent). It provides a command line interface to connect with peers in the network, find and download files over the network.

 ### BASIC USAGE
Download binary `tgnutella` and run below command:

tgnutella <servent_port>  

Eg: `tgnutella 8100`
servent_port: the port where current servent will run and listen for command and descriptors


## Commands:
> |commands         | usage                                                                 |
> |:--------------- |:----------------------------------------------------------------------|
> |help             |lists details of available commands                                    |
> |open <host:port> |connects to a node on the network and returns its id                   |
> |close <id>       |closes connection by connection id (see info command)                  |
> |info connections |prints list of connected hosts with an id for each                     |
> |find <keyword>   |search files on the network and lists results with an id for each entry|
> |get <id>         |download a file by id                                                  |