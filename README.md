Oh man.... this code is messy and fragile. If I ever become an Italian chef I could probably put a picture of this code on the menu as the spaghetti.

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

##### Distributed Server
###### Proposed Server - Local Machine
This test case does not outperform the traditional machine locally. It might be a bottleneck with my machine though

```go test -run Internal_Distributed```
