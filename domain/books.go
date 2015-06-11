package domain

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// BookFormat type
type BookFormat string

// constants for BookFormat type values
const (
	BookFormatEBook     BookFormat = "ebook"
	BookFormatPaperBack BookFormat = "paperback"
	BookFormatHardCover BookFormat = "hardcover"
)

// New Book struct
type NewBook struct {
	Author      string     `json:"author" bson:"author"`
	Name        string     `json:"name,omitempty" bson:"name"`
	Description string     `json:"description,omitempty" bson:"description"`
	ISBN        string     `json:"isbn,omitempty" bson:"isbn"`
	Format      BookFormat `json:"format,omitempty" bson:"format"`
}

// Update Book struct
type UpdateBook NewBook

// Book model struct
type Book struct {
	ID               bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Creator          bson.ObjectId `json:"creator" bson:"creator"`
	Author           string        `json:"author" bson:"author"`
	Name             string        `json:"name,omitempty" bson:"name"`
	Description      string        `json:"description,omitempty" bson:"description"`
	ISBN             string        `json:"isbn,omitempty" bson:"isbn"`
	Format           BookFormat    `json:"format,omitempty" bson:"format"`
	LastModifiedDate time.Time     `json:"lastModifiedDate" bson:"lastModifiedDate"`
	CreatedDate      time.Time     `json:"createdDate,omitempty" bson:"createdDate"`
}

// Books struct
type Books []Book
