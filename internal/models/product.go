package models

type Product struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Barcode     string `json:"barcode"`
	Weight      string `json:"weight"`
	Cert        string `json:"cert"`
	Address     string `json:"address"`
	DateCreate  string `json:"date_create"`
	DateCode    string `json:"date_code"`
	Composition string `json:"composition"`
}
