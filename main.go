package main

import (
	"fmt"
	"log"
	"net/http"

	// "github.com/crud_api/api"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

var Employ = []Employee{
	Employee{
		Name: "Chirag Gupta",
		ID:   764,
		City: "Delhi",
	},
	Employee{
		Name: "Ankit Jain",
		ID:   987,
		City: "Bangalore",
	},
}

var server = Server{}

func main() {

	server.Initialize()
	// Seed()
	// Server.Load()
	Load(server.DB)
	Run(":8099")
}

func (server *Server) Initialize() {
	var err error
	server.DB, err = gorm.Open("postgres", "host=localhost port=5432 user=aicumendeveloper dbname=postgres password=dev sslmode=disable")
	if err != nil {
		fmt.Printf("Cannot connect to databaase")
	}
	err = server.DB.DB().Ping()
	if err != nil {
		panic(err)
		fmt.Println("Not connected with Database")

	} else {
		fmt.Println("Server connection is success")
	}

	server.DB.AutoMigrate(&Employee{})
	server.Router = mux.NewRouter()
	server.InitalizeRoutes()
	Load(server.DB)

}

func (s *Server) InitalizeRoutes() {
	s.Router.HandleFunc("/Get", GetEmployee).Methods("GET")
	s.Router.HandleFunc("/Create", CreateEmpployee).Methods("POST")
}

// func (s *Server)Load() {
// 	err := s.DB.De
// }

func Load(DB *gorm.DB) {
	err := DB.Debug().DropTableIfExists(&Employee{}).Error
	if err != nil {
		panic(err)
	}
	err = DB.Debug().AutoMigrate(&Employee{}).Error
	if err != nil {
		log.Fatal("Cannot automigrate the Table")
	}

	for i, _ := range Employ {
		err = DB.Debug().Model(&Employee{}).Create(&Employ[i]).Error
		if err != nil {
			log.Fatal("Cannot push employee detail in DB %v", err)
		}
		// err = DB.Debug().Model(&Employee{}).Create(&users[i])

	}

}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome to Get employee details endpoint")
}

func CreateEmpployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome to Create Employee End Point")
}

func Run(addr string) {

	fmt.Println("Listening on Port %d", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))

}
