package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	getcollection "github.com/YoumnaSalloum/golang-test/Collection"
	database "github.com/YoumnaSalloum/golang-test/databases"
	models "github.com/YoumnaSalloum/golang-test/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Database("myTaskDB").Collection("Currs")
func CreatePost(c *gin.Context, cur *models.BasicExchange) (primitive.ObjectID, error) {
	var DB = database.ConnectDB()
	var curCollection = getcollection.GetCollection(DB, "Currs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := curCollection.InsertOne(ctx, cur)
	if err != nil {
		log.Fatal(err)
		return primitive.NilObjectID, err
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": map[string]interface{}{"data": res}})

	return res.InsertedID.(primitive.ObjectID), nil
}