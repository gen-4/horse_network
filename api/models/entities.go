package models

type Horse struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	Age    uint   `json:"age"`
	Gender string `json:"gender"`
	Owner  *uint  `json:"owner"`
}

type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type User struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	Username string  `json:"username"`
	Mail     string  `json:"mail" gorm:"uniqueIndex:compositeindex"`
	Age      uint    `json:"age"`
	Gender   string  `json:"gender"`
	Password string  `json:"password"`
	Country  string  `json:"country"`
	Roles    []Role  `json:"roles" gorm:"many2many:user_roles; constraint:OnDelete:SET NULL;"`
	Horses   []Horse `json:"horses" gorm:"foreignKey:Owner; constraint:OnDelete:SET NULL;"`
}
