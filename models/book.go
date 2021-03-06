package models

import (
	"database/sql"
)

//Book struct
type Book struct {
	ID    int    `json:"id"`
	Isbn  int    `json:"isbn"`
	Title string `json:"title"`
}

//BookResponse struct
type BookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []Book `json:"data"`
}

//ReadBooks function
func ReadBooks(db *sql.DB) BookResponse {

	data := []Book{}
	book := Book{}
	res := BookResponse{}

	results, err := db.Query("SELECT * FROM books")

	if err != nil {
		return handleBookError(res, err.Error())
	}

	for results.Next() {
		var id, isbn int
		var title string
		err = results.Scan(&id, &isbn, &title)
		if err != nil {
			return handleBookError(res, err.Error())
		}
		book.ID = id
		book.Isbn = isbn
		book.Title = title
		data = append(data, book)
	}
	res.Success = true
	res.Message = ""
	res.Data = data

	return res
}

//ReadBook function
func ReadBook(db *sql.DB, parmID int) BookResponse {

	data := []Book{}
	book := Book{}
	res := BookResponse{}

	results, err := db.Query("SELECT * FROM books WHERE id=?", parmID)

	if err != nil {
		return handleBookError(res, err.Error())
	}

	for results.Next() {
		var id, isbn int
		var title string
		err = results.Scan(&id, &isbn, &title)
		if err != nil {
			return handleBookError(res, err.Error())
		}
		book.ID = id
		book.Isbn = isbn
		book.Title = title

		data = append(data, book)
	}

	//If no data then it means the book was not found
	if len(data) == 0 {
		res.Data = data
		return handleBookError(res, "User not found")
	}

	res.Success = true
	res.Message = ""
	res.Data = data

	return res
}

//CreateBook function
func CreateBook(db *sql.DB, book Book) BookResponse {
	res := BookResponse{}

	insForm, err := db.Prepare("INSERT INTO books(isbn, title) VALUES(?,?)")
	if err != nil {
		return handleBookError(res, "Error on inserting: "+err.Error())
	}

	res2, _ := insForm.Exec(book.Isbn, book.Title)

	/*Logic needed for autogenerated IDs */
	lastInserted, _ := res2.LastInsertId()
	book.ID = int(lastInserted)

	res.Data = append(res.Data, book)
	res.Success = true
	res.Message = "Book inserted"

	return res

}

//UpdateBook func
func UpdateBook(db *sql.DB, paramID int, book Book) BookResponse {
	data := []Book{}
	res := BookResponse{}

	q := ReadBook(db, paramID)

	if len(q.Data) == 0 {
		return handleBookError(res, "Book not found")
	}

	insForm, err := db.Prepare("UPDATE books SET isbn=?, title=? WHERE id=?")
	if err != nil {
		return handleBookError(res, "Error updating: "+err.Error())
	}

	insForm.Exec(book.Isbn, book.Title, paramID)

	book.ID = paramID
	data = append(data, book)
	res.Success = true
	res.Message = "Successfull update"
	res.Data = data

	return res
}

//DeleteBook func
func DeleteBook(db *sql.DB, paramID int) BookResponse {
	data := []Book{}

	res := BookResponse{}
	res.Data = data

	q := ReadBook(db, paramID)

	if len(q.Data) == 0 {
		return handleBookError(res, "Book not found")
	}

	insForm, err := db.Prepare("DELETE FROM books WHERE id=?")

	if err != nil {
		return handleBookError(res, "Error deleting: "+err.Error())
	}

	insForm.Exec(paramID)

	res.Message = "Document deleted"
	res.Success = true

	return res
}

func handleBookError(res BookResponse, err string) BookResponse {
	res.Success = false
	res.Message = err

	return res
}
