package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"test/internal/config"
	"test/internal/repository"
	"test/internal/services/printer"
	"test/internal/services/scale"
	"time"
)

func main() {

	//testLabel := printer.Label{
	//	Name:         "«Қойдың қарыны» «Требуха баранья»",
	//	Id:           "1234",
	//	Description:  "Мұздатылған өңделген ет субөнімі. Құрамы:қойдың қарыны. 100 г өнімнің тағамдық құндылығы: ақуыз 18 г, май 4 г. Энергетикалық құндылығы/ құнарлылығы 108 ккал/443,8 кДж. Орташа көрсеткіштер келтірілген. Вакуумда қапталған. Жарамдылық мерзімі: минус 18°С сақтау температурасында қаптамада 180 тәуліктен артық емес, қаптамасын ашқан соң- жарамдылық мерзімі шегінде 12 сағат. Өнімді қайта мұздатуға болмайды. Дайындау тәсілі: өнімді термиялық өңдеуге (пісіруге) болады. Жасап шығарылған күнін термочектен қараңыз.",
	//	Manufacturer: "Өндіруші:«Первомайские деликатесы» ЖШС, Қазақстан Республикасы, Алматы облысы, Іле ауданы, Қоянқұс ауылы, Абай көшесі, №200\nИзготовитель: ТОО«Первомайские Деликатесы», Республика Казахстан, Алматинская область, Илийский район, село Коянкус,улица Абай, №200. т:+7 775 256 22 55",
	//	CreateDate:   "Дайындалған күні/Дата изготовления: 17/04/2024",
	//	DateCode:     "1604",
	//	Weight:       "таза салмағы: 250 гр +/-3%",
	//	Cert:         "ГОСТ 32244-2013",
	//	Barcode:      "1233321321",
	//	Paper:        "70",
	//}
	//
	//testLabel.Print(cfg.PrinterName)
	//
	//time.Sleep(time.Second * 65)

	cfg := config.NewConfig()
	repo := repository.NewRepo(cfg)
	a := app.New()
	w := a.NewWindow("Весовой Печать этикеток")
	w.Resize(fyne.NewSize(700, 800))
	bindingWeight := binding.NewString()

	//весы

	scale, scaleErr := scale.Connect(cfg.WeightAddress)
	if scaleErr != nil {
		log.Println("error connect to scale")
		errorMessage(scaleErr, w)
	}
	categories, err := repo.GetCategories()
	if err != nil {
		errorMessageQuit(err, w, a)
	}
	stable := make(chan bool)
	go func() {
		var oldValue int64
		for {
			if scale == nil {
				continue
			}
			value, stb, err := scale.GetWeight()
			log.Println(value)
			if err != nil {
				log.Println("error get weight")
				log.Println(err)
			} else {
				bindingWeight.Set(strconv.FormatInt(value, 10))

				if !stb {
					continue
				}
				if value <= 10 {
					oldValue = 0
					continue
				}
				if value == oldValue {
					continue
				}

				stable <- stb
				oldValue = value
				time.Sleep(time.Second * 1)
			}

		}
	}()

	var products []string

	var selectedCategory string
	var selectedProduct string
	var selectedPaper string
	var lang = "kz"
	var productPlaceHolder = "выберите продукт"
	papers := []string{"58", "70"}
	productsWidget := widget.NewSelect(products, func(selected string) {
		selectedProduct = selected
	})
	paperWidget := widget.NewSelect(papers, func(selected string) {
		selectedPaper = selected
	})
	productsWidget.PlaceHolder = productPlaceHolder
	categoriesWidget := widget.NewSelect(categories, func(selected string) {

		selectedCategory = selected

		products, err = repo.GetProducts(selected, lang)
		if err != nil {
			log.Fatal(err)
		}
		productsWidget.Options = products
		productsWidget.SetSelected("")
		productsWidget.Selected = ""
		productPlaceHolder = "выберите продукт"
		productsWidget.Refresh()

	})
	categoriesWidget.PlaceHolder = "выберите категорию"
	paperWidget.PlaceHolder = "выберите размер бумаги"

	langWidget := widget.NewSelect([]string{"kz", "en"}, func(selected string) {
		lang = selected
		products, err = repo.GetProducts(selectedCategory, lang)
		if err != nil {
			log.Fatal(err)
		}
		productsWidget.Options = products
		productsWidget.SetSelected("")
		productsWidget.Selected = ""
		productPlaceHolder = "выберите продукт"
		productsWidget.Refresh()
	})
	langWidget.SetSelected(lang)

	dateWidget := widget.NewEntry()
	dateWidget.Text = time.Now().Format("2006-01-02")

	dateBool := widget.NewCheck("Дата", func(b bool) {
		if b {
			dateWidget.Enable()
		} else {
			dateWidget.Disable()
		}
	})
	dateBool.SetChecked(true)

	content := container.NewVBox(
		widget.NewLabel("Категория"),
		categoriesWidget,
		widget.NewLabel("Язык"),
		langWidget,
		widget.NewLabel("Продукт"),
		productsWidget,
		widget.NewLabel("Размер бумаги"),
		paperWidget,

		container.NewGridWithColumns(2, dateBool, dateWidget),
		widget.NewLabelWithData(bindingWeight),
	)

	go func() {
		for {
			<-stable
			if selectedProduct == "" {
				continue
				//errorMessage(errors.New("выберите продукт"), w)
				//return
			}
			if selectedPaper == "" {
				continue
			}
			_weight, _ := bindingWeight.Get()
			fmt.Println(selectedPaper)
			labelProduct, err := repo.ProductCreate(selectedProduct, lang, _weight, dateBool.Checked, dateWidget.Text)
			if err != nil {
				log.Println(err)
				//errorMessage(err, w)
				continue
			}
			label := printer.Label{
				Name:         labelProduct.Name,
				Description:  labelProduct.Composition,
				Id:           labelProduct.Id,
				DateCode:     labelProduct.DateCode,
				Manufacturer: labelProduct.Address,
				Cert:         labelProduct.Cert,
				CreateDate:   labelProduct.DateCreate,
				Weight:       labelProduct.Weight,
				Barcode:      labelProduct.Barcode,
				Paper:        selectedPaper,
				Measure:      labelProduct.Measure,
			}
			err = label.Print(cfg.PrinterName)
			if err != nil {
				log.Println("name: " + labelProduct.Name)
				log.Println("description: " + labelProduct.Composition)
				log.Println("Id: " + labelProduct.Id)
				log.Println("DateCode: " + labelProduct.DateCode)
				log.Println("Manufacturer: " + labelProduct.Address)
				log.Println("Cert: " + labelProduct.Cert)
				log.Println("CreateDate: " + labelProduct.DateCreate)
				log.Println("Weight: " + labelProduct.Weight)
				log.Println("Barcode: " + labelProduct.Barcode)
				log.Println("Paper: " + selectedPaper)
				log.Println("Measure: " + labelProduct.Measure)

				log.Println(err)
				continue
				//errorMessage(err, w)
			}
			log.Println("print success")
		}
	}()

	fmt.Println(selectedCategory, selectedProduct)
	w.SetContent(content)
	w.ShowAndRun()
}

func errorMessageQuit(err error, window fyne.Window, app fyne.App) {
	infoDialog := dialog.NewInformation("Ошибка", err.Error(), window)
	infoDialog.Show()
	infoDialog.SetOnClosed(func() {
		app.Quit()
	})
}
func errorMessage(err error, window fyne.Window) {
	infoDialog := dialog.NewInformation("Ошибка", err.Error(), window)
	infoDialog.Show()
}
