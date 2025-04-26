package models

type Horse struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"not null;" binding:"required,min=3"`
	Breed  string `json:"breed"`
	Age    uint   `json:"age" binding:"omitempty,min=1"`
	Gender string `json:"gender" binding:"omitempty,oneof=m f"`
	Owner  *uint  `json:"owner" gorm:"not null;"`
}

type RoleEnum string

const (
	ADMIN RoleEnum = "ADMIN"
	USER  RoleEnum = "USER"
)

type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"unique;not null;"`
}

type User struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Username string  `json:"username" gorm:"not null;" binding:"required,min=3"`
	Mail     string  `json:"mail" gorm:"unique;not null;" binding:"required,email"`
	Age      uint    `json:"age" binding:"omitempty,min=1"`
	Gender   string  `json:"gender" binding:"omitempty,oneof=m f"`
	Password string  `json:"password" gorm:"not null;" binding:"required"`
	Country  string  `json:"country" gorm:"not null;" binding:"required,iso3166_1_alpha2"`
	Roles    []Role  `json:"roles" gorm:"many2many:user_roles; constraint:OnDelete:SET NULL;"`
	Horses   []Horse `json:"horses" gorm:"foreignKey:Owner; constraint:OnDelete:SET NULL;"`
}

type Group struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" gorm:"not null;" binding:"required,min=3"`
	Users []User `json:"users" gorm:"many2many:group_users; constraint:OnDelete:SET NULL;"`
}
