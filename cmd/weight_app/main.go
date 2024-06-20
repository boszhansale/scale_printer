package main

import (
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
	"test/internal/services/scale"
	"test/internal/utils"
	"time"
)

func main() {
	jsonStr := repository.Get()
	db := repository.New(jsonStr)

	cfg := config.NewConfig()
	a := app.New()
	w := a.NewWindow("Весовой Печать этикеток")
	w.Resize(fyne.NewSize(900, 700))

	weightBinding := binding.BindString(nil)
	weightBinding.Set("0")

	var products []string
	var selectedCategory string
	var selectedProduct string
	var selectedPaper string
	var selectedLang = "kz"
	papers := []string{"58", "70"}
	categories := db.GetCategoryNames()
	stable := make(chan bool)

	//весы
	scale, scaleErr := scale.Connect(cfg.WeightAddress)
	if scaleErr != nil {
		log.Println("error connect to scale")
		utils.ErrorMessage(scaleErr, w)
	}

	go func() {
		var oldValue int64
		for {
			if scale == nil {
				continue
			}
			value, stb, err := scale.GetWeight()
			log.Println(value)
			if err != nil {
				log.Println("error get weight_app")
				log.Println(err)
			} else {
				weightBinding.Set(strconv.FormatInt(value, 10))

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

		widget.NewLabelWithData(weightBinding),

		&layout.Spacer{},
		container.NewGridWithColumns(3, &widget.Button{Text: "скачать базу", OnTapped: func() {
			repository.Download()
			utils.MessageQuit("Перезапустите программу", w, a)
		}}),
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
			weightStr, err := weightBinding.Get()
			product, err := db.GetProduct(selectedCategory, selectedProduct, selectedLang)
			if err != nil {
				log.Println(err)
				//errorMessage(err, w)
				continue
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
				label.DescriptionRu = product.CompositionRu
				label.KzRuMargin = product.KzRuMargin
				label.Cert = product.CertKz
				label.CreateDate = dateWidget.Text
				label.Weight = weightStr
				label.Barcode = newBarcode
				label.Paper = selectedPaper
				label.Measure = product.Measure
				label.DateCode = utils.DateToCode()
				label.Lang = selectedLang
				label.DateCode = product.DateType
				label.DateBool = dateCheckWidget.Checked
			} else {
				label.Name = product.NameEn
				label.Description = product.CompositionEn
				label.DescriptionRu = product.CompositionRu
				label.KzRuMargin = product.KzRuMargin
				label.Cert = product.CertEn
				label.CreateDate = dateWidget.Text
				label.Weight = weightStr
				label.Barcode = newBarcode
				label.Paper = selectedPaper
				label.Measure = product.Measure
				label.DateCode = utils.DateToCode()
				label.Lang = selectedLang
				label.DateType = product.DateType
				label.DateBool = dateCheckWidget.Checked
			}
			err = label.Print(cfg.PrinterName, "1")
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
