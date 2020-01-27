package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
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
	Name       string `json:"name"`
	ID         int    `json:"id"`
	Department string `json:"department"`
	Location   Adress `json:"location"`
}

type Adress struct {
	HouseNo   int    `json:"houseno"`
	Apartment string `json:"apartment"`
	Street    string `json:"street"`

	City    string `json:"city"`
	Pincode int    `json:"pincode"`
}

var Employ = []Employee{
	Employee{
		Name:       "Chirag Gupta",
		ID:         764,
		Department: "Development",
		Location: Adress{
			HouseNo:   645,
			Apartment: "NewMarvel",
			Street:    "OldRoad",
			City:      "Delhi",
			Pincode:   600923,
		},
	},
	Employee{
		Name:       "Ankit Jain",
		ID:         987,
		Department: "Management",
		Location: Adress{
			HouseNo:   776,
			Apartment: "NewJerkey",
			Street:    "Boriwali",
			City:      "Mumbai",
			Pincode:   848742,
		},
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
	// s.Router.HandleFunc("/GetByValue/{city}/{name}/{department}/{street}", GetEmployeeByInfo).Methods("GET")
	s.Router.HandleFunc("/GetByValue/{value}", GetEmployeeByInfo).Methods("GET")
}

// func (s *Server)Load() {
// 	err := s.DB.De
// }

func Load(DB *gorm.DB) {
	err := DB.Debug().DropTableIfExists(&Employee{}).Error
	if err != nil {
		panic(err)
	}
	// This function will drop Employee Table if exits
	err = DB.Debug().AutoMigrate(&Employee{}).Error
	if err != nil {
		log.Fatal("Cannot automigrate the Table")
	}

	for i, _ := range Employ {
		err = DB.Debug().Model(&Employee{}).Create(&Employ[i]).Error
		if err != nil {
			log.Fatal("Cannot push employee detail in DB %v", err)
		}
		// This function will re-create dummy data
		// err = DB.Debug().Model(&Employee{}).Create(&users[i]).Error()

	}
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

	if emp.Name == "" {
		// fmt.Println("employee name is compulsory")
		os.Exit(1)
		// break
	} else if emp.Department == "" {

		fmt.Println("Department is a compulsory")

	}
	//Have to close program execution
	// break
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

func GetEmployeeByInfo(w http.ResponseWriter, r *http.Request) {
	// You should be able to list employees by location (either only city or both city and street), by name or by department

	vars := mux.Vars(r)
	// city := vars["city"]
	// street := vars["street"]
	// name := vars["name"]
	// department := vars["department"]

	value := vars["value"]
	var emp Employee
	// server.DB.Where("City = ? || Street = ? || Name= ? || Deparment = ?", value).Find(&emp)

	err := server.DB.Where("City = ?", value).Or("City = ?", value).Or("name = ?", value).Or("Department = ?", value).Scan(&emp)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(&emp)

}

func DeleteById(w http.ResponseWriter, r *http.Request) {
	// var err error
	vars := mux.Vars(r)
	id := vars["id"]
	pid, _ := strconv.ParseInt(id, 10, 64)
	// fmt.Println(pid)

	// err = server.DB.Debug().Model(&Employee{}).Where("ID = ?", pid).Take(&Employee{}).Delete(&Employee{}).Error
	// if err != nil {
	// 	panic(err)
	// }

	var emp Employee

	server.DB.Where("ID = ?", pid).Find(&emp)
	server.DB.Delete(&emp)
}

func Run(addr string) {

	fmt.Println("Listening on Port %d", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))

}
