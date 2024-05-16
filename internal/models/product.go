package models

type Product struct {
	LabelCategoryId int    `json:"label_category_id"`
	NameKz          string `json:"name_kz"`
	NameEn          string `json:"name_en"`
	CompositionKz   string `json:"composition_kz"`
	CompositionEn   string `json:"composition_en"`
	Barcode         string `json:"barcode"`
	CertKz          string `json:"cert_kz"`
	CertEn          string `json:"cert_en"`
	Measure         string `json:"measure"`
}

//type Product struct {
//	Id          string `json:"id"`
//	Name        string `json:"name"`
//	Barcode     string `json:"barcode"`
//	Weight      string `json:"weight_app"`
//	Cert        string `json:"cert"`
//	Address     string `json:"address"`
//	DateCreate  string `json:"date_create"`
//	DateCode    string `json:"date_code"`
//	Composition string `json:"composition"`
//	Measure     string `json:"measure"`
//}
