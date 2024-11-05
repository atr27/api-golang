package main

import (
	"auth/auth"
	"auth/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type ListStudent struct {
	Student_id       uint64 `json:"student_id" binding:"required"`
	Student_name     string `json:"student_name" binding:"required"`
	Student_age      uint64 `json:student_age" binding:"required"`
	Student_address  string `json:"student_address" binding:"required"`
	Student_phone_no string `json:"student_phone_no" binding:"required"`
}

func postHandler(c *gin.Context, db *gorm.DB) {
	var newStudent ListStudent
	c.Bind(&newStudent)
	db.Create(&newStudent)
	c.JSON(http.StatusOK, gin.H{
		"message": "success created student", 
		"data": newStudent,
	})

}

func getAllHandler(c *gin.Context, db *gorm.DB) {
	var newStudent []ListStudent
	db.Find(&newStudent)
	c.JSON(http.StatusOK, gin.H{"message":"success get all", "data": newStudent})

}

func getHandler(c *gin.Context, db *gorm.DB) {
	var newStudent ListStudent
	studentId := c.Param("student_id")
	id, _ := strconv.ParseUint(studentId, 10, 64)

	if db.Find(&newStudent, "student_id = ?", id).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "no student data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success get data student", "data": newStudent})
}

func putHandler(c *gin.Context, db *gorm.DB) {
	var newStudent = ListStudent{}
	//integer to string param
	studentId := c.Param("student_id")

	if db.Find(&newStudent, "student_id = ?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{"message": "no student"})
		return
	}

	var reqStudent = newStudent

	c.Bind(&reqStudent)

	db.Model(&newStudent).Update(reqStudent)

	c.JSON(http.StatusOK, gin.H{"message": "success updated student", "data": reqStudent})
	
}

func deleteHandler(c *gin.Context, db *gorm.DB) {
	var newStudent ListStudent
	studenId := c.Param("student_id")
	db.Delete(&newStudent, "student_id = ?", studenId)
	c.JSON(http.StatusOK, gin.H{"message": "success deleted student"})
}


func setupRouter() *gin.Engine {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	conn := os.Getenv("POSTGRES_URL")
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})

	router.POST("/login", auth.LoginHandler)

	router.POST("/student", func(c *gin.Context) {
		postHandler(c, db)
	})

	router.GET("/student", middleware.AuthValidator, func(c *gin.Context) {
		getAllHandler(c, db)
	})

	router.GET("/student/:student_id", middleware.AuthValidator, func(c *gin.Context) {
		getHandler(c, db)
	})

	router.PUT("/student/:student_id", func (c *gin.Context){
		putHandler(c, db)
	})

	router.DELETE("/student/:student_id", func (c *gin.Context){
		deleteHandler(c, db)
	})

	return router
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&ListStudent{})

	data := ListStudent{}
	if db.Find(&data).RecordNotFound() {
		fmt.Println("====run seeder user====")
		seederUser(db)
	}
}

func seederUser(db *gorm.DB) {
	data := ListStudent {
		Student_id: 1,
		Student_name: "Budi",
		Student_age: 22,
		Student_address: "Bandung",
		Student_phone_no: "081443322111",	
	}

	db.Create(&data)
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
