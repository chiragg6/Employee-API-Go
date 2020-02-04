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
	s.Router.HandleFunc("/api/get", GetAllEmployee).Methods("GET")
	s.Router.HandleFunc("/api/create", CreateEmpployee).Methods("POST")
	s.Router.HandleFunc("/api/delete/{id}", DeleteById).Methods("DELETE")
	// s.Router.HandleFunc("/GetByValue/{city}/{name}/{department}/{street}", GetEmployeeByInfo).Methods("GET")
	s.Router.HandleFunc("/api/getbyvalue/", GetEmployeeByInfo).Methods("GET")
	s.Router.HandleFunc("/api/getbyid/{id}", GetEmployeeByID).Methods("GET")
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
	var AllEmp []Employee //AllEmp is an slice of type Employee struct
	err := server.DB.Debug().Model(&Employee{}).Limit(100).Find(&AllEmp).Error
	if err != nil {
		panic(err)
		w.Write([]byte(err.Error()))
		// return
	} else {

		// fmt.Println(w, )
		response, _ := json.Marshal(&AllEmp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}

	// json.NewEncoder(w).Encode(AllEmp)

	// Get Employees is working totally fine

}

func CreateEmpployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "Welcome to Create Employee End Point")
	// queryParamas := r.URL.Query()
	// id := queryParamas["id"]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var emp Employee
	// emp := new(Employee)
	err = json.Unmarshal(data, &emp) // Unmarshalling post content in struct

	defer r.Body.Close()

	if emp.Name == "" {
		// fmt.Println("employee name is compulsory")
		// os.Exit(1)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Employee name is a compulsory field"))
		return
		// break
	} else if emp.Department == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Department is a compulsory field"))
		// w.Write([]byte(Content not available)
		// fmt.Println("Department is a compulsory")
		return
	}

	//Have to close program execution
	// break
	if emp.ID == 0 {

		var employ Employee
		var number int
		number = CreatingRandomNumber()
		// emp.ID = number
		err := server.DB.Debug().Model(&Employee{}).Where("id= ?", number).Take(&employ).Error
		if err != nil {
			fmt.Println("Concerted employee id doesnt exits in DB")
			// w.Write([]byte("Employee with same employee id already exits"))

			emp.ID = number
			err = server.DB.Save(&emp).Error
			if err != nil {
				panic(err)
			}

			response, _ := json.Marshal(&emp)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		} else {
			w.Write([]byte("ID exits please retry"))
		}
	} else {

		err = server.DB.Debug().Model(&Employee{}).Where("id= ?", emp.ID).Take(&emp).Error
		if err != nil {
			err = server.DB.Save(&emp).Error
			if err != nil {
				panic(err)
			}

			response, _ := json.Marshal(&emp)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))

		} else {
			w.Write([]byte("Employee with entered ID already available"))
			return
		}

	}

	// var AllEmp []Employee
	// err := server.DB.Debug().Model(&Employee{}).Limit(100).Find(&AllEmp).Error
	// if err != nil {
	// 	panic(err)
	// 	// return

	// }
	// for i, _ := range AllEmp {
	// 	var number int
	// 	number = CreatingRandomNumber()
	// 	if number == AllEmp[i].ID {
	// 		w.Write([]byte("ID is already taken"))
	// 		return
	// 	} else {
	// 		emp.ID = number
	// 	}
	// }

	// if err != nil {
	// 	panic(err)
	// }

	// } else {
	// 	// check logic if entered employee id already present in db
	// 	// var All []Employee

	// 	// err := server.DB.Debug().Model(&Employee{}).Limit(100).Find(&All).Error
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// 	// return
	// 	// }

	// 	err = server.DB.Debug().Model(&Employee{}).Where("id= ?", emp.ID).Take(&emp).Error
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		// w.Write()
	// 		return
	// 	} else {
	// 		w.Write([]byte("ID exits"))
	// 	}

	// }
}

