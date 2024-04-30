package main

import (
	"errors"
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
	"time"
)

func main() {

	cfg := config.NewConfig()
	repo := repository.NewRepo(cfg)
	a := app.New()
	w := a.NewWindow("Штучный Печать этикеток")
	w.Resize(fyne.NewSize(700, 800))

	categories, err := repo.GetCategories()
	if err != nil {
		errorMessageQuit(err, w, a)
	}

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
	_weightBinding := binding.BindString(nil)
	_weightBinding.Set("0")
	_weightWidget := widget.NewEntryWithData(_weightBinding)
	_weightWidget.OnChanged = func(text string) {
		if _, err := strconv.ParseFloat(text, 64); err != nil {
			// Если ввод не является числом, очистите виджет
			_weightWidget.SetText("")
		}
	}

	_countPrintBinding := binding.BindString(nil)
	_countPrintWidget := widget.NewEntryWithData(_countPrintBinding)

	_countPrintBinding.Set("1")

	_countPrintWidget.OnChanged = func(text string) {
		if _, err := strconv.Atoi(text); err != nil {
			// Если ввод не является целым числом, очистите виджет
			_countPrintWidget.SetText("")
		}
	}

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

		container.NewGridWithColumns(2, widget.NewLabel("вес гр."), _weightWidget),

		container.NewGridWithColumns(2, widget.NewLabel("количество печати"), _countPrintWidget),

		&widget.Button{
			Text: "печать",
			OnTapped: func() {

				if selectedProduct == "" {
					errorMessage(errors.New("выберите продукт"), w)
					return
				}
				if selectedPaper == "" {
					errorMessage(errors.New("выберите размер бумаги"), w)
					return
				}
				_weight, err := _weightBinding.Get()
				if err != nil {
					errorMessage(err, w)
					return
				}
				fmt.Println(selectedPaper)
				labelProduct, err := repo.ProductCreate(selectedProduct, lang, _weight, dateBool.Checked, dateWidget.Text)
				if err != nil {
					log.Println(err)
					errorMessage(err, w)
					return
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
				_countPrint, err := _countPrintBinding.Get()
				if err != nil {
					log.Println(err)
					errorMessage(err, w)
					return
				}
				if err != nil {
					log.Println(err)
					errorMessage(err, w)
					return
				}
				err = label.Print(cfg.PrinterName, _countPrint, _weight)

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
					errorMessage(err, w)
					return
				}

				log.Println("print success")
				_countPrintBinding.Set("1")
				_weightBinding.Set("0")
			}},
	)

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
