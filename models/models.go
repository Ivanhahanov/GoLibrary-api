package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID     primitive.ObjectID
	Title  string
	Author string
	Tags   []string
	Path   string
}
