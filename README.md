# Table of Contents
- [What i've learned]()
- [Important notes](#important-notes)
- [Tensorflow](#tensorflow)
- [Golang notes](#golang-notes)


# Important notes
* Server only consumes jpeg images
* Nohup is better than disown for startig server in Linux [More detailed answer](https://unix.stackexchange.com/questions/3886/difference-between-nohup-disown-and)
* Needed files:
  * <code>geoloc/GeoLite2-City.mmdb</code> file with ip-localization database
  * <code>certificate.crt</code> and <code>private.key</code> file with SSL certificate
  * <code>QuicPos-key.json</code> with Google cloud storage key
  * <code>data/passwords.go</code> with 2 passwords, mongoSRV and Stripe private key
  * <code>out/</code> directory with 3x saved tensorflow models and 2x json dictionaries to work with this models


# Tensorflow
Tensorflow notes, how to use Tensorflow with golang.

### Installation
Initially I couldn't install Tensorflow. The reason was, that golang package links to official tensorflow repo, but doesn't support 2.0 version. So to install tensorflow I needed to execute:

```sh
go get github.com/tensorflow/tensorflow/tensorflow/go@v1.15.4
```

 Also I skipped naming layers when saving model to pb file in keras. To check default names execute this command:
 ```
 saved_model_cli show --dir out/ --all
 ```

Check fields names without ":0". 

### Additional steps that may help:
* Installing tensorflow C [tutorial](https://www.tensorflow.org/install/lang_c)
* Interesting [tutorial](https://tonytruong.net/running-a-keras-tensorflow-model-in-golang/)
* All values needs to be float32

### Net structure
To see detailed net structure go to [QuicPos-Microservice]() repository.



# Golang notes
Working with golang, important notes:
- Simple installation - only remember about the PATH variable. To change it I needed to edit <code>/etc/environment</code>

- Example go package download <code>go get go.mongodb.org/mongo-driver/mongo</code>

- Important! Go to file -> preferences -> settings -> go (extension) -> change format tool to "gofmt"

### GraphQL Workflow:
GraphQL schema is important, it defines how results will be sent, data structure.
- In <code>schema.graphqls</code> define your models and operations - mutations and queries.
- Generate functions with <code>go run github.com/99designs/gqlgen generate</code>
- Go to <code>schema.resolvers.go</code> and implement functions.



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
Important handling of ObjectID - see post.go. I don't know if it is good but works well. <br>
Localization stats:
```
[{
    $match: {
        "views": {
            "$ne": null
        }
    }
}, {
    $project: {
        "views.localization": 1,
        "_id": 0
    }
}, {
    $unwind: {
        path: "$views"
    }
}, {
    $group: {
        _id: "$views.localization",
        count: {
            $sum: 1
        }
    }
}]
```

Ciekawe query usuwające wszystkie linki w tekście posta.
```sh
db.posts.find({text: {$regex: "https:[^ ]+"}}).forEach(function(e, i){ 
    const regex = /https:[^ ]+/gi; 
    e.text = e.text.replace(regex, ''); 
    db.posts.save(e); 
})
```

## Deploy
You don't have to create docker image. Just build golang project, send to virtual machine and run. Then execute:
```sh
[1]+  Stopped                 myprogram
$ disown -h %1
$ bg 1
[1]+ myprogram &
$ logout
```

## Docker* napisanie statystyk na serwerze (moze być bez implementacji na stronie )
Make docker image:
* check if go builds correctly <code>go build -o bin/</code> and then run <code>bin/QuicPos</code>
* build image <code>docker build -t quicpos .</code>
* check if works <code>docker run -p 8080:8080 quicpos</code>
* export to file <code>docker save -o ./bin/quicpos.tar quicpos</code>
* transfer exported tar to your virtual machine <code>scp ./bin/quicpos.tar root@128.199.45.42:~/quicpos.tar</code>
* connect with virtual machine <code>ssh root@128.199.45.42</code>
* load tar <code>docker load -i quicpos.tar</code>
* run in background exposing ports <code>docker run -d --rm -p 80:8080 quicpos</code>
* check if is running <code>docker ps -a</code>
* stop <code>docker stop \<name\></code>

## API
Save post
```graphql
mutation create {
  createPost(
    input: {
      text: "My new post"
      userId: -1
      image: "base64-string"
    }
    password: "kuba"
  ) {
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
Get post
```graphql
query {
  post(userId: -1, normalMode: true, ad: false, password: "kubad") {
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
  createUser(password: "frw")
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
  unReviewed(password: "funia", type: 0) {
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
    spam
  }
}
```

---
Get the most reported post
```graphql
query {
  unReviewed(password: "funia", type: 1) {
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
    spam
  }
}
```

---
Get post without humanreview
```graphql
query {
  unReviewed(password: "funia", type: 2) {
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
    spam
  }
}
```


---
Block user by user
```graphql
mutation{
  blockUser(input: {
    reqUser: ""
    blockUser: ""
  }, password: "")
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
  view(
    input: {
      postID: "5f79cc2f9ec125d75f2e36e4"
      userId: -1
      time: 1.0
      deviceDetails: 1
    }
    password: "ddd"
  )
}


```


---
New report
```graphql
mutation review {
  report(input: { userID: 1, postID: "5f79cc2f9ec125d75f2e36e4" })
}

```


---
New share
```graphql
mutation share {
  share(
    input: { userID: 1, postID: "5f79cc2f9ec125d75f2e36e4" }
    password: "kub"
  )
}

```


---
Update nets
```graphql
mutation learning {
  learning(input: { recommender: 1, detector: 1 }, password: "dd")
}

```



---
Get stats
```graphql
query stats {
  getStats(id: "5fa53fa53c01bd8b20cd13f9"){
    text
    userid
    views{
      localization
      date
    }
  }
}

```



---
Check google storage images integrity
```graphql
query {
  storageIntegrity(password: "")
}
```