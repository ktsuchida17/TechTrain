package main

import (
	"../pkg/data"
	"../pkg/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	engine := gin.Default()

	engine.POST("/user/create", func(c *gin.Context) {
		var form data.User
		if err := c.BindJSON(&form); err != nil {
			c.String(http.StatusBadRequest, "Bad request\n")
			return
		}
		user := data.User{
			Id:    mysql.GenerateID(),
			Name:  form.Name,
			Token: mysql.Generatoken(),
		}
		if !mysql.Create(user) {
			c.String(http.StatusBadRequest, "Failed to create your data\n")
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"token": user.Token,
			})
		}
	})
	engine.GET("/user/get", func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if token == "" {
			c.String(http.StatusBadRequest, "Bad request\n")
			return
		}
		name := mysql.Get(token)
		if name == "" {
			c.String(http.StatusBadRequest, "Failed to search your name\n")
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"name": name,
			})
		}
	})
	engine.PUT("/user/update", func(c *gin.Context) {
		token := c.GetHeader("x-token")
		var form data.User
		if err := c.BindJSON(&form); token == "" || err != nil {
			c.String(http.StatusBadRequest, "Bad request\n")
			return
		}
		if !mysql.Update(token, form.Name) {
			c.String(http.StatusBadRequest, "Failed to change your name\n")
			return
		} else {
			c.Status(http.StatusOK)
		}
	})
	engine.Run(":8080")
}
