package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

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
	_ = json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)

	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("update a book")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}

	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	log.Println("delete a book")
	var params = mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
		}
	}

	json.NewEncoder(w).Encode(books)
}
