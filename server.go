package main

import (
	"fmt"
	"log"
)

func main() {
	lib, err := OpenDB("handy_book_library.db")
	if err != nil {
		log.Fatal(err)
	}
	defer lib.db.Close()

	//// Example of querying all books
	//books, err := lib.queryAllBooks()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, book := range books {
	//	fmt.Printf("Title: %s, Author: %s\n", book.Title.String, book.Author.String)
	//}

	//// Example of full-text search in books
	searchResults, err := lib.fullTextSearchBooks("MVC")
	if err != nil {
		log.Fatal(err)
	}
	for _, book := range searchResults {
		fmt.Printf("Search Result - Title: %s, Author: %s\n", book.Title, book.Author)
	}
}
