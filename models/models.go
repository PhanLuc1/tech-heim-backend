package models

type User struct {
	Id         *int    `json:"id"`
	Email      *string `json:"email"`
	First_Name *string `json:"firstName"`
	Last_Name  *string `json:"lastName"`
	Password   *string `json:"Password"`
}
type Image struct {
	Url         *string `json:"imageURL"`
	Description *string `json:"description"`
}
type Category struct {
	Id_category   *int    `json:"idCategory"`
	Name_category *string `json:"nameCategory"`
	Desc_category *string `json:"descCategory"`
}
type Technical struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
type Type struct {
	Id          *int    `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
type Group struct {
	Id    *int    `json:"id"`
	Title *string `json:"title"`
	Type  []Type  `json:"productType"`
}
