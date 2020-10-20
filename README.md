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

GeoLite2-City.mmdb file needed!

## GraphQL
GraphQL schema is important, it defines how results will be sent, data structure.

## Google cloud storage
Upload your account key! Name: QuicPos-key.json

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
none

## API
Save post
```graphql
mutation create {
  createPost(
    input: {
      text: "My new post"
      userId: "kuba"
      image: "base64-string"
    }
  ) {
    ID
    text
    userId
    shares
    views {
      userId
      time
    }
    creationTime
    initialReview
    image
    reports
  }
}
```

---
Get post
```graphql
query {
  post(userId: "ee", normalMode: true) {
    ID
    text
    userId
    shares
    views
    creationTime
    initialReview
    image
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

---
Get post by id
```graphql
query {
  viewerPost(id: "5f79cc689ec125d75f2e36e5") {
    ID
    text
    userId
    shares
    views
    creationTime
    initialReview
    image
  }
}
```
---
Get the oldest without initial review
```graphql
query {
  unReviewed(password: "funia", new: true) {
    post {
      ID
      text
      userId
      shares
      views
      creationTime
      initialReview
      image
    }
    left
  }
}
```

---
Get the most reported post
```graphql
query {
  unReviewed(password: "funia", new: false) {
    post {
      ID
      text
      userId
      shares
      views
      creationTime
      initialReview
      image
    }
    left
  }
}
```

---
Initial review without delete
```graphql
mutation review {
  review(
    input: {
      postID: "5f79cc2f9ec125d75f2e36e4"
      new: true
      delete: false
      password: "funia"
    }
  )
}
```


---
Reported review with delete
```graphql
mutation review {
  review(
    input: {
      postID: "5f79cc2f9ec125d75f2e36e4"
      new: false
      delete: true
      password: "funia"
    }
  )
}
```


---
New view
```graphql
mutation view {
  view(input: { postID: "5f79cc2f9ec125d75f2e36e4", userId: "kuba", time: 1.0, deviceDetails: "iPhone" })
}

```


---
New report
```graphql
mutation review {
  report(input: { userID: "kuba2", postID: "5f79cc2f9ec125d75f2e36e4" })
}

```


---
New share
```graphql
mutation share {
  share(input: {userID: "kuba", postID: "5f79cc2f9ec125d75f2e36e4"})
}
```