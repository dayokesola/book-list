package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load()
}
func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	psqlInfo := os.Getenv("DB_CONN")

	log.Println(psqlInfo)
	db, err = sql.Open("postgres", psqlInfo)
	logFatal(err)
	//defer db.Close()

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("get all books")

	var book Book
	books = []Book{}
	rows, err := db.Query("select * from books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	log.Println("get a book")
	params := mux.Vars(r)
	var book Book
	rows := db.QueryRow("select * from books where id=$1", params["id"])
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)
	json.NewEncoder(w).Encode(&book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("add a book")
	var book Book
	var id int
	_ = json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("insert into books (title, author, year) values($1,$2,$3) returning id;",
		book.Title, book.Author, book.Year).Scan(&id)
	logFatal(err)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("update a book")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	res, err := db.Exec("update books set title = $1, author = $2, year = $3 where id = $4;",
		book.Title, book.Author, book.Year, book.ID)
	logFatal(err)

	rowsUpdated, err := res.RowsAffected()
	logFatal(err)
	json.NewEncoder(w).Encode(rowsUpdated)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	log.Println("delete a book")
	var params = mux.Vars(r)

	res, err := db.Exec("delete from books where id = $1;", params["id"])
	logFatal(err)

	rowsUpdated, err := res.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}
