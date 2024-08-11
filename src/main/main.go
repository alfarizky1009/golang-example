package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
);

func main() {
	fmt.Print("Welcome to the server!")

	//Create Echo object
	e := echo.New()

	g := e.Group("/admin")

	//How to setup middleware
	// middleware.Logger() used to loged the server interaction

	//1. Put it inside Grouping
	//g := e.Group("/admin", middleware.Logger())

	//2. Put it inside Method only for spesific use
	// g.GET("/main", mainAdmin, middleware.Logger() )

	//3. Using Use function **Recommended method
	// g.Use(middleware.Logger())

	//Middleware custom logger
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	g.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		// Check in DB

		if username == "alfa1" && password == "password" {
			return true, nil
		}

		return false, nil
	}))

	g.GET("/main", mainAdmin)

	e.GET("/", halloWeb)
	e.GET("/cats/:data", getCats)
	
	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/hamster", addHamster)

	e.Start(":8000")
}

func halloWeb(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from the web side")
}

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}



func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	
	dataType := c.Param("data")
	
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is: %s\nyour cat type is: %s", catName, catType))
	}

	//map[key]value{}
	if (dataType == "json") {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string {
		"error": "Error type not described",
	})	
}

// Faster method
func addCat(c echo.Context) error {
	cat := Cat{} 

	//defer is called after all the return function is executed. Just like finally
	defer c.Request().Body.Close()
	
	//b = body, err = error. Just like try catch
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading request body for addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	//Marshaling is convert Go Data into Json Data.
	//Unmarshaling is convert Json Data to Go Data.
	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	// %#v prints a Go syntax representation of the value, i.e. the source code snippet that would produce that value.
	log.Printf("This is your cat: %#v", cat)
	return c.String(http.StatusOK, "We got your cat!")
}

// Recommended for project 
func addDog(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()

	// Create new decoder from request body and decode to dog class
	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed reading request body for addDog: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("This is your dog: %#v", dog)
	return c.String(http.StatusOK, "We got your dog!")
}

// Slowest Method
func addHamster(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("Failed reading request body for addHamster: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("This is your hamster: %#v", hamster)
	return c.String(http.StatusOK, "We got your hamster!")
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "You are at secret admin page")
}