package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID     primitive.ObjectID `json:"id"`
	Title  string             `json:"title"`
	Author string             `json:"author"`
	Tags   []string           `json:"tags"`
	Path   string             `json:"path"`
}
