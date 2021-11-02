# ChittyChat
Miniproject 2 


run the following commands to run the server and clients 
in different terminals run:
```make serv```

```make cli```

```make cli1```

```make cli2```

```make cli3```


to remake the proto generated files run 
    ```make gen```


## Task list 
  - done!!  Set up gRPC 
  - done!! Set up server 
  - done!! Set up client(s)
  - Lamport time
    - done!! add lamport time to messages 
    - Done!! update lamport time in clients (and server)
    - Nope! maybe vector clock - not sure how to do the with increasing/decreasing number of clients
 - Done!! need to make it possible to have several users named anon. 
  - Done! -added unique identifierchange the way the map works. or add an unique identifier to the user in the message. 
    - Nope! This will probably need a different implimentation of the streams. 
    suggested flow of streams 
        - rpc publish(User) retruns (sream message)
        - rpc broadcast(stream Message) retruns (Empty)
        well maybe not.. i don't know.

