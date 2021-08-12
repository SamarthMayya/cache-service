# Service Caching

This project is an implementation of caching using golang and redis-server. 
Below are the steps involved to set up this project locally, assuming that
your machine has protobuf, go, and redis set up:
1. First, make sure that you have all the packages required, by running 
   ```
   go mod tidy
   ```
2. In order to ensure that the server can set and fetch keys from the database, 
   redis-server has to be started. Command to start the redis server with the required configuration is present 
   inside the init directory. To start up redis-server, make sure that you are in 
   the cache-service directory itself, and then run
   ```
   source init/initialize_redis.sh
   ```
3. This project has two servers and two clients (both within the client folder). 
   One is a user client, and another is an experimental one to test redis. 
   Based upon your choice, say `x.go`, run
   ```
   go run server/handlers/x.go
   ```
4. Now run the client of your choice, say `y.go`, as
   ```
   go run client/y.go
   ```
> The client should be running properly now. 

Some of the configurations of the system used while developing this project are as follows:
* OS : Ubuntu 20.04.2
* Repository path in local: `$HOME/cache-service`
* `GOPATH`: `$HOME/samarth/go` 