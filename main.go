package main

import (
	"fmt"

	"github.com/crud_api/api"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type Employee struct {
	Name string
	ID   int
	City string
}

type Employ Employee{
	Employee{
		Name: "Chirag Gupta",
		ID: 764,
		City: "Delhi",
	},
	Employee{
		Name: "Ankit Jain",
		ID: 987,
		City:"Bangalore",
	},
}

func main() {

	Server.Initialize()
	// Seed()
	Server.Load()
	Run(:8099)
}

func (server *Server) Initialize() {
	server.DB, err = gorm.Open("postgres", "host=localhost port=5432 user=aicumendeveloper dbname=postgres password=dev sslmode=disable")
	if err != nil {
		fmt.Printf("Cannot connect to databaase")
	}
	err := server.DB.DB().Ping()
	if err != nil {
		panic(err)
		fmt.Println("Not connected with Database")

	}

	server.DB.AutoMigrate(&Employee{})
	server.Router = mux.NewRouter()
	Server.InitalizeRoutes()
	Server.Load()

}

func (s *Server) InitalizeRoutes() {
	s.Router.HandleFunc("/Get", GetEmployee).Methods("GET")
	s.Router.HandleFunc("/Create", CreateEmpployee).Methods("POST")
}

// func (s *Server)Load() {
// 	err := s.DB.De
// }

func Run(addr string) {

	fmt.Printn("Listening on Port %d", add)
	log.Fatal(http.ListenAndServe(addr, Server.Router))

}