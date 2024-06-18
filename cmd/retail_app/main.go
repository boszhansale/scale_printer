package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"test/internal/components"
	"test/internal/config"
	"test/internal/repository"
	"test/internal/services/printer"
	"test/internal/utils"
	"time"
)

func main() {

	printer.GetImage()
	jsonStr := repository.Get()

	cfg := config.NewConfig()
	db := repository.New(jsonStr)

	a := app.New()
	w := a.NewWindow("Штучный Печать этикеток")
	w.Resize(fyne.NewSize(900, 700))

	categories := db.GetCategoryNames()

	var products []string
	var papers = []string{"58", "70", "30"}

	var selectedCategory string
	var selectedProduct string
	var selectedPaper string
	var selectedLang = "kz"

	productsWidget := widget.NewSelect(products, func(selected string) {
		selectedProduct = selected
	})
	productsWidget.PlaceHolder = "выберите продукт"

	paperWidget := widget.NewSelect(papers, func(selected string) {
		selectedPaper = selected
	})
	paperWidget.PlaceHolder = "выберите размер бумаги"

	categoriesWidget := widget.NewSelect(categories, func(selected string) {
		selectedCategory = selected
		products = db.GetProductNames(selected, selectedLang)
		productsWidget.Options = products
		productsWidget.Selected = ""
		productsWidget.Refresh()
	})
	categoriesWidget.PlaceHolder = "выберите категорию"

	langWidget := widget.NewSelect([]string{"kz", "en"}, func(selected string) {
		selectedLang = selected
		products = db.GetProductNames(selectedCategory, selectedLang)

		productsWidget.Options = products
		productsWidget.Selected = ""
		productsWidget.Refresh()
	})
	langWidget.SetSelected(selectedLang)

	dateWidget := widget.NewEntry()
	dateWidget.Text = time.Now().Format("2006-01-02")

	dateCheckWidget := widget.NewCheck("Дата", func(b bool) {
		if b {
			dateWidget.Enable()
		} else {
			dateWidget.Disable()
		}
	})
	dateCheckWidget.SetChecked(true)

	weightBinding := binding.BindString(nil)
	weightBinding.Set("0")
	countPrintBinding := binding.BindString(nil)
	countPrintBinding.Set("1")

	weightWidget := widget.NewEntryWithData(weightBinding)
	weightWidget.OnChanged = func(text string) {
		if _, err := strconv.ParseFloat(text, 64); err != nil {
			weightWidget.SetText("")
		}
	}

	countPrintWidget := widget.NewEntryWithData(countPrintBinding)

	countPrintWidget.OnChanged = func(text string) {
		if _, err := strconv.Atoi(text); err != nil {
			// Если ввод не является целым числом, очистите виджет
			countPrintWidget.SetText("")
		}
	}
	margin := fyne.NewSize(30, 20)

	content := container.NewVBox(
		container.NewGridWithColumns(2,
			container.NewGridWithRows(2, widget.NewLabel("Категория"), categoriesWidget),
			container.NewGridWithRows(2, widget.NewLabel("Продукт"), productsWidget),
		),
		container.New(components.NewMarginLayout(margin)),
		container.NewGridWithColumns(2,
			container.NewGridWithRows(2, widget.NewLabel("Язык"), langWidget),
			container.NewGridWithRows(2, widget.NewLabel("Размер бумаги"), paperWidget),
		),
		container.New(components.NewMarginLayout(margin)),
		container.NewGridWithColumns(2, dateCheckWidget, dateWidget),

		container.NewGridWithColumns(2, widget.NewLabel("вес гр."), weightWidget),

		container.NewGridWithColumns(2, widget.NewLabel("количество печати"), countPrintWidget),
		container.New(components.NewMarginLayout(margin)),

		container.New(components.NewMarginLayout(margin)),

		&widget.Button{
			Text: "печать",
			OnTapped: func() {

				if selectedProduct == "" {
					utils.ErrorMessage(errors.New("выберите продукт"), w)
					return
				}
				if selectedPaper == "" {
					utils.ErrorMessage(errors.New("выберите размер бумаги"), w)
					return
				}
				weightStr, err := weightBinding.Get()
				if err != nil {
					utils.ErrorMessage(err, w)
					return
				}

				product, err := db.GetProduct(selectedCategory, selectedProduct, selectedLang)
				if err != nil {
					utils.ErrorMessage(err, w)
					return
				}

				newBarcode, err := utils.BarcodeGenerate(product.Barcode, weightStr)

				if err != nil {
					utils.ErrorMessage(err, w)
					return
				}
				label := printer.Label{}

				if selectedLang == "kz" {
					label.Name = product.NameKz
					label.Description = product.CompositionKz
					label.Cert = product.CertKz
					label.CreateDate = dateWidget.Text
					label.Weight = weightStr
					label.Barcode = newBarcode
					label.Paper = selectedPaper
					label.Measure = product.Measure
					label.DateCode = utils.DateToCode()
					label.Lang = selectedLang
				} else {
					label.Name = product.NameEn
					label.Description = product.CompositionEn
					label.Cert = product.CertEn
					label.CreateDate = dateWidget.Text
					label.Weight = weightStr
					label.Barcode = newBarcode
					label.Paper = selectedPaper
					label.Measure = product.Measure
					label.DateCode = utils.DateToCode()
					label.Lang = selectedLang
				}

				countPrintStr, err := countPrintBinding.Get()
				if err != nil {
					utils.ErrorMessage(err, w)
					return
				}
				countPrint, err := strconv.Atoi(countPrintStr)

				if err != nil {
					utils.ErrorMessage(err, w)
					return
				}
				batchSize := 50
				if countPrint <= batchSize {

					err = label.Print(cfg.PrinterName, countPrintStr)

					if err != nil {
						utils.ErrorMessage(err, w)
						return
					}
				} else {
					for i := 0; i < countPrint; i += batchSize {
						batchPrintCount := countPrint - i
						if batchPrintCount > batchSize {
							batchPrintCount = batchSize
						}

						err = label.Print(cfg.PrinterName, strconv.Itoa(batchPrintCount))
						if err != nil {
							utils.ErrorMessage(err, w)
							return
						}
					}
				}

				log.Println("send to printer")
				countPrintBinding.Set("1")
				weightBinding.Set("0")
			}},
		&layout.Spacer{},
		container.NewGridWithColumns(3, &widget.Button{Text: "скачать базу", OnTapped: func() {
			repository.Download()
			utils.MessageQuit("Перезапустите программу", w, a)
		}}),
	)

	fmt.Println(selectedCategory, selectedProduct)
	w.SetContent(content)
	w.ShowAndRun()

}
