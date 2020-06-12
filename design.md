## HIGH LEVEL DESIGN

The communication between servents are almost similar to GNUTELLA protocol mentioned [here](http://rfc-gnutella.sourceforge.net/developer/stable/index.html)

The servent uses four descriptors:

> |descriptor |description                                                                            |
> |:----------|:--------------------------------------------------------------------------------------|
> |ping       |Used to actively discover hosts on the network. A servent receiving a Ping descriptor responds with pong     descriptor|
> |pong       |The response to a Ping. Includes the address of a connected servent and information regarding the amount of files and total size it is making available to the network|
> |query      |The primary mechanism for searching the distributed network. A servent receiving a Query descriptor will respond with a QueryHit if a match is found against its files in current working directory|
> |queryHit   |The response to a Query. This descriptor provides the recipient with enough information like http port to acquire the file matching the corresponding Query.|

All the four descriptor data are encoded and decoded in Little Endian format with the exception of IP address which is in Big Endian format as mentioned in the GNUTELLA protocol. 

Default servent Ping Interval and waiting time for query results are set as 5 seconds.
Descriptor TTL is set at 5 i.e. after 5 hops, the descriptor will not be forwarded.

There will be id's generated for each connection and each file.
`open` accepts connection string and gives connection_id
`close` accepts connection_id
`find` accepts file name and gives file_id
`get` accepts file_id and downloads file

If a particular connection is closed, no descriptor will be sent to that servent.
The connection has to be opened again to resume communication.
If a particular servent goes down, the other servents remove the connection from their system.

### LIMITATIONS:

Currently works with servents in local ports.
Exact file name search is only possible.
Servent file system comprises of files (excluding nested directories) in the current working directory.