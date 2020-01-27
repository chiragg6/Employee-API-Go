package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
	// gorm.Model
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
	// Server.Initialize()
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
	s.Router.HandleFunc("/Delete/{id}", DeleteById).Methods("DELETE")
}

// func (s *Server)Load() {
// 	err := s.DB.De
// }

func Load(DB *gorm.DB) {
	// err := DB.Debug().DropTableIfExists(&Employee{}).Error
	// if err != nil {
	// 	panic(err)
	// }
	// This function will drop Employee Table if exits
	err := DB.Debug().AutoMigrate(&Employee{}).Error
	if err != nil {
		log.Fatal("Cannot automigrate the Table")
	}

	// for i, _ := range Employ {
	// 	err = DB.Debug().Model(&Employee{}).Create(&Employ[i]).Error
	// 	if err != nil {
	// 		log.Fatal("Cannot push employee detail in DB %v", err)
	// 	}
	// This function will re-create dummy data
	// err = DB.Debug().Model(&Employee{}).Create(&users[i])

}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome to Get employee details endpoint")

	// emp := Employee{}
	var AllEmp []Employee
	err := server.DB.Debug().Model(&Employee{}).Limit(100).Find(&AllEmp).Error
	if err != nil {
		panic(err)
		// return
	}
	// fmt.Println(w, )
	json.NewEncoder(w).Encode(AllEmp)

}

func CreateEmpployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome to Create Employee End Point")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var emp Employee
	err = json.Unmarshal(data, &emp)

	if emp.ID == 0 {
		// emp.ID = rand.Intn(1000)
		n1 := rand.NewSource(time.Now().UnixNano())
		random := rand.New(n1)
		emp.ID = random.Int()
		// Logic to get random number every time, if id is given 0

	}
	if err != nil {
		panic(err)
	}

	err = server.DB.Debug().Create(&emp).Error
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(&emp)

}

// func FindByCondition(w http.ResponseWriter, r *http.Request) {
// 	var err error
// 	vars := mux.Vars(r)
// 	// condition := vars["id"] || vars["city"]

// }

func DeleteById(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	id := vars["id"]
	pid, _ := strconv.ParseInt(id, 10, 64)
	// fmt.Println(pid)

	err = server.DB.Debug().Model(&Employee{}).Where("ID = ?", pid).Take(&Employee{}).Delete(&Employee{}).Error
	if err != nil {
		panic(err)
	}

	server.DB.Where("ID = ?", pid).Find(&Employee)
	server.DB.Delete(&Employee)
}

func Run(addr string) {

	fmt.Println("Listening on Port %d", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))

}
