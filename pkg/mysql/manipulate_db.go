package mysql

import (
	"../data"
	"crypto/rand"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"math/big"
)

func Connect() *gorm.DB {
	DBMS := "mysql"
	user := "docker"
	password := "docker"
	host := "localhost"
	port := "1217"
	database_name := "db_TechTrain"

	dbconfig := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database_name + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	db, err := gorm.Open(DBMS, dbconfig)
	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(&data.User{})
	return db
}

func GenerateID() uint {
	db := Connect()
	user := new(data.User)
	if result := db.Order("id desc").Take(&user); result.Error != nil {
		return 1
	}
	return user.Id + 1
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Generatoken() string {
	db := Connect()
	user := new(data.User)
	var x_token string
	token := make([]byte, 15)
	length := int64(len(letters))
	for {
		for i := range token {
			n, err := rand.Int(rand.Reader, big.NewInt(length))
			if err != nil {
				panic("failed generate random token")
			}
			token[i] = letters[n.Int64()]
			x_token = string(token)
		}
		if result := db.Where("token = ?", x_token).Find(&user); result.Error != nil {
			break
		}
	}
	return x_token
}

func Create(info data.User) bool {
	db := Connect()
	if result := db.Create(&info); result.Error != nil {
		return false
	}
	db.Close()
	return true
}

func Get(token string) string {
	db := Connect()
	user := new(data.User)
	if err := db.Where("token = ?", token).Find(&user).Error; err != nil {
		return ""
	}
	db.Close()
	return user.Name
}

func Update(token string, name string) bool {
	db := Connect()
	user := new(data.User)
	if err := db.Where("token = ?", token).Find(&user).Error; err != nil {
		return false
	}
	updated_user := data.User{
		Id:    user.Id,
		Name:  name,
		Token: token,
	}
	if err := db.Model(&user).Updates(updated_user).Error; err != nil {
		return false
	}
	db.Close()
	return true
}
