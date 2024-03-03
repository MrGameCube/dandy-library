package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type BookTemplate struct {
	BookDetail BookLibrary
}

func main() {
	lib, err := OpenDB("handy_book_library.db")
	if err != nil {
		log.Fatal(err)
	}
	defer lib.db.Close()
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*.gohtml")
	router.GET("/book", func(c *gin.Context) {
		books, err := lib.queryAllBooks()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, books)
	})
	router.Static("/web/", "web/assets")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.gohtml", nil)
	})
	router.GET("/book/:id", func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}
		book, err := lib.getByID(idInt)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.HTML(200, "book.gohtml", BookTemplate{BookDetail: book})
	})

	router.GET("/templates/table", func(c *gin.Context) {
		books, err := lib.queryAllBooks()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.HTML(200, "table.gohtml", books)
	})

	router.GET("/search", func(c *gin.Context) {
		term := c.Query("term")
		books, err := lib.fullTextSearchBooks(term)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.HTML(200, "search-result.gohtml", books)
	})
	//// Example of querying all books
	//books, err := lib.queryAllBooks()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, book := range books {
	//	fmt.Printf("Title: %s, Author: %s\n", book.Title.String, book.Author.String)
	//}

	//// Example of full-text search in books
	//searchResults, err := lib.fullTextSearchBooks("MVC")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, book := range searchResults {
	//	fmt.Printf("Search Result - Title: %s, Author: %s\n", book.Title, book.Author)
	//}
	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}

}
