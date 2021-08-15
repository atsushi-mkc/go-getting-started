package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/api/controllers"
	"github.com/heroku/go-getting-started/api/models"
	"github.com/heroku/go-getting-started/db"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	"github.com/russross/blackfriday"
)

func repeatHandller(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

// func dbFunc(db *sql.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
// 			c.String(http.StatusInternalServerError,
// 				fmt.Sprintf("Error creating database table: %q", err))
// 			return
// 		}
// 		if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
// 			c.String(http.StatusInternalServerError,
// 				fmt.Sprintf("Error incrementing tick: %q", err))
// 			return
// 		}
// 		rows, err := db.Query("SELECT tick FROM ticks")
// 		if err != nil {
// 			c.String(http.StatusInternalServerError,
// 				fmt.Sprintf("Error reading ticks: %q", err))
// 			return
// 		}

// 		defer rows.Close()
// 		for rows.Next() {
// 			var tick time.Time
// 			if err := rows.Scan(&tick); err != nil {
// 				c.String(http.StatusInternalServerError,
// 					fmt.Sprintf("Error scanning ticks: %q", err))
// 				return
// 			}
// 			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
// 		}
// 	}
// }

type Tick struct {
	Tick time.Time
}

func gormDBFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tick = Tick{Tick: time.Now()}
		db := db.GetDB()
		db.Create(&tick)
		var ticks []Tick
		result := db.Find(&ticks)
		if result.Error != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", result.Error))
		}
		for _, t := range ticks {
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", t.Tick.String()))
		}
	}
}

const location = "Asia/Tokyo"

func initTimeLocation() {
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	initTimeLocation()
	db.Init()
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	// })

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.Run([]byte("**hi!**"))))
	})

	router.GET("/repeat", repeatHandller(repeat))

	router.GET("/db", gormDBFunc())

	db.GetDB().AutoMigrate(&models.User{})
	user := controllers.UserController{}
	router.POST("/api/user", user.Add)
	router.GET("/api/user", user.GetAll)
	router.GET("/api/user/:id", user.Get)
	router.PUT("/api/user/:id", user.Update)

	db.GetDB().AutoMigrate(&models.BasicAuthUser{})
	basicUser := controllers.BasicAuthUserController{}
	router.POST("/signup", basicUser.Signup)
	router.POST("/users/:id", basicUser.Get)

	router.Run(":" + port)
}
