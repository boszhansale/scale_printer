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

	cfg := config.NewConfig()
	repo := repository.NewRepo(cfg)

	a := app.New()
	w := a.NewWindow("Печать этикеток")
	w.Resize(fyne.NewSize(600, 700))
	bindingWeight := binding.NewString()
	//весы

	s, err := scale.Connect(cfg.WeightAddress)
	if err != nil {
		errorMessageQuit(err, w, a)
	}
	categories, err := repo.GetCategories()
	if err != nil {
		errorMessageQuit(err, w, a)
	}
	stable := make(chan bool)
	go func() {
		var oldValue int64
		for {

			value, stb, err := s.GetWeight()
			if err != nil {
				log.Println(err)
			}
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
	}()

	var products []string

	var selectedCategory string
	var selectedProduct string
	var lang = "kz"
	var productPlaceHolder = "выберите продукт"

	productsWidget := widget.NewSelect(products, func(selected string) {
		selectedProduct = selected
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
			_weight, _ := bindingWeight.Get()

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
			}
			err = label.Print(cfg.PrinterName)
			if err != nil {
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
