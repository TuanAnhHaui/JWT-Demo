package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	router          = gin.Default()
	secretKey       = os.Getenv("JWT_SECRET")
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	Payload TokenPayload `json:"payload"`
	jwt.StandardClaims
}
type TokenPayload struct {
	Id       int    `json:"id"`
	UserName string `json:"user_name"`
}

func main() {
	router.POST("/login", Login)
	router.GET("/verify", VerifyToken)
	log.Fatal(router.Run(":8080"))
}

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//A sample use
var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
func VerifyToken(c *gin.Context) {
	//h := model.AuthHeader{}
	var user User
	header := c.GetHeader("Authorization")
	//if err := c.ShouldBindHeader(&h); err != nil {
	//	fmt.Println(err)
	//}
	tokenHeader := strings.Split(header, " ")
	if tokenHeader[0] != "Bearer" || len(tokenHeader) < 2 || strings.TrimSpace(tokenHeader[1]) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Authorization",
		})
		return
	}
	validate, err := user.VerifyToken(tokenHeader[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid token",
		})
		return
	}
	c.JSON(http.StatusOK, validate)
}

func CreateToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Second * 60).Unix()
	td.RtExpires = time.Now().Add(time.Minute * 15).Unix()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	os.Setenv("REFRESH_TOKEN", "adfsdfsdfasd")
	rtClaims := jwt.MapClaims{}
	rtClaims["authorized"] = true
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_TOKEN")))
	if err != nil {
		return nil, err

	}
	return td, nil
}
func (u User) VerifyToken(token string) (*TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	clamis, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &clamis.Payload, nil
}
