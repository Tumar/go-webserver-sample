package models

import "database/sql"

type Book struct {
	Isbn   string  `json:"isbn" binding:"required,isbn"`
	Title  string  `json:"title" binding:"required"`
	Author string  `json:"author" binding:"required"`
	Price  float32 `json:"price" binding:"required,gt=0"`
}

func AllBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book

		err := rows.Scan(&book.Isbn, &book.Title, &book.Author, &book.Price)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func CreateBook(db *sql.DB, b *Book) (*Book, error) {
	_, err := db.Exec(`
		INSERT INTO books (isbn, title, author, price)
		VALUES ($1, $2, $3, $4)`,
		b.Isbn, b.Title, b.Author, b.Price,
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func FindBookByIsbn(db *sql.DB, isbn string) (*Book, error) {
	var book Book
	row := db.QueryRow("SELECT * FROM books WHERE isbn = $1", isbn)
	err := row.Scan(&book.Isbn, &book.Title, &book.Author, &book.Price)
	if err != nil {
		return nil, err
	}

	return &book, nil
}
