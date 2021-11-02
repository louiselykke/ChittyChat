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


# Docker
run command:

$ docker build -t chitty --no-cache .

to build the code. what you write after -t will be the name of your image, so the name of the image here is chitty. The name doesn't matter, but it helps you identify it in the docker desktop app - It has to be lowercase though.

run command:

$ docker run -p 9080:9080 -tid chitty

To run the code in a docker container. if you changed the name of the image from chitty, make sure to change it in this command as well.

Now clients can connect to the docker container for example by running the make commands from above.


## Task list 
  - done!!  Set up gRPC 
  - done!! Set up server 
  - done!! Set up client(s)
  - Lamport time
    - done!! add lamport time to messages 
    - Done!! update lamport time in clients (and server)
    - Nope! maybe vector clock - not sure how to do the with increasing/decreasing number of clients
 - Done!! need to make it possible to have several users named anon. 
  - Done! (added unique identifier) Change the way the map works. or add an unique identifier to the user in the message. 
    - Nope! This will probably need a different implimentation of the streams. 
    suggested flow of streams 
        - rpc publish(User) retruns (sream message)
        - rpc broadcast(stream Message) retruns (Empty)
        well maybe not.. i don't know.

