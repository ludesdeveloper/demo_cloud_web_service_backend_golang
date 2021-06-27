package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	UserRequest struct {
		NIK     string `json:"nik"`
		Name    string `json:"name"`
		Company string `json:"company"`
	}
	User struct {
		gorm.Model
		NIK     string
		Name    string
		Company string
	}
)

func connectDB() *gorm.DB {
	// Connect to DB
	//dsn := "root:example@tcp(127.0.0.1:3306)/demo_cloud_web_service?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:example@tcp(databaseservice:3306)/demo_cloud_web_service?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func createUser(c echo.Context) error {
	userrequest := &UserRequest{}
	if err := c.Bind(userrequest); err != nil {
		return err
	}
	// Call Function connectDB
	db := connectDB()
	user := &User{}
	result := db.First(&user, "nik = ?", userrequest.NIK)
	// Checking Duplicate Value
	if result.RowsAffected > 0 {
		fmt.Println("Duplicate Detected")
		return c.JSON(http.StatusInternalServerError, "Duplicate NIK Detected !!!")
	} else {
		// Create DB
		db.Create(&User{NIK: userrequest.NIK, Name: userrequest.Name, Company: userrequest.Company})
	}
	// Close DB connection
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.Close()
	return c.JSON(http.StatusCreated, userrequest)
}

func getUser(c echo.Context) error {
	type NIK struct {
		NIK string `json:"nik"`
	}
	id := &NIK{}
	if err := c.Bind(id); err != nil {
		return err
	}
	if id.NIK == "" {
		// Call Function connectDB
		db := connectDB()
		// Read DB
		//var users []User
		users := &[]User{}
		db.Find(&users)
		// Close DB connection
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
		return c.JSON(http.StatusOK, users)
	} else {
		// Call Function connectDB
		db := connectDB()
		// Read DB
		user := &User{}
		db.First(&user, "nik = ?", id.NIK)
		// Close DB connection
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
		return c.JSON(http.StatusOK, user)
	}

}

func updateUser(c echo.Context) error {
	userrequest := &UserRequest{}
	if err := c.Bind(userrequest); err != nil {
		return err
	}
	// Call Function connectDB
	db := connectDB()
	// Update DB
	user := &User{}
	db.Model(&user).Where("nik = ?", userrequest.NIK).Updates(User{Company: userrequest.Company, Name: userrequest.Name})
	// Close DB connection
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.Close()
	return c.JSON(http.StatusOK, userrequest)
}

func deleteUser(c echo.Context) error {
	type NIK struct {
		NIK string `json:"nik"`
	}
	id := &NIK{}
	if err := c.Bind(id); err != nil {
		return err
	}
	// Call Function connectDB
	db := connectDB()
	// Read DB
	user := &User{}
	// Soft Delete
	db.Delete(&user, "nik = ?", id.NIK)
	// Permanently Delete
	// db.Unscoped().Delete(&user, "nik = ?", id.NIK)
	// Close DB connection
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.Close()
	return c.JSON(http.StatusOK, id)
}

func main() {
	// Call Function connectDB
	db := connectDB()
	// Generate Table
	db.AutoMigrate(&User{})
	// Close DB connection
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("DB Error")
	}
	sqlDB.Close()
	// Echo
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	// Routes
	e.POST("/users", createUser)
	e.GET("/users", getUser)
	e.PUT("/users", updateUser)
	e.DELETE("/users", deleteUser)
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
