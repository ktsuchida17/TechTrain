package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ktsuchida17/TechTrain/pkg/db/mysql"
	"github.com/ktsuchida17/TechTrain/pkg/model"
)

type Server struct {
	DB     *gorm.DB
	engine *gin.Engine
}

type result map[string]interface{}

func NewServer() (*Server, *model.Error) {
	DB, err := mysql.ConnectToDB()
	if err != nil {
		Error := &model.Error{
			Title:  "Error establishing the Database Connection",
			Code:   http.StatusInternalServerError,
			Msg:    err.Error(),
			Detail: "Could not establish the Database Connection. Check the configuration file.",
		}
		return nil, Error
	}
	return &Server{
		DB:     DB,
		engine: gin.Default(),
	}, nil
}

func main() {
	Server, Error := NewServer()
	if Error != nil {
		panic(*Error)
	}
	defer Server.DB.Close()

	Server.engine.POST("/user/create", func(c *gin.Context) {
		var req model.Request
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.Error{
				Title:  "Invalid Request Format",
				Code:   http.StatusBadRequest,
				Msg:    err.Error(),
				Detail: "The format of the request body is JSON.",
			})
			return
		}

		Server.DB.Exec("LOCK TABLES `users` WRITE")
		defer Server.DB.Exec("UNLOCK TABLES")
		ID, err := mysql.GenerateID(Server.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Error{
				Title:  "Error generating userID",
				Code:   http.StatusInternalServerError,
				Msg:    err.Error(),
				Detail: "Failed to generate usercharacterID of yours. Read the GORM documentation.",
			})
			return
		}
		token, err := mysql.GenerateToken(Server.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Error{
				Title:  "Error generating x-token",
				Code:   http.StatusInternalServerError,
				Msg:    err.Error(),
				Detail: "Failed to generate x-token of yours.",
			})
			return
		}
		user := model.User{
			ID:    ID,
			Name:  req.Name,
			Token: token,
		}
		if err := mysql.CreateUser(Server.DB, user); err != nil {
			c.JSON(http.StatusInternalServerError, model.Error{
				Title:  "Error creating user data",
				Code:   http.StatusInternalServerError,
				Msg:    err.Error(),
				Detail: "Failed to create user data of yours. Read the GORM documentation.",
			})
			return
		}
		c.JSON(http.StatusOK, result{
			"token": user.Token,
		})
	})

	Server.engine.GET("/user/get", func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusBadRequest, model.Error{
				Title:  "Invalid Request Format",
				Code:   http.StatusBadRequest,
				Msg:    "token not found",
				Detail: "The request header must contain a x-token. Put the token in the header.",
			})
			return
		}

		name, err := mysql.GetUserName(Server.DB, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Error{
				Title:  "Error getting user name",
				Code:   http.StatusInternalServerError,
				Msg:    err.Error(),
				Detail: "Failed to get user name of your data. Check the x-token.",
			})
			return
		}
		c.JSON(http.StatusOK, result{
			"name": name,
		})
	})

	Server.engine.PUT("/user/update", func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusBadRequest, model.Error{
				Title:  "Invalid Request Format",
				Code:   http.StatusBadRequest,
				Msg:    "token not found",
				Detail: "The request header must contain a x-token. Put the token in the header.",
			})
			return
		}
		var req model.Request
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.Error{
				Title:  "Invalid Request Format",
				Code:   http.StatusBadRequest,
				Msg:    err.Error(),
				Detail: "The format of the request body is JSON.",
			})
			return
		}

		if err := mysql.UpdateUserName(Server.DB, token, req.Name); err != nil {
			c.JSON(http.StatusInternalServerError, model.Error{
				Title:  "Error updating user name",
				Code:   http.StatusInternalServerError,
				Msg:    err.Error(),
				Detail: "Failed to update user name of your data. Check the x-token.",
			})
			return
		}
		c.Status(http.StatusOK)
	})

	Server.engine.POST("gacha/draw", func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusBadRequest, model.Error{
				Title:  "Invalid Request Format",
				Code:   http.StatusBadRequest,
				Msg:    "token not found",
				Detail: "The request header must contain a x-token. Put the token in the header.",
			})
			return
		}
		var req model.Request
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.Error{
				Title:  "Invalid Request Format",
				Code:   http.StatusBadRequest,
				Msg:    err.Error(),
				Detail: "The format of the request body is JSON.",
			})
			return
		}

		ID, err := mysql.GetUserID(Server.DB, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Error{
				Title:  "Error getting user characterID",
				Code:   http.StatusInternalServerError,
				Msg:    err.Error(),
				Detail: "Failed to get user characterID of your data. Check the x-token.",
			})
			return
		}

		results := []model.Character{}
		for i := 0; i < req.Times; i++ {
			character, err := mysql.Gacha(Server.DB, token)
			if err != nil {
				c.JSON(http.StatusInternalServerError, model.Error{
					Title:  "Gacha Error",
					Code:   http.StatusInternalServerError,
					Msg:    err.Error(),
					Detail: "Failed to get a character by rolling the gacha.",
				})
				return
			}
			if err := mysql.SaveGachaResults(Server.DB, ID, character); err != nil {
				c.JSON(http.StatusInternalServerError, model.Error{
					Title:  "Error saving Gacha results",
					Code:   http.StatusInternalServerError,
					Msg:    err.Error(),
					Detail: "Failed to save gacha results. ",
				})
			}
			results = append(results, character)
		}

		c.JSON(http.StatusOK, result{
			"results": results,
		})
	})

	Server.engine.GET("character/list", func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusBadRequest, model.Error{
				Title:  "Invalid Request Format",
				Code:   http.StatusBadRequest,
				Msg:    "token not found",
				Detail: "The request header must contain a x-token. Put the token       in the header.",
			})
			return
		}

		ID, err := mysql.GetUserID(Server.DB, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Error{
				Title:  "Error getting user characterID",
				Code:   http.StatusInternalServerError,
				Msg:    err.Error(),
				Detail: "Failed to get user characterID of your data. Check the x-      token.",
			})
			return
		}

		list, err := mysql.GetUsersCharacterList(Server.DB, ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Error{
				Title:  "Error getting user's character list",
				Code:   http.StatusInternalServerError,
				Msg:    err.Error(),
				Detail: "Failed to get user's character list of your data. Check the x-token.",
			})
			return
		}
		c.JSON(http.StatusOK, result{
			"characters": list,
		})
	})

	Server.engine.Run(":8080")
}
