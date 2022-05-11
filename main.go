package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var (
	router *mux.Router
	secretkey string ="secretkeyjwt"
)



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
	Password string `json:"-"`
	Role     string `json:"role"`
}
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"-"`
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

//set error message in Error struct
func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}
func GetDatabase() *gorm.DB {
	databasename := "userdb"
	database := "postgres"
	databasepassword := "babygetmygun98"
	databaseurl := "postgres://postgres:" + databasepassword + "@localhost/" + databasename + "?sslmode=disable"
	connection, err := gorm.Open(database, databaseurl)
	if err != nil {
		log.Fatal("Invalid database url ")
	}
	sqldb := connection.DB

	err = sqldb.Ping()
	if err != nil {
		log.Fatal("database connected")
	}
	fmt.Println("Database connection successful")
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
func GeneratehashPassword(password string) (string,error) {
	bytes,err :=bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes),err
}
func CheckPasswordHash(password ,hash string) bool {
	err :=bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err == nil
}
//GenerateJWT
func GenerateJWT(email,role string)(string,error){
	var mySigningKey =[]byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"]=true
	claims["email"]=email
	claims["role"]=role
	claims["expires"]=time.Now().Add(time.Minute *30).Unix()

	tokenString,err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something went wrong :%s",err.Error())
		return "",err
	}
	return tokenString,nil
}

//Middleware function
//check wether user is authorized or not
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		if r.Header["Token"] == nil{
			var err Error
			err =SetError(err,"No Token found")
			json.NewDecoder(w).Encode(err)
			return
		}
		var mySigningKey = []byte(secretkey)

		token,err := jwt.Parse(r.Header["Token"][0],func(token *jwt.Token)(interface{},error))
	}
}





func SignUp(w http.ResponseWriter, r *http.Request) {
	connection := GetDatabase()
	defer Closedatabase(connection)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err Error
		err = SetError (err, "Error reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewDecoder(w).Encode(err)
		return
	}

}
var dbuser User
connection.Where("email=?", user.Email).First((&dbuser))

//check if email is already register or not
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
