package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `json:"title"`
	Publisher   string             `json:"publisher"`
	Author      string             `json:"author"`
	Tags        []string           `json:"tags"`
	Description string             `json:"description"`
	Slug        string             `json:"slug"`
	Path        string             `json:"path"`
}
