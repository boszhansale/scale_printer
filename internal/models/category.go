package models

type CategoriesWithProduct struct {
	Name     string    `json:"name"`
	Products []Product `json:"label_products"`
}
