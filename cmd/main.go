package main

import (
	"../pkg/model"
	"../pkg/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	engine := gin.Default()

	engine.POST("/user/create", func(c *gin.Context) {
		test, err := mysql.ConnectToDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Err{
				Msg: "Could not establish a Database Connection. Check the configuration file.",
				Err: err.Error(),
			})
			return
		}
		test.Close()

		var req model.User
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.Err{
				Msg: "Invalid Request Format. The format of the request body is JSON.",
				Err: err.Error(),
			})
			return
		}
		id, err := mysql.GenerateID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Err{
				Msg: "Failed to generate userID of yours.",
				Err: err.Error(),
			})
			return
		}
		token, err := mysql.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Err{
				Msg: "Failed to generate x-token of yours.",
				Err: err.Error(),
			})
			return
		}
		user := model.User{
			Id:    id,
			Name:  req.Name,
			Token: token,
		}
		if err := mysql.CreateUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, model.Err{
				Msg: "Failed to create user data of yours.",
				Err: err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"token": user.Token,
			})
		}
	})

	engine.GET("/user/get", func(c *gin.Context) {
		test, err := mysql.ConnectToDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Err{
				Msg: "Could not establish a Database Connection. Check the configuration file.",
				Err: err.Error(),
			})
			return
		}
		test.Close()

		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusBadRequest, model.Err{
				Msg: "Invalid Request Format. Put the x-token in the request header.",
				Err: "No token in the request header",
			})
			return
		}
		name, err := mysql.GetUserInfo(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Err{
				Msg: "Failed to find user data of yours. Check the x-token.",
				Err: err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"name": name,
			})
		}
	})

	engine.PUT("/user/update", func(c *gin.Context) {
		test, err := mysql.ConnectToDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Err{
				Msg: "Could not establish a Database Connection. Check the configuration file.",
				Err: err.Error(),
			})
			return
		}
		test.Close()

		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusBadRequest, model.Err{
				Msg: "Invalid Request Format. Put the x-token in the request header.",
				Err: "No token in the request header",
			})
			return
		}
		var req model.User
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.Err{
				Msg: "Invalid Request Format. The format of the request body is JSON.",
				Err: err.Error(),
			})
			return
		}
		if err := mysql.UpdateUserInfo(token, req.Name); err != nil {
			c.JSON(http.StatusInternalServerError, model.Err{
				Msg: "Failed tp Update user data of yours.",
				Err: err.Error(),
			})
		} else {
			c.Status(http.StatusOK)
		}
	})

	engine.Run(":8080")
}
