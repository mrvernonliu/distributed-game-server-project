# Distributed game server for CSC462
The idea of this project is to use a distributed system to increase the tick rate for games that utilize a large amount of players.

### Instructions
##### Traditional Server - Local Machine
Simply running the test case will give you all the statistics you need
```go test -run Internal_Traditional```
##### Traditional Server - External Machine
In this case we are going to run the server on a different machine. To do this simply copy the files over to that machine:

```scp -r distributed-game-server-project/ vernon@10.0.0.55:```

And then start only the server using the following command:

```go run main/start_traditional_server.go```

Make sure you change the ip in your test case to match your server ip and then run the following command from your local machine

```go test -run External_Traditional```


##### Distributed Server - Local Machine
This method of running the distributed server is just meant to make sure the code works before moving the servers to different clusters.
In this case just run:

```go test -run Internal_Distributed```

##### Distributed Server - External Machines
In this case we are going to use an external machine for the server and an external machine for all of the workers on top of our client machine.

First you will need to make changes to the starting files and the test files to properly connect to the IP address of the machines

Then copy the code to each of the respective machines

```scp -P 2200 -r distributed-game-server-project/ vernon@192.168.0.18:```

```scp -r distributed-game-server-project/ vernon@192.168.0.20:```

Now, on your worker machine run:

```go run main/start_distributed_workers.go```

Then on your server machine run:

```go run main/start_distributed_server.go```

And finally on your local machine run the test script to start the players:

```go test -run External_Distributed```