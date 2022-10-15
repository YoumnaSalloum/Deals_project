package main

// ? Require the packages
import (
	"context"
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/YoumnaSalloum/golang-test/config"
	models "github.com/YoumnaSalloum/golang-test/models"
	routes "github.com/YoumnaSalloum/golang-test/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Create required variables that we'll re-assign later

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client
)

// ? Init function that will run before the "main" function
func init() {

	// ? Load the .env variables
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// ? Create a context
	ctx = context.TODO()

	// ? Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// ? Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// ? Create the Gin Engine instance
	server = gin.Default()
}

func main() {
	//all availbe rates for testing code
	var rates = []models.Rate{
		{Cfrom: "USD", Cto: "GBP", Conv: 1.5},
		{Cfrom: "GBP", Cto: "USD", Conv: 0.5},
	}
	// type request struct {
	//     t time.Time  `json:"t"`
	//     tcur  string `json:"tcur"`
	//     fcur string  `json:"fcur"`
	//     amount  float64 `form:"amount"`
	// }
	// var txs = []request{
	//     {t: time.Now() , fcur: "USD", tcur: "GBP", amount: 1000},
	// 	// {t: time.Now() , fcur: "GBP", tcur: "USD", amount: 1000},
	// }
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})
	///////////API
	router.POST("/exchange", func(c *gin.Context) {
		var result float32
		curr := new(models.BasicExchange)
		c.ShouldBindJSON(curr)
		var flag bool = true

		//check created_at is before our current time not in future
		if curr.CreatedAt.After(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Sorry but you cannot do the exchange for the currencies the time of createdAt is in future"})
			return
		}
		// check if amount is positive number
		if curr.Amount > 0 {

			//loop over availble rates
			for _, v := range rates {
				//check all elements in array of rates in request.
				for i, _ := range curr.Rates {
					if v.Cfrom == curr.Rates[i].Cfrom && v.Cto == curr.Rates[i].Cto {
						result = float32(curr.Amount) * v.Conv
						flag = false
						break
					}
				}
			}

		}
		//check if id is empty then create new unique id else create id from request id
		var er error
		var id primitive.ObjectID
		if curr.Id.Hex() == "" {
			id = primitive.NewObjectID()
		} else {
			id, er = primitive.ObjectIDFromHex(curr.Id.Hex())
			if er != nil {
				return
			}
		}
		//||curr.CreatedAt==""||curr.Rates==""||curr.CodeId==""||curr.Id==""
		// if(curr.Amount==0){
		// panic("feild is missing or null!")
		// }
		curObj := &models.BasicExchange{
			CreatedAt: curr.CreatedAt,
			Rates:     curr.Rates,
			Id:        id,
			Amount:    curr.Amount,
			NewAmount: result,
			CodeId:    curr.CodeId,
		}
		//save in db
		_, err := routes.CreatePost(c, curObj)
		if err != nil {
			log.Print(err)
			return
		}
		if flag {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No Rate Available"})
			return
		} else {
			c.JSON(http.StatusCreated, gin.H{"success": true, "message": "data created successfully", "data": curObj})
		}

	})

	log.Fatal(server.Run(":" + config.Port))
}

//docker-compose up -d && go run main.go
