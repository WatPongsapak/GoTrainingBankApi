package main

import (
	"strconv"
	"bankapi/user"
	"bankapi/secret"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type UserApiService interface {
	AllUsers() ([]user.User, error)
	CreateUser(user *user.User) error
	GetUserByID(id int) (*user.User, error)
	DeleteUser(id int) error
	UpdateUser(id int, user *user.User) (*user.User, error)
}


type SecretService interface {
	Insert(s *secret.Secret) error
	FindSecretKey(s *secret.Secret) error
}

type Server struct {
	userApiService UserApiService
	secretService  SecretService
}

func (s *Server) AllUsers(c *gin.Context) {
	users, err := s.userApiService.AllUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("db: query error: %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (s *Server) CreateUser(c *gin.Context) {
	var user user.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("json: wrong params: %s", err),
		})
		return
	}

	if err := s.userApiService.CreateUser(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (s *Server) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	todo, err := s.userApiService.GetUserByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (s *Server) DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.userApiService.DeleteUser(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) UpdateUser(c *gin.Context) {
	h := map[string]string{}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	userinput := &user.User{
		ID        :int64(id),
		FirstName :h["firstname"],
		LastName  :h["lastname"],
	}
	user, err := s.userApiService.UpdateUser(id, userinput)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *Server) CreateSecret(c *gin.Context) {
	var secret secret.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	if err := s.secretService.Insert(&secret); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, secret)
}

func (s *Server) AuthTodo(c *gin.Context) {
	var secret secret.Secret
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	if err := s.secretService.FindSecretKey(&secret); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
}

func setupRoute(s *Server) *gin.Engine {
	r := gin.Default()
	// todos := r.Group("/todos")
	// admin := r.Group("/admin")

	// admin.Use(gin.BasicAuth(gin.Accounts{
	// 	"admin": "1234",
	// }))
	// r.Use(s.AuthTodo)
	r.GET("/users", s.AllUsers)
	r.POST("/users", s.CreateUser)
	r.GET("/users/:id", s.GetUserByID)
	r.PUT("/users/:id", s.UpdateUser)
	r.DELETE("/users/:id", s.DeleteUser)

	// r.POST("/users/:id/bankAccounts", s.Create)
	// r.GET("/users/:id/bankAccounts", s.GetByID)
	// r.PUT("/users/:id/bankAccounts", s.GetByID)
	// r.PUT("/users/:id/bankAccounts", s.GetByID)
	// r.DELETE("/users/:id/bankAccounts", s.GetByID)

	return r
}
func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		userApiService: &user.UserApiServiceImp{
			db: db,
		},
		secretService: &secret.SecretServiceImp{
			db: db,
		},
	}

	r := setupRoute(s)

	r.Run(":" + os.Getenv("PORT"))
}
