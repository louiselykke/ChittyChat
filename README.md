# ChittyChat
Miniproject 2 



## Task list 
 - done!!  Set up gRPC following NaddiNadja guide https://github.com/NaddiNadja/grpc101 
  - done!! Set up server 
  - done!! Set clients up 
  - (no?)maybe use publish and broadcast methods? 
  - Lamport time
    - (done) add lamport time to messages 
    - Done!! update lamport time in clients (and server)
    - Nope! maybe vector clock - not sure how to do the with increasing/decreasing number of clients
 - need to make it possible to have several users named anon. 
  - change the way the map works. or add an unique identifier to the user in the message. This will probably 
    need a different implimentation of the streams. 
    suggested flow of streams 
        - rpc publish(User) retruns (sream message)
        - rpc broadcast(stream Message) retruns (Empty)
        well maybe not.. i don't know.

