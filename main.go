package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

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
	Name       string `json:"name"` //To add column name `gorm:"column:beast_name"`
	ID         int    `gorm:"unique" json:"id"`
	Department string `json:"department"`
	Location   Adress `json:"location" gorm:"embedded"`
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

	// server.DB.AutoMigrate(&Employee{})
	server.Router = mux.NewRouter()
	server.InitalizeRoutes()
	Load(server.DB)

}

func (s *Server) InitalizeRoutes() {
	s.Router.HandleFunc("/Get", GetAllEmployee).Methods("GET")
	s.Router.HandleFunc("/Create", CreateEmpployee).Methods("POST")
	s.Router.HandleFunc("/Delete/{id}", DeleteById).Methods("DELETE")
	// s.Router.HandleFunc("/GetByValue/{city}/{name}/{department}/{street}", GetEmployeeByInfo).Methods("GET")
	s.Router.HandleFunc("/GetByValue/{value}", GetEmployeeByInfo).Methods("GET")
	s.Router.HandleFunc("/GetByID/{id}", GetEmployeeByID).Methods("GET")
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
	// if !DB.HasTable
	err = DB.Debug().AutoMigrate(&Employee{}).Error
	// err = DB.Debug().AutoMigrate(&Adress{}).Error
	if err != nil {
		log.Fatal("Cannot automigrate the Table")
	}
	// DB.Model(&Employee{}).Related(&Adress{})
	for i, _ := range Employ {
		err = DB.Debug().Model("Employee").Create(&Employ[i]).Error
		// err = DB.Create(&Employ[len(Employ)-1]).Error
		//err := DB.Model(&Employee{}).Related(&Adress{}).Save(&Employ[i]).Error
		if err != nil {
			log.Fatal("Cannot push employee detail in DB %v", err)
		}
		// This function will re-create dummy data
		// err = DB.Debug().Model(&Employee{}).Create(&users[i]).Error()

	}
}

func GetAllEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome to Get employee details endpoint")

	// emp := Employee{}
	var AllEmp []Employee
	err := server.DB.Debug().Model(&Employee{}).Limit(100).Find(&AllEmp).Error
	if err != nil {
		panic(err)
		// return
	}
	// fmt.Println(w, )
	response, _ := json.Marshal(&AllEmp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

	// json.NewEncoder(w).Encode(AllEmp)

	// Get Employees is working totally fine

}

func CreateEmpployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome to Create Employee End Point")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var emp Employee
	err = json.Unmarshal(data, &emp) // Unmarshalling post content in struct

	defer r.Body.Close()

	if emp.Name == "" {
		// fmt.Println("employee name is compulsory")
		// os.Exit(1)
		w.WriteHeader(http.StatusBadRequest)
		return
		// break
	} else if emp.Department == "" {
		w.WriteHeader(http.StatusBadRequest)
		// w.Write([]byte(Content not available)
		// fmt.Println("Department is a compulsory")
		return

	}
	//Have to close program execution
	// break
	if emp.ID == 0 {
		// emp.ID = rand.Intn(1000)
		// n1 := rand.NewSource(time.Now().UnixNano())
		// random := rand.New(n1)
		// emp.ID = random.Int()
		// emp.ID = rand.Intn(10000)
		// random := rand.Int63n(time.Now().Unix()-94608000) + 94608000
		// emp.ID = int(random)
		// emp.ID = rand.Intn(100)
		// if
		// Logic to get random number every time, if id is given 0

		var AllEmp []Employee
		err := server.DB.Debug().Model(&Employee{}).Limit(100).Find(&AllEmp).Error
		if err != nil {
			panic(err)
			// return

		}
		for _, employee := range AllEmp {
			var number int
			number = CreatingRandomNumber()
			if employee.ID == number {
				w.Write([]byte("ID is already taken"))
				return
			} else {
				emp.ID = number
			}
		}

		if err != nil {
			panic(err)
		}

		// server.DB.Save()

		// err = server.DB.Debug().Create(&emp).Error
		// if err != nil {
		// 	panic(err)
		// }

		err = server.DB.Save(&emp).Error
		if err != nil {
			panic(err)
		}
		// fmt.Print;ln(w, "Showing latest updated employee ", (json.NewEncoder(w).Encode(&emp)))
		response, _ := json.Marshal(&emp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}

func CreatingRandomNumber() int {
	random := rand.Intn(100)
	return random

}

func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["id"]
	pid, _ := strconv.ParseInt(value, 10, 64)

	var emp Employee
	err := server.DB.First(&emp, Employee{ID: int(pid)}).Error
	if err != nil {
		panic(err)
	}

	response, _ := json.Marshal(&emp)

	// json.NewEncoder(w).Encode(&emp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

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

	// err := server.DB.Where("City = ?", value).Or("City = ?", value).Or("name = ?", value).Or("Department = ?", value).Scan(&emp)
	// if err != nil {
	// 	panic(err)
	// }

	err := server.DB.First(&emp, Employee{Name: value}).Error
	if err != nil {
		panic(err)
	}
	response, _ := json.Marshal(&emp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	// json.NewEncoder(w).Encode(&emp)

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

	server.DB.Where("id = ?", pid).Find(&emp)
	server.DB.Delete(&emp)
	response, err := json.Marshal(&emp)
	if err != nil {
		panic(err)
	}

	// RespondJSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func Run(addr string) {

	fmt.Println("Listening on Port %d", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))

}
