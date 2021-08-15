package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/api/models"
	"github.com/heroku/go-getting-started/db"
)

type BasicAuthUserController struct{}

func (u BasicAuthUserController) Signup(c *gin.Context) {
	user := models.BasicAuthUser{}
	now := time.Now()
	user.UpdatedAt = now
	user.CreatedAt = now
	err := c.BindJSON(&user)
	if err != nil || !user.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account creation failed",
			"cause":   "required user_id and password"})
		return
	}
	db := db.GetDB()
	tempUser := models.BasicAuthUser{}
	result := db.Where("ID = ?", user.ID).Find(&tempUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account creation failed",
			"cause":   "already same user_id is used"})
		return
	}
	user.Nickname = &user.ID
	db.Create(&user)

	resUser := models.ResponseUser{ID: user.ID, Nickname: *user.Nickname}
	resMesage := models.ResponseMessage{
		Message: "Account successfully created",
		User:    resUser,
	}
	c.JSON(http.StatusOK, resMesage)
}
func (u BasicAuthUserController) Get(c *gin.Context) {

}
