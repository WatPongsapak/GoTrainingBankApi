package main

import (
	"strconv"
	"bankapi/user"
	"bankapi/secret"
	"bankapi/bankaccount"
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

type AccountApiService interface {
	AllAccounts() ([]bankaccount.Account, error)
	CreateAccount(acc *bankaccount.Account) error
	GetAccountByUserID(id int) ([]bankaccount.Account, error)
	GetAccountByID(id int) (*bankaccount.Account, error)
	DeleteAccount(id int) error
	UpdateAccount(id int, acc *bankaccount.Account) (*bankaccount.Account, error)
}

type SecretService interface {
	Insert(s *secret.Secret) error
	FindSecretKey(s *secret.Secret) error
}

type Server struct {
	userApiService UserApiService
	secretService  SecretService
	accountApiService  AccountApiService
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

	fmt.Println(user)

	if err := s.userApiService.CreateUser(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (s *Server) GetUserByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	acc, err := s.userApiService.GetUserByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, acc)
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
	fmt.Println(h)
	id, _ := strconv.Atoi(c.Param("id"))
	userinput,err := s.userApiService.GetUserByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(userinput)
	if h["firstname"] != ""{
		userinput.FirstName = h["firstname"]
	}
	if h["lastname"] != ""{
		userinput.LastName = h["lastname"]
	}
	fmt.Println(userinput)
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

func (s *Server) Auth(c *gin.Context) {
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

func (s *Server) CreateAccount(c *gin.Context) {
	var acc bankaccount.Account
	err := c.ShouldBindJSON(&acc)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"object":  "error",
			"message": fmt.Sprintf("json: wrong params: %s", err),
		})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	acc.UserID = int64(id)
	fmt.Println(acc)

	if err := s.accountApiService.CreateAccount(&acc); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, acc)
}

func (s *Server) GetAccountByUserID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	fmt.Println(id)
	acc, err := s.accountApiService.GetAccountByUserID(id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, acc)
}

func (s *Server) AccountWithdraw(c *gin.Context) {
	h := map[string]int{}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	acc, err := s.accountApiService.GetAccountByID(id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	acc.Balance = acc.Balance - int64(h["amount"]) 
	newacc,err := s.accountApiService.UpdateAccount(int(acc.ID),acc)
	c.JSON(http.StatusOK, newacc)
}

func (s *Server) AccountDeposit(c *gin.Context) {
	h := map[string]int{}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	acc, err := s.accountApiService.GetAccountByID(id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	acc.Balance = acc.Balance + int64(h["amount"])  
	newacc,err := s.accountApiService.UpdateAccount(int(acc.ID),acc)
	c.JSON(http.StatusOK, newacc)
}

func (s *Server) DeleteAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.accountApiService.DeleteAccount(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) Transfer(c *gin.Context) {
	h := map[string]int{}
	if err := c.ShouldBindJSON(&h); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	accfrom, err := s.accountApiService.GetAccountByID(h["from"])
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	accto, err := s.accountApiService.GetAccountByID(h["to"])
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	accfrom.Balance = accfrom.Balance - int64(h["amount"])
	accto.Balance = accto.Balance + int64(h["amount"])    
	newaccfrom,err := s.accountApiService.UpdateAccount(int(accfrom.ID),accfrom)
	newaccto,err := s.accountApiService.UpdateAccount(int(accto.ID),accto)
	c.JSON(http.StatusOK, []*bankaccount.Account{newaccfrom, newaccto})
}


func setupRoute(s *Server) *gin.Engine {
	r := gin.Default()
	// accs := r.Group("/accs")
	admin := r.Group("/admin")
	user := r.Group("/users")
	account := r.Group("/bankAccounts")
	admin.Use(gin.BasicAuth(gin.Accounts{
		"admin": "1234",
	}))
	admin.POST("/secrets", s.CreateSecret)
	user.Use(s.Auth)
	account.Use(s.Auth)
	user.GET("", s.AllUsers)
	user.POST("", s.CreateUser)
	user.GET("/:id", s.GetUserByID)
	user.PUT("/:id", s.UpdateUser)
	user.DELETE("/:id", s.DeleteUser)
	user.POST("/:id/bankAccounts", s.CreateAccount)
	user.GET("/:id/bankAccounts", s.GetAccountByUserID)
	account.PUT("/:id/withdraw", s.AccountWithdraw)
	account.PUT("/:id/deposit", s.AccountDeposit)
	account.DELETE("/:id", s.DeleteAccount)

	return r
}

func main() {
	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", "postgres://suyhzbwz:zMMdsNufLoJGLzdVphQt9qb6pwjI02Wu@elmer.db.elephantsql.com:5432/suyhzbwz")
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		userApiService: &user.UserApiServiceImp{
			DB: db,
		},
		secretService: &secret.SecretServiceImp{
			DB: db,
		},
		accountApiService: &bankaccount.AccountApiServiceImp{
			DB: db,
		},
	}

	r := setupRoute(s)

	r.Run(":" + os.Getenv("PORT"))
}
