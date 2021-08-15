package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/api/models"
	"github.com/heroku/go-getting-started/db"
)

type UserController struct{}

func (u UserController) Add(c *gin.Context) {
	user := models.User{}
	now := time.Now()
	user.UpdatedAt = now
	user.CreatedAt = now

	err := c.BindJSON(&user)
	if err != nil {
		c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		return
	}
	db := db.GetDB()
	result := db.Create(&user)
	if result.Error != nil {
		c.String(http.StatusInternalServerError, "Error: "+result.Error.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u UserController) GetAll(c *gin.Context) {
	users := []models.User{}
	db := db.GetDB()
	result := db.Find(&users)
	if result.Error != nil {
		c.String(http.StatusInternalServerError, "Error: "+result.Error.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

func (u UserController) Get(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	db := db.GetDB()
	user := models.User{}
	result := db.Where("ID = ?", uid).First(&user)

	if result.Error != nil {
		c.String(http.StatusInternalServerError, "Error: "+result.Error.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u UserController) Update(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	data := models.User{}
	err := c.BindJSON(&data)
	if err != nil {
		c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		return
	}
	db := db.GetDB()
	user := models.User{}
	result := db.Where("ID = ?", uid).First(&user).Updates(&data)

	if result.Error != nil {
		c.String(http.StatusInternalServerError, "Error: "+result.Error.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}
