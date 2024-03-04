package models

type User struct {
	Id         *int    `json:"id"`
	Email      *string `json:"email"`
	First_Name *string `json:"firstName"`
	Last_Name  *string `json:"lastName"`
	Password   *string `json:"Password"`
}
