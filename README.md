# Gin Server

Sample web server built using [Gin](https://github.com/gin-gonic/gin) (a framework for building web applications in [Go](https://golang.org))

## Getting Started

### Prerequisites
---
* Docker (v20.10.5)

### Start up
---
To build and run this project locally, run the following commands at the root directory of the project (task) directory
```bash
docker-compose up -d && go run main.go
```
This will build the docker image, start up a [MongoDB](https://www.mongodb.com) container for the project and then start the project at port `http://localhost:8000`
MongoDb on 'localhost:27017'
this is main API (POST :' http://localhost:8000/api/exchange')
request body example:
{
    "rates": [
        {
            "cfrom": "USD",
            "cto": "GBP",
            "conv": 1.5
        }
    ],
    "_id": "63499d39985b28b42b43c263",
    "amount": 5000,
    "created_at": "2022-10-14T17:31:03.735+00:00"
}

## Built With
---
* [Golang](https://golang.org) - Language used
* [Gin](https://github.com/gin-gonic/gin) - Web Framework used
* [MongoDB](https://www.mongodb.com) - Database used
  
## Documentation
---
Documentation for the Project can be found at [gin-server-doc](https://documenter.getpostman.com/view/8916756/SztG3mKJ?version=latest)

