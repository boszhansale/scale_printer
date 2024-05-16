package repository

import (
	"encoding/json"
	"errors"
	"log"
	"test/internal/models"
)

type DB struct {
	Data []models.CategoriesWithProduct
}

func New(str []byte) *DB {
	var data []models.CategoriesWithProduct
	err := json.Unmarshal(str, &data)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return nil
	}
	return &DB{data}
}

func (db *DB) GetCategoryNames() []string {
	names := make([]string, len(db.Data))

	for i, category := range db.Data {
		names[i] = category.Name
	}
	return names
}

func (db *DB) GetProductNames(categoryName string, lang string) []string {
	for _, category := range db.Data {
		if category.Name == categoryName {
			productNames := make([]string, len(category.Products))
			for i, product := range category.Products {
				if lang == "kz" {
					productNames[i] = product.NameKz
				} else {
					productNames[i] = product.NameEn
				}
			}
			return productNames
		}
	}
	return nil
}

func (db *DB) GetProduct(categoryName, productName, lang string) (models.Product, error) {

	if lang == "kz" {
		for _, category := range db.Data {
			if category.Name == categoryName {
				for _, product := range category.Products {
					if product.NameKz == productName {
						return product, nil
					}
				}
			}
		}
	} else {
		for _, category := range db.Data {
			if category.Name == categoryName {
				for _, product := range category.Products {
					if product.NameEn == productName {
						return product, nil
					}
				}
			}
		}
	}

	return models.Product{}, errors.New("product not found")
}
