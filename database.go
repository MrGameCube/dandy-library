package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type BookLibrary struct {
	ID                 sql.NullInt64  `db:"_id"`
	ISBN               sql.NullString `db:"ISBN"`
	ISBN10             sql.NullString `db:"ISBN10"`
	Title              sql.NullString `db:"Title"`
	Author             sql.NullString `db:"Author"`
	Author2            sql.NullString `db:"Author2"`
	Author3            sql.NullString `db:"Author3"`
	Translator         sql.NullString `db:"Translator"`
	Publisher          sql.NullString `db:"Publisher"`
	PublishedDate      sql.NullString `db:"Published_Date"`
	PagesNumber        sql.NullInt32  `db:"Pages_Number"`
	Series             sql.NullString `db:"Series"`
	Volume             sql.NullInt32  `db:"Volume"`
	Language           sql.NullString `db:"Language"`
	ImageUrl           sql.NullString `db:"Image_Url"`
	IconPath           sql.NullString `db:"Icon_Path"`
	PhotoPath          sql.NullString `db:"Photo_Path"`
	Summary            sql.NullString `db:"Summary"`
	HardCover          sql.NullInt32  `db:"Hard_Cover"`
	Ebook              sql.NullString `db:"Ebook"`
	Width              sql.NullInt32  `db:"Width"`
	Height             sql.NullInt32  `db:"Height"`
	Price              sql.NullInt32  `db:"Price"`
	Category           sql.NullString `db:"Category"`
	Rating             sql.NullString `db:"Rating"`
	RatingCount        sql.NullInt32  `db:"Rating_Count"`
	Reviews            sql.NullString `db:"Reviews"`
	ReviewsCount       sql.NullInt32  `db:"Reviews_Count"`
	ItemUrl            sql.NullString `db:"Item_Url"`
	ItemAffiliateIdUrl sql.NullString `db:"Item_AffiliateId_Url"`
	Copy               sql.NullInt32  `db:"Copy"`
	Read               sql.NullInt32  `db:"Read"`
	Favorite           sql.NullInt32  `db:"Favorite"`
	Wish               sql.NullInt32  `db:"Wish"`
	LendOrBorrow       sql.NullInt32  `db:"Lend_or_Borrow"`
	Person             sql.NullString `db:"Person"`
	StartDate          sql.NullString `db:"Start_Date"`
	DueDate            sql.NullString `db:"Due_Date"`
	ReturnDate         sql.NullString `db:"Return_Date"`
	Comment            sql.NullString `db:"Comment"`
	Location           sql.NullString `db:"Location"`
}

type HandyLibraryDB struct {
	db *sqlx.DB
}

func OpenDB(dataSourceName string) (*HandyLibraryDB, error) {
	db, err := sqlx.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	err = initFTS(db)
	return &HandyLibraryDB{db}, err
}

func initFTS(db *sqlx.DB) error {
	// Define the FTS virtual table creation query
	createFTSTableSQL := `
	CREATE VIRTUAL TABLE IF NOT EXISTS book_library_fts USING fts4(
		Title,
		Author,
		Author2,
		Author3,
		Translator,
		Publisher,
		Summary,
		Category,
	);`

	// Execute the FTS table creation
	_, err := db.Exec(createFTSTableSQL)
	if err != nil {
		return err
	}
	// Populate the FTS table with data from the main table
	populateFTSTableSQL := `
	INSERT INTO book_library_fts(rowid, Title, Author, Author2, Author3, Translator, Publisher, Summary, Category)
	SELECT _id, Title, Author, Author2, Author3, Translator, Publisher, Summary, Category FROM book_library
	WHERE _id NOT IN (SELECT rowid FROM book_library_fts);
	`
	_, err = db.Exec(populateFTSTableSQL)
	if err != nil {
		return err
	}
	return nil
}

// queryAllBooks retrieves all books from the book_library table using sqlx
func (lib *HandyLibraryDB) queryAllBooks() ([]BookLibrary, error) {
	var books []BookLibrary
	err := lib.db.Select(&books, "SELECT * FROM book_library")
	if err != nil {
		return nil, err
	}
	return books, nil
}

// fullTextSearchBooks performs a full-text search on the book_library table using sqlx.
func (lib *HandyLibraryDB) fullTextSearchBooks(searchTerm string) ([]BookLibrary, error) {
	var books []BookLibrary
	query := "SELECT * FROM book_library_fts WHERE book_library_fts MATCH ?"
	err := lib.db.Select(&books, query, searchTerm)
	if err != nil {
		return nil, err
	}
	return books, nil
}
