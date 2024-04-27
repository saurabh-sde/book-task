package model

type Book struct {
	BookName        string `json:"bookName"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publicationYear"`
}
