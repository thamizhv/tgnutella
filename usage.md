### EXAMPLE USAGE

Let's assume the below folder structure.

folder_1:
    a.txt
folder_2:
    b.txt
folder_3:
    c.txt


Only the current working directory files (excluding inner directories) are taken into account.
And we are running tgnutella in these three folders separately by below commands.

folder_1 -> tgnutella 8100
folder_2 -> tgnutella 8200
folder_3 -> tgnutella 8300

TCP servers will be started in mentioned ports and http servers will be started in +1 ports.
Eg: TCP running in 8100 and HTTP running in 8101 ports


Give the below commands in folder_1 session.
> |commands               |action                                                                                 |
> |:---------------       |:--------------------------------------------------------------------------------------|
> |open :8200             |Establishes TCP connection to 8200 server and returns connection id.                   |
> |open :8300             |Establishes TCP connection to 8300 server and returns connection id.                   |
> |info connections       |List all TCP connections                                                               |
> |find c.txt             |search for a file in the network. If found within 5 seconds, an id will appear. If the file is present in network, id will appear when it is found by current server.                                            |
> |get <file_id>          |Makes a TCP connection to 8300 server and downloads c.txt file and saves it in folder_1|  
> |close <connection_id>  |closes the connection to the corresponding server.                                     |