func CreatingRandomNumber() int {
	random := rand.Intn(100)
	return random

}

func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["id"]
	pid, _ := strconv.ParseInt(value, 10, 64)
	// w.Write([]byte(strconv.Itoa(int(pid))))

	// var AllEmp []Employee
	// err := server.DB.Debug().Model(&Employee{}).Limit(100).Find(&AllEmp).Error
	// if err != nil {
	// 	panic(err)
	// }
	var emp Employee

	err := server.DB.Debug().Model(&Employee{}).Where("id= ?", pid).Take(&emp).Error
	if err != nil {

		// panic(err)
		fmt.Println(err)
		w.Write([]byte("Content with concerned id is not available in db"))
		return
	} else {

		result, _ := json.Marshal(&emp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}

	// for i, _ := range AllEmp {
	// 	var num int
	// 	num = 88
	// 	if num != AllEmp[i].ID {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		w.Write([]byte("Data with concerned id doesnt exits"))
	// 		// os.Exit(1)
	// 		return
	// 	} else {
	// 		var emp Employee
	// 		err = server.DB.First(&emp, Employee{ID: int(pid)}).Error
	// 		if err != nil {
	// 			panic(err)
	// 			w.WriteHeader(http.StatusNotFound)
	// 			w.Write([]byte("Not Database entry with concerned employee id"))
	// 		}

	// 		response, _ := json.Marshal(&emp)

	// 		// json.NewEncoder(w).Encode(&emp)
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte(response))

	// 	}
	// }

}

func GetEmployeeByInfo(w http.ResponseWriter, r *http.Request) {
	// You should be able to list employees by location (either only city or both city and street), by name or by department

	query := r.URL.Query()

	city := query.Get("city")
	street := query.Get("street")
	name := query.Get("name")
	department := query.Get("department")

	var emp Employee
	if city != "" && street != "" {
		err := server.DB.Debug().Model(&Employee{}).Where("city = ? AND steet = ?", city, street).Find(&emp).Error
		if err != nil {
			fmt.Println(err)
		}

	} else if city != "" {
		err := server.DB.Debug().Model(&Employee{}).Where("city = ? ", city).Find(&emp).Error
		if err != nil {
			fmt.Println(err)
		}
	} else if name != "" {
		err := server.DB.Debug().Model(&Employee{}).Where("name = ?", name).Find(&emp).Error
		if err != nil {
			fmt.Println(err)
		}
	} else if department != "" {
		err := server.DB.Debug().Model(&Employee{}).Where("department = ?", department).Find(&emp).Error
		if err != nil {
			fmt.Println(err)
		}
	} else {
		w.Write([]byte("No employee info is present in db matching the query values"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, _ := json.Marshal(&emp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
	// json.NewEncoder(w).Encode(&emp)

}

// func queryParams(w http.ResponseWriter, r *http.Request) {
// 	// http://localhost:8080/products?filter=color
// 	// there is one query parameter filter

// 	// Getting query parameters

// 	// localhost:9090/products?filter=color&filter=price&filter=brand
// 	query := r.URL.Query()
// 	filters, present := query["filters"] // query will get multiple values with filters key
// 	if !present || len(filters) == 0 {
// 		fmt.Println("filters not present")
// 	}
// 	query.Get("city")
// 	w.Write([]byte(strings.Join(filters, ",")))

// }

func DeleteById(w http.ResponseWriter, r *http.Request) {
	// var err error
	vars := mux.Vars(r)
	id := vars["id"]
	pid, _ := strconv.ParseInt(id, 10, 64)
	// fmt.Println(pid)

	err := server.DB.Debug().Model(&Employee{}).Where("ID = ?", pid).Take(&Employee{}).Delete(&Employee{}).Error
	if err != nil {
		// panic(err)
		fmt.Println(err)

		w.Write([]byte("Employee with concerned ID not present in database"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Employee with concerned ID has been deleted"))

	}

}

func Run(addr string) {

	fmt.Println("Listening on Port %d", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))

}
