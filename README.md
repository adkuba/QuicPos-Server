# QuicPos
Application backend in golang


## Golang notes
Simple installation - only remember about path variable. I needed to edit <code>/etc/environment</code>

Example package download <code>go get go.mongodb.org/mongo-driver/mongo</code>

Important! Go to file -> preferences -> settings -> go (extension) -> change format tool to "gofmt"

Workflow:
- in <code>schema.graphqls</code> define your models and operations - mutations and queries
- generate functions with <code>go run github.com/99designs/gqlgen generate</code>
- go to <code>schema.resolvers.go</code> and implement functions.

## GraphQL
GraphQL schema is important, it defines how results will be sent, data structure.


## Steps to recreate project
In this directory run <code>go mod init QuicPos</code><br>
Im using graphQL package <br>
Then <code>go run github.com/99designs/gqlgen init</code> creates project structure.<br>
Delete CreateTodos and Todos from schema.resolvers.go.

## Mongodb
In internal/mongodb is file with all mongo database scripts.

## TODO
Integrate mongodb with schema!