package repository

import (
	"fmt"
	"io"
	"net/http"
)

//	type Repository struct {
//		cfg *config.Config
//	}
//
//	func NewRepo(cfg *config.Config) *Repository {
//		return &Repository{
//			cfg: cfg,
//		}
//	}
func send(method, url string, params io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		method, url, params,
	)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return nil, err
	}
	return body, nil

}

//
//func (r *Repository) GetCategories() ([]string, error) {
//
//	data, err := send("GET", r.cfg.CategoriesApi, nil)
//	if err != nil {
//		return nil, err
//	}
//	var categories []string
//
//	if err = json.Unmarshal(data, &categories); err != nil {
//		return nil, err
//	}
//	return categories, nil
//}
//
//func (r *Repository) GetProducts(category string, lang string) ([]string, error) {
//
//	url := r.cfg.ProductsApi + "?category=" + category + "&lang=" + lang
//	data, err := send("GET", url, nil)
//	if err != nil {
//		return nil, err
//	}
//	var products []string
//
//	if err = json.Unmarshal(data, &products); err != nil {
//		return nil, err
//	}
//	return products, nil
//}
//
//func (r *Repository) ProductCreate(LabelProductName, lang, weight string, dateShow bool, date string) (models.Product, error) {
//
//	// Создаем структуру для данных, которые вы хотите отправить
//	requestData := struct {
//		LabelProductName string `json:"label_product_name"`
//		Lang             string `json:"lang"`
//		Weight           string `json:"weight_app"`
//		DateShow         bool   `json:"date_show"`
//		Date             string `json:"date"`
//	}{
//		LabelProductName: LabelProductName,
//		Lang:             lang,
//		Weight:           weight,
//		DateShow:         dateShow,
//		Date:             date,
//	}
//
//	requestBody, err := json.Marshal(requestData)
//	if err != nil {
//		return models.Product{}, err
//	}
//
//	data, err := send("POST", r.cfg.ProductCreateApi, bytes.NewReader(requestBody))
//	if err != nil {
//		return models.Product{}, err
//	}
//
//	var product models.Product
//	if err = json.Unmarshal(data, &product); err != nil {
//		return models.Product{}, err
//	}
//	return product, nil
//}
