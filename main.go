package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)
func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}
var router *mux.Router

func CreateRouter() {
	router = mux.NewRouter()
}
func InitializeRoute() {
	router.HandleFunc("/signup", SignUp).Methods("POST")
	router.HandleFunc("/signin", SignIn).Methods("POST")
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}
type Error struct {
	IsError bool `json:"isError"`
	Message string `json:"message"`
}


func GetDatabase() *gorm.DB {
	databasename := "userdb"
	database := "postgres"
	databasepassword := "babygetmygun98"
	databaseurl := "postgres://postgres:" + databasepassword + "@localhost/" + databasename + "?sslmode=disable"
	connection, err := gorm.Open(database, databaseurl)
	if err != nil {
		log.Fatal("wrong database url ")
	}
	sqldb := connection.DB

	err = sqldb.Ping()
	if err != nil {
		log.Fatal("database connected")
	}
	fmt.Println("connected to database")
	return connection
}
func InitialMigration() {
	connection := GetDatabase()
	defer Closedatabase(connection)
	connection.AutoMigrate(User{})
}
func Closedatabase(connection *gorm.DB) {
	sqldb := connection.DB()
	sqldb.Close()
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	connection := GetDatabase()
	defer Closedatabase(connection)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err error
		err = SetError (err, "Error reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewDecoder(w).Encode(err)
		return
	}

}
var dbuser User
connection.Where("email=?", user.Email).First((&dbuser))

//cehck if email is already register or not
if dbuser.Email != ""{
	var err error
	err = SetError(err, "Email already registered")
	w.Header().Set("Content-Type", "application/")
	json.NewDecoder(w).Encode(err)
	return
}

user.Password,err = GeneratehashPassword(user.Password)
if err != nil {
	log.Fatal("Error in password hash")
}

//insert user details in database
connection.Create(&user)
w.Header().Set("Content-Type", "application/json")
json.NewDecoder(w).Encode(user)

func main() {
	CreateRouter()
	InitializeRoute()

}
