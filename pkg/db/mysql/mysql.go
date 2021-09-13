package mysql

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/ktsuchida17/TechTrain/pkg/model"
)

func ConnectToDB() (*gorm.DB, error) {
	if err := godotenv.Load("envfiles/TechTrain.env"); err != nil {
		return nil, err
	}
	DBMS := os.Getenv("DBMS")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	DBname := os.Getenv("DB_NAME")
	dbconfig := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + DBname + "?charset=utf8mb4&parseTime=true&loc=Asia%2FTokyo"

	DB, err := gorm.Open(DBMS, dbconfig)
	if err != nil {
		return nil, err
	}
	DB.AutoMigrate(&model.User{})
	DB.Table("users_character_list").AutoMigrate(&model.UsersCharacterList{})
	return DB, nil
}

func GenerateID(DB *gorm.DB) (uint, error) {
	user := new(model.User)
	if err := DB.Order("id desc").Take(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return 1, nil
		} else {
			return 0, err
		}
	}
	return user.ID + 1, nil
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateToken(DB *gorm.DB) (string, error) {
	var xtoken string
	user := new(model.User)
	token := make([]byte, 5)
	length := int64(len(letters))
	for {
		for i := range token {
			n, err := rand.Int(rand.Reader, big.NewInt(length))
			if err != nil {
				return "", err
			}
			token[i] = letters[n.Int64()]
		}
		xtoken = string(token)
		if err := DB.Where("token = ?", xtoken).Order("id desc").Take(&user).Error; err != nil {
			if err.Error() == "record not found" {
				break
			} else {
				return "", err
			}
		}
	}
	return xtoken, nil
}

func CreateUser(DB *gorm.DB, user model.User) error {
	if err := DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserName(DB *gorm.DB, token string) (string, error) {
	user := new(model.User)
	if err := DB.Where("token = ?", token).Find(&user).Error; err != nil {
		return "", err
	}
	return user.Name, nil
}

func UpdateUserName(DB *gorm.DB, token string, name string) error {
	user := new(model.User)
	if err := DB.Where("token = ?", token).Find(&user).Error; err != nil {
		return err
	}
	user.Name = name
	user.Token = token
	if err := DB.Model(&user).Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func Gacha(DB *gorm.DB, token string) (model.Character, error) {
	list := []model.Gacha{}
	if err := DB.Table("gacha").Find(&list).Error; err != nil {
		return model.Character{}, err
	}
	var sum float64
	sum = 0
	for _, character := range list {
		sum += character.Weight
	}
	if sum != 100 {
		return model.Character{}, errors.New("the sum of weights in the gacha didn't add up to 100")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return model.Character{}, err
	}
	decimal := float64(n.Int64()) / (1 << 63)
	n, err = rand.Int(rand.Reader, big.NewInt(99))
	if err != nil {
		return model.Character{}, err
	}
	number := float64(n.Int64()) + decimal
	if !(number >= 0 && number <= 100) {
		return model.Character{}, errors.New("error generating the random number")
	}

	sum = 0
	for _, character := range list {
		sum += float64(character.Weight)
		if sum >= number {
			return model.Character{
				CharacterID: character.CharacterID,
				Name:        character.Name,
			}, nil
		}
	}
	return model.Character{}, errors.New("logical error in the process of the gacha")
}

func GetUserID(DB *gorm.DB, token string) (uint, error) {
	user := new(model.User)
	if err := DB.Where("token = ?", token).Find(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func SaveGachaResults(DB *gorm.DB, ID uint, character model.Character) error {
	list := model.UsersCharacterList{
		UserCharacterID: fmt.Sprint(ID),
		CharacterID:     fmt.Sprint(character.CharacterID),
		Name:            character.Name,
	}
	if err := DB.Table("users_character_list").Create(&list).Error; err != nil {
		return err
	}
	return nil
}

func GetUsersCharacterList(DB *gorm.DB, ID uint) ([]model.UsersCharacterList, error) {
	list := []model.UsersCharacterList{}
	if err := DB.Table("users_character_list").Where("user_character_id = ?", ID).Find(&list).Error; err != nil {
		return []model.UsersCharacterList{}, err
	}
	return list, nil
}
