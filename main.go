package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "这是主页")
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "这是用户")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "这是创建用户")
}

func order(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "这是订单")
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/user", user)
	http.HandleFunc("/user/create", createUser)
	http.HandleFunc("/order", order)
	log.Fatal(http.ListenAndServe(":8080", nil))
}