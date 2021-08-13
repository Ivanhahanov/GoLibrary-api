package database

import (
	"github.com/Ivanhahanov/GoLibrary/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type MongoSearch struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Publisher   string   `json:"publisher"`
}

func Search(request *MongoSearch) (books []*models.Book) {
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("books")
	collection := db.Collection("books")
	result, err := collection.Find(ctx, bson.M{
		"$and": []interface{}{
			bson.M{"title": bson.M{"$regex": request.Title, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": request.Description, "$options": "i"}},
			bson.M{"author": bson.M{"$regex": request.Author, "$options": "i"}},
			bson.M{"publisher": bson.M{"$regex": request.Publisher, "$options": "i"}},
		},
	})
	defer result.Close(ctx)
	if err != nil {
		log.Println("here", err)
	}
	err = result.All(ctx, &books)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil
	}
	return
	// {"title": {"$regex": title, "$options": "i"}}
}
