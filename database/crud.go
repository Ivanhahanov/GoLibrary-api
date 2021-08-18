package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ivanhahanov/GoLibrary/config"
	"github.com/Ivanhahanov/GoLibrary/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Timeout operations after N seconds
	connectTimeout           = 5
	connectionStringTemplate = "mongodb://%s:%s@%s"
)

type MongoCredentials struct {
	Username string
	Password string
	Enrypoint string
}

var mongoCfg MongoCredentials

func InitConnection(cfg *config.Config){
	mongoCfg.Enrypoint = cfg.Database.Address
	mongoCfg.Username = cfg.Database.Username
	mongoCfg.Password = cfg.Database.Password
}

// GetConnection Retrieves a client to the MongoDB
func GetConnection() (*mongo.Client, context.Context, context.CancelFunc) {
	username := mongoCfg.Username
	password := mongoCfg.Password
	clusterEndpoint := mongoCfg.Enrypoint

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, clusterEndpoint)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client, ctx, cancel
}

// GetAllBooks Retrives all books from the db
func GetAllBooks() ([]*models.Book, error) {
	var books []*models.Book

	client, ctx, cancel := GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("books")
	collection := db.Collection("books")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &books)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return books, nil
}

// GetBookByID Retrives a task by its id from the db
func GetBookByID(id primitive.ObjectID) (*models.Book, error) {
	var task *models.Book

	client, ctx, cancel := GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("books")
	collection := db.Collection("books")
	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result == nil {
		return nil, errors.New("Could not find a Book")
	}
	err := result.Decode(&task)

	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	log.Printf("Books: %v", task)
	return task, nil
}

//Create creating a task in a mongo
func Create(task *models.Book) (primitive.ObjectID, error) {
	client, ctx, cancel := GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	task.ID = primitive.NewObjectID()

	result, err := client.Database("books").Collection("books").InsertOne(ctx, task)
	if err != nil {
		log.Printf("Could not create Book: %v", err)
		return primitive.NilObjectID, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

//Update updating an existing task in a mongo
func Update(task *models.Book) (*models.Book, error) {
	var updatedBook *models.Book
	client, ctx, cancel := GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": task,
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database("books").Collection("books").FindOneAndUpdate(ctx, bson.M{"_id": task.ID}, update, &opt).Decode(&updatedBook)
	if err != nil {
		log.Printf("Could not save Book: %v", err)
		return nil, err
	}
	return updatedBook, nil
}

func OneDelete(bookId string)  (string, error){
	client, ctx, cancel := GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	objectId, _ := primitive.ObjectIDFromHex(bookId)
	fmt.Println(objectId)
	result, err := client.Database("books").Collection("books").DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil{
		log.Printf("Could not delete Book: %v", err)
	}
	if result.DeletedCount == 0 {
		log.Println("DeleteOne() document not found:", result)
	} else {
		// Print the results of the DeleteOne() method
		log.Println("DeleteOne Result:", result)
	}
	return bookId, nil
}

