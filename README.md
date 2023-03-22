# Challenge GoExpert

## Objective

Create an application client/server to get dollar quotation.

### How to execute

#### Requirements

- Install dependencies of app "Server"

```shell
cd ./server && go mod tidy
```

- Install dependencies of app "Client"

```shell
cd ./client && go mod tidy
```

#### Server

- Start server with the following command

```shell
cd ./server && go run main.go
```

#### Client

- Run client with the following command

```shell
cd ./client && go run main.go
```

##### Expected behavior

- Should exists `./client/Quotation.txt` file
  - To see file content run `cat ./client/Quotation.txt`
- Should exists `./server/database.s3db` file
  - To see database records can use VSCode extension called [SQLite Viewer](https://marketplace.visualstudio.com/items?itemName=qwtel.sqlite-viewer)
