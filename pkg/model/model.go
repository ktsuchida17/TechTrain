package model

type User struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name" gorm:"not null"`
	Token string `json:"token" gorm:"not null"`
}

type Character struct {
	CharacterID uint   `gorm:"->;column:characterID"`
	Name        string `gorm:"->"`
}

type Gacha struct {
	CharacterID uint    `gorm:"->;column:characterID"`
	Name        string  `gorm:"->"`
	Weight      float64 `gorm:"->"`
}

type UsersCharacterList struct {
	UserCharacterID string `gorm:"not null"`
	CharacterID     string `gorm:"not null"`
	Name            string `gorm:"not null"`
}

type Request struct {
	Name  string `json:"name"`
	Times int    `json:"times"`
}

type Error struct {
	Title  string
	Code   uint
	Msg    string
	Detail string
}
