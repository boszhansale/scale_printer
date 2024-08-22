package main

import (
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
	"test/internal/logger"
	"test/internal/repository"
	"test/internal/services/label"
	"test/internal/services/printer"
	"test/internal/services/scale"
	"test/internal/utils"
	"time"
)

func main() {
	err := logger.Init("app.log")
	if err != nil {
		log.Fatal("Error initializing logger:", err)
	}
	jsonStr := repository.Get()
	cfg := config.NewConfig()
	db := repository.New(jsonStr)

	prt, err := printer.NewPrinter(cfg.PrinterName)
	if err != nil {
		logger.Error("error connect to printer: " + err.Error())
	}
	prt.Start(cfg.PrinterName)

	defer prt.Close()
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
		logger.Error("error connect to scale")
		utils.ErrorMessage(scaleErr, w)
	}

	go func() {
		var oldValue int64
		for {
			if scale == nil {
				continue
			}
			value, stb, err := scale.GetWeight()
			if err != nil {
				logger.Error("error get weight_app: ", err.Error())
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
			if err != nil {
				logger.Error("error get weightStr: ", err.Error())
				utils.ErrorMessage(err, w)
				return
			}
			product, err := db.GetProduct(selectedCategory, selectedProduct, selectedLang)
			if err != nil {
				logger.Error("get product error: ", err.Error())
				//errorMessage(err, w)
				continue
			}
			newBarcode, err := utils.BarcodeGenerate(product.Barcode, weightStr)

			if err != nil {
				logger.Error("error generate barcode: ", err.Error())
				utils.ErrorMessage(err, w)
				return
			}
			labelData := label.Label{}

			labelData.CreateDate = dateWidget.Text
			labelData.Weight = weightStr
			labelData.Barcode = newBarcode
			labelData.Paper = selectedPaper
			labelData.DateCode = utils.DateToCode()
			labelData.Lang = selectedLang
			labelData.DateBool = dateCheckWidget.Checked
			labelData.CountCopy = 1
			if selectedLang == "kz" {
				labelData.Name = product.NameKz
				labelData.Description = product.CompositionKz
				labelData.DescriptionRu = product.CompositionRu
				labelData.KzRuMargin = product.KzRuMargin
				labelData.Cert = product.CertKz
				labelData.Measure = product.Measure
				labelData.DateCode = product.DateType
			} else {
				labelData.Name = product.NameEn
				labelData.Description = product.CompositionEn
				labelData.DescriptionRu = product.CompositionRu
				labelData.KzRuMargin = product.KzRuMargin
				labelData.Cert = product.CertEn
				labelData.Measure = product.Measure
				labelData.DateType = product.DateType
			}
			command := labelData.Generate()
			err = prt.Print(command)
			if err != nil {
				logger.Error("label generate: ", err.Error())
				utils.ErrorMessage(err, w)
				return
			}
			logger.Info("print success")
		}
	}()

	w.SetContent(content)
	w.ShowAndRun()
}
