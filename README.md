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
In internal/mongodb is file with all mongo database scripts. <br>
If mongo can't find DNS edit <code>/etc/resolv.conf</code> and change nameserver to 8.8.8.8 <br>
Important handling of ObjectID - see post.go. I don't know if it is good but works well.

## TODO
Rozpisac jak ma wygladac dokladnie api, komunikacja  z aplikacja i dokladnie zaimplementowac te endpointy, tylko bez uczenia maszynowego - czyli podrzuca randomowe posty a całą reszta normalnie!

## Examples
Save post
```graphql
mutation create {
  createPost(input: { text: "My new post", userId: "kuba" }) {
    ID
    text
    userId
    shares
    views{
      userId
      time
    }
  }
}
```

---
Get post
```graphql
query {
  post(userId: "ee") {
    ID
    text
    userId
    shares
    views{
      userId
      time
    }
  }
}
```

---
Get new UserID
```graphql
query {
  createUser
}
```