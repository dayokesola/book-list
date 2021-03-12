package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book

func main() {
	router := mux.NewRouter()

	books = append(books, Book{ID: 1, Title: "How to fly", Author: "Charles", Year: "2005"},
		Book{ID: 2, Title: "How to Pray", Author: "Sam", Year: "2004"},
		Book{ID: 3, Title: "How to Cook", Author: "Lekan", Year: "2006"},
		Book{ID: 4, Title: "How to Dance", Author: "Idowu", Year: "2007"},
		Book{ID: 5, Title: "How to Code", Author: "Dayo", Year: "2005"})

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("get all books")
	json.NewEncoder(w).Encode(books)
}
func getBook(w http.ResponseWriter, r *http.Request) {
	log.Println("get a book")
	var params = mux.Vars(r)
	log.Println(params)
}
func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("add a book")
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("update a book")
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	log.Println("delete a book")
}
