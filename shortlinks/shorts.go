package shortlinks

import (
	"errors"
	"fmt"
	"github.com/Ivanhahanov/GoLibrary/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Shorter struct {
	ShortName    string `json:"short_name"`
	OriginalName string `json:"original_name"`
	VisitCount   int    `json:"visit_count"`
}

func CreateShortLink(originalLink string) (string, error) {
	shortLink := GenerateShorLink(4)
	for {

		if checkShortLinkExists(shortLink) {
			break
		}
		shortLink = GenerateShorLink(4)
	}
	shorter := Shorter{
		OriginalName: originalLink,
		ShortName:  shortLink,
		VisitCount: 0,
	}
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("shorter")
	collection := db.Collection("shorter")
	_, err := collection.InsertOne(ctx, shorter)
	if err != nil {
		log.Printf("Could not create Book: %v", err)
		return "", err
	}
	return shorter.ShortName, nil
}

func checkShortLinkExists(shortLink string) bool {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("shorter")
	collection := db.Collection("shorter")
	result := collection.FindOne(ctx, bson.M{"shortname": shortLink})
	if result == nil {
		return false
	}
	return true
}

func getDocumentByShortLink(shortLink string) (*Shorter, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("shorter")
	collection := db.Collection("shorter")
	result := collection.FindOne(ctx, bson.M{"shortname": shortLink})
	if result == nil {
		return nil, errors.New(fmt.Sprintf("could not find a Record for %s", shortLink))
	}

	var shorter *Shorter
	err := result.Decode(&shorter)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return shorter, nil
}

func GetOriginalLink(shortLink string) (originalLink string, err error) {
	doc, err := getDocumentByShortLink(shortLink)
	if err != nil {
		return "", err
	}
	return doc.OriginalName, nil
}

func GetAllDocuments() ([]*Shorter, error) {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("shorter")
	collection := db.Collection("shorter")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var shorters []*Shorter
	err = cursor.All(ctx, &shorters)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return shorters, nil
}

func WriteVisit(shortLink string) error {
	client, ctx, cancel := database.GetConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("shorter")
	collection := db.Collection("shorter")
	_, err := collection.UpdateOne(ctx,
		bson.M{"shortname": shortLink}, bson.D{
			{"$inc", bson.M{"visitcount": 1}},
		}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return err
	}
	return nil
}
