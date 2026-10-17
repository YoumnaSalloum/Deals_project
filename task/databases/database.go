package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/YoumnaSalloum/golang-test/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
)

// ConnectDB returns one shared MongoDB client (created on first use).
// The URI comes from MONGODB_LOCAL_URI instead of a hardcoded 127.0.0.1,
// so it works inside Kubernetes where Mongo is another pod.
func ConnectDB() *mongo.Client {
	once.Do(func() {
		cfg, err := config.LoadConfig(".")
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBUri))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to mongoDB")
	})
	return client
}
