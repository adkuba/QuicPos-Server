# Table of Contents
- [What I've learned]()
- [Important notes](#important-notes)
- [Tensorflow](#tensorflow)
  - [Installation](#installation)
  - [Additional steps](#additional-steps)
  - [Net structure](#net-structure)
- [Golang notes](#golang-notes)
  - [Recreate project](#recreate-project)
  - [Working with go](#working-with-go)
  - [GraphQL](#GraphQL)
- [Mongodb](#mongodb)
  - [Localization stats](#localization-stats)
  - [Daily users](#daily-users)
  - [Links del](#links-del)
- [Deploy](#deploy)
  - [Build server](#build-server)
  - [On server](#on-server)
- [API](#api)


# What I've learned
* How to create backend application in **Go** with **GraphQL**
* Working with **MongoDB** and **Google cloud**
* Implementing **Stripe** for payments and **Tensorflow** for machine learning
* Developing **recommender system**
* Deploying server with **SSL**


# Important notes
* Server only consumes jpeg images
* Go to [Playground](https://www.api.quicpos.com)
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

### Additional steps
* Installing tensorflow C [tutorial](https://www.tensorflow.org/install/lang_c)
* Interesting [tutorial](https://tonytruong.net/running-a-keras-tensorflow-model-in-golang/)
* All values needs to be float32

### Net structure
To see detailed net structure go to [QuicPos-Microservice]() repository.



# Golang notes

### Recreate project
* In main directory run <code>go mod init QuicPos</code>
* Im using graphQL package
* Create project structure with: <code>go run github.com/99designs/gqlgen init</code>
* Delete CreateTodos and Todos from <code>schema.resolvers.go</code>

### Working with go
Important notes:
- Simple installation - only remember about the PATH variable. To change it I needed to edit <code>/etc/environment</code>

- Example go package download <code>go get go.mongodb.org/mongo-driver/mongo</code>

- Important! Go to file -> preferences -> settings -> go (extension) -> change format tool to "gofmt"

### GraphQL
GraphQL schema is important, it defines how results will be sent, data structure.
- In <code>schema.graphqls</code> define your models and operations - mutations and queries.
- Generate functions with <code>go run github.com/99designs/gqlgen generate</code>
- Go to <code>schema.resolvers.go</code> and implement functions.



# Mongodb
- <code>internal/mongodb</code> - file with all mongo database scripts.
- If mongo can't find DNS edit <code>/etc/resolv.conf</code> and change nameserver to 8.8.8.8 <br>

### Localization stats
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


### Daily users
```
[{
    $match: {
      views: {
        $ne: null
      }
    }
}, {
    $project: {
      'views.user': 1,
      'views.date': 1,
      _id: 0
    }
}, {
    $unwind: {
      path: '$views'
    }
}, {
    $project: {
      "views.user": "$views.user",
      "views.date": {
        $dateToString: {
          date: "$views.date"
        }
      }
    }
}, {
    $project: {
      "views.user": "$views.user",
      "views.date": {
        $substr: ["$views.date", 0, 10]
      }
    }
}, {
    $group: {
      _id: {
        date: '$views.date',
        user: '$views.user'
      },
      count: {
        $sum: 1
      }
    }
}, {
    $group: {
      _id: "$_id.date",
      count: {
        $sum: 1
      }
    }
}, {
    $sort: {
      _id: 1
    }
}]
```

### Links del
Query deleting all links in post's text:
```sh
db.posts.find({text: {$regex: "https:[^ ]+"}}).forEach(function(e, i){ 
    const regex = /https:[^ ]+/gi; 
    e.text = e.text.replace(regex, ''); 
    db.posts.save(e); 
})
```



# Deploy

### Build server
Build server executable:
* Build with <code>go build -o bin/</code> and check with <code>bin/QuicPos</code>
* Send to server with scp <code>scp bin/QuicPos root@142.93.232.180:~/QuicPos</code>
* Remember to send needed files! May need scp with -r flag

### On server
You don't have to create docker image. Just build golang project, send to virtual machine and run. Then execute:
* Better option, Nohup is better than disown for startig server in Linux [More detailed answer](https://unix.stackexchange.com/questions/3886/difference-between-nohup-disown-and)

  ```sh
  sudo nohup ./QuicPos &
  ```
* Second option

  ```sh
  ./Quicpos &>> server.txt
  [1]+  Stopped                 myprogram
  $ disown -h %1
  $ bg 1
  [1]+ myprogram &
  $ logout
  ```




# API
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
