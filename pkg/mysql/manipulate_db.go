package mysql

import (
	"../model"
	"crypto/rand"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"math/big"
	"os"
)

func ConnectToDB() (*gorm.DB, error) {
	if err := godotenv.Load("../envfiles/TechTrain.env"); err != nil {
		return nil, err
	}
	DBMS := os.Getenv("DBMS")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	dbconfig := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + db_name + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	db, err := gorm.Open(DBMS, dbconfig)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&model.User{})
	return db, nil
}

func GenerateID() (uint, error) {
	db, _ := ConnectToDB()
	defer db.Close()

	user := new(model.User)
	if err := db.Order("id desc").Take(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return 1, nil
		} else {
			return 0, err
		}
	}
	return user.Id + 1, nil
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" //62 characters
func GenerateToken() (string, error) {
	db, _ := ConnectToDB()
	defer db.Close()

	var x_token string
	user := new(model.User)
	token := make([]byte, 5) //62 ** 5 is about 1 billion.
	length := int64(len(letters))
	for {
		for i := range token {
			n, err := rand.Int(rand.Reader, big.NewInt(length))
			if err != nil {
				return "", err
			}
			token[i] = letters[n.Int64()]
			x_token = string(token)
		}
		if err := db.Where("token = ?", x_token).Order("id desc").Take(&user).Error; err != nil {
			if err.Error() == "record not found" {
				break
			} else {
				return "", err
			}
		}
	}
	return x_token, nil
}

func CreateUser(info model.User) error {
	db, _ := ConnectToDB()
	defer db.Close()

	if err := db.Create(&info).Error; err != nil {
		return err
	}
	return nil
}

func GetUserInfo(token string) (string, error) {
	db, _ := ConnectToDB()
	defer db.Close()

	user := new(model.User)
	if err := db.Where("token = ?", token).Find(&user).Error; err != nil {
		return "", err
	}
	return user.Name, nil
}

func UpdateUserInfo(token string, name string) error {
	db, _ := ConnectToDB()
	defer db.Close()

	user := new(model.User)
	if err := db.Where("token = ?", token).Find(&user).Error; err != nil {
		return err
	}
	updated_user := model.User{
		Id:	user.Id,
		Name:	name,
		Token:	token,
	}
	if err := db.Model(&user).Updates(updated_user).Error; err != nil {
		return err
	}
	return nil
}
