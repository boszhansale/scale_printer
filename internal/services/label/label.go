package label

import (
	"fmt"
	_ "image/png"
	"strconv"
	"test/internal/services/zpl"
	"test/internal/services/zpl/element"
	"time"
)

const (
	measurerKz = "Өндіруші: «Первомайские деликатесы» ЖШС, Қазақстан Республикасы, Алматы облысы, Іле ауданы, Қоянқұс ауылы, Абай көшесі, №200"
	measurerRu = "Изготовитель: ТОО«Первомайские Деликатесы», Республика Казахстан, Алматинская область, Илийский район, село Коянкус,улица Абай, №200. т:+7(727)260-36-48"
	measurerEn = "Manufacturer: Pervomayskie Delikatesy LLP, Republic of Kazakhstan, Almaty region, Ili district,Koyankus village,Abay Street, No. 200 tel: +7(727)260-36-48"
)

type Label struct {
	Name,
	Description,
	DescriptionRu,
	KzRuMargin,
	CreateDate,
	DateCode,
	Weight,
	Cert,
	Barcode,
	Paper,
	Lang,
	Measure,
	DateType string
	DateBool  bool
	CountCopy int
}

func (l *Label) Generate() string {

	date, _ := time.Parse("2006-01-02", l.CreateDate)

	switch l.Paper {
	case "70":
		data := l.paper70(date)
		return data
	case "30":
		data := l.paper30(date)
		return data
	default:
		data := l.paper58(date)
		return data
	}
}

// 58*92
func (l *Label) paper58(date time.Time) string {
	data := ""
	data += zpl.Start(525, 463)
	data += zpl.StaticImage(&element.Orientation{X: 300, Y: 590})

	data += zpl.Text(&element.Text{
		Text:     l.Name,
		X:        5,
		Y:        5,
		Width:    450,
		Lines:    4,
		FontSize: 20,
	})

	data += zpl.Text(&element.Text{
		Text:     l.Description,
		X:        5,
		Y:        78,
		Width:    445,
		Lines:    35,
		FontSize: 16,
	})

	data += zpl.BoxText(&element.Text{
		Text:     l.DateCode,
		X:        399,
		Y:        508,
		Width:    445,
		Lines:    1,
		FontSize: 16,
	}, &element.BorderBox{
		X:         390,
		Y:         495,
		Width:     55,
		Height:    30,
		Thickness: 1,
	})

	if l.DateType == "1" {
		data += zpl.Text(&element.Text{
			Text:     l.Cert,
			X:        10,
			Y:        525,
			Width:    445,
			Lines:    1,
			FontSize: 16,
		})
		if l.DateBool {
			data += zpl.Text(&element.Text{
				Text:     "Дайындалған күні/Дата изготовления " + date.Format("02/01/2006"),
				X:        10,
				Y:        540,
				Width:    445,
				Lines:    1,
				FontSize: 16,
			})
		}

	} else {
		data += zpl.Text(&element.Text{
			Text:     l.Cert,
			X:        10,
			Y:        510,
			Width:    445,
			Lines:    1,
			FontSize: 16,
		})

		if l.DateBool {
			data += zpl.Text(&element.Text{
				Text:     "Дайындалған және оралған күні" + date.Format("02/01/2006"),
				X:        10,
				Y:        525,
				Width:    445,
				Lines:    1,
				FontSize: 16,
			})
			data += zpl.Text(&element.Text{
				Text:     "Дата изготовления и упаковывания" + date.Format("02/01/2006"),
				X:        10,
				Y:        540,
				Width:    445,
				Lines:    1,
				FontSize: 16,
			})
		}
	}

	if l.Measure == "2" && l.Weight != "0" && l.Weight != "" {
		data += zpl.Text(&element.Text{
			Text:     fmt.Sprintf("таза салмағы/масса нетто: %s г", l.Weight),
			X:        10,
			Y:        555,
			Width:    445,
			Lines:    1,
			FontSize: 16,
		})
	}
	if l.Barcode != "" {
		data += zpl.Barcode(&element.Barcode{
			Code:   l.Barcode,
			X:      10,
			Y:      570,
			Height: 70,
		})
	}

	if l.Lang == "kz" {
		data += zpl.Text(&element.Text{
			Text:     measurerKz,
			X:        9,
			Y:        667,
			Width:    445,
			Lines:    6,
			FontSize: 15,
		})
		data += zpl.Text(&element.Text{
			Text:     measurerRu,
			X:        9,
			Y:        710,
			Width:    445,
			Lines:    6,
			FontSize: 15,
		})
	} else {
		data += zpl.Text(&element.Text{
			Text:     measurerEn,
			X:        9,
			Y:        667,
			Width:    445,
			Lines:    6,
			FontSize: 15,
		})
	}

	data += zpl.CopyCount(l.CountCopy)
	data += zpl.End()
	return data
}

// 68*108
func (l *Label) paper70(date time.Time) string {
	data := ""
	data += zpl.Start(900, 500)
	data += zpl.StaticImage(&element.Orientation{X: 365, Y: 690})
	data += zpl.Text(&element.Text{
		Text:     l.Name,
		X:        5,
		Y:        5,
		Width:    520,
		Lines:    3,
		FontSize: 20,
	})
	data += zpl.Text(&element.Text{
		Text:     l.Description,
		X:        5,
		Y:        65,
		Width:    515,
		Lines:    35,
		FontSize: 15,
	})
	y, _ := strconv.Atoi(l.KzRuMargin)
	data += zpl.Text(&element.Text{
		Text:     l.DescriptionRu,
		X:        5,
		Y:        y,
		Width:    515,
		Lines:    35,
		FontSize: 15,
	})

	data += zpl.BoxText(&element.Text{
		Text:     l.DateCode,
		X:        455,
		Y:        610,
		Width:    445,
		Lines:    1,
		FontSize: 16,
	}, &element.BorderBox{
		X:         445,
		Y:         599,
		Width:     55,
		Height:    30,
		Thickness: 1,
	})

	if l.DateType == "2" {
		data += zpl.Text(&element.Text{
			Text:     l.Cert,
			X:        10,
			Y:        625,
			Width:    445,
			Lines:    1,
			FontSize: 16,
		})
		if l.DateBool {
			data += zpl.Text(&element.Text{
				Text:     "Дайындалған күні/Дата изготовления " + date.Format("02/01/2006"),
				X:        10,
				Y:        640,
				Width:    445,
				Lines:    1,
				FontSize: 16,
			})
		}
	} else {
		zpl.Text(&element.Text{
			Text:     l.Cert,
			X:        10,
			Y:        610,
			Width:    445,
			Lines:    1,
			FontSize: 16,
		})
		if l.DateBool {
			data += zpl.Text(&element.Text{
				Text:     "Дайындалған және оралған күні " + date.Format("02/01/2006"),
				X:        10,
				Y:        625,
				Width:    445,
				Lines:    1,
				FontSize: 16,
			})
			data += zpl.Text(&element.Text{
				Text:     "Дата изготовления и упаковывания " + date.Format("02/01/2006"),
				X:        10,
				Y:        640,
				Width:    445,
				Lines:    1,
				FontSize: 16,
			})
		}
	}

	if l.Measure == "2" && l.Weight != "0" && l.Weight != "" {
		data += zpl.Text(&element.Text{
			Text:     fmt.Sprintf("таза салмағы/масса нетто: %s г", l.Weight),
			X:        10,
			Y:        655,
			Width:    445,
			Lines:    1,
			FontSize: 16,
		})
	}
	if l.Barcode != "" {
		data += zpl.Barcode(&element.Barcode{
			Code:   l.Barcode,
			X:      10,
			Y:      670,
			Height: 70,
		})
	}
	if l.Lang == "kz" {
		data += zpl.Text(&element.Text{
			Text:     measurerKz,
			X:        5,
			Y:        765,
			Width:    515,
			Lines:    6,
			FontSize: 15,
		})
		data += zpl.Text(&element.Text{
			Text:     measurerRu,
			X:        5,
			Y:        808,
			Width:    515,
			Lines:    6,
			FontSize: 15,
		})
	} else {
		data += zpl.Text(&element.Text{
			Text:     measurerEn,
			X:        5,
			Y:        765,
			Width:    515,
			Lines:    6,
			FontSize: 15,
		})
	}
	data += zpl.CopyCount(l.CountCopy)
	data += zpl.End()

	return data
}
func (l *Label) paper30(date time.Time) string {
	data := ""
	data += "^XA^CI28^LL360^PW463"
	data += zpl.Start(360, 463)
	data += zpl.StaticImage(&element.Orientation{X: 300, Y: 200})

	data += zpl.BoxText(&element.Text{
		Text:     l.DateCode,
		X:        399,
		Y:        100,
		Width:    300,
		Lines:    1,
		FontSize: 16,
	}, &element.BorderBox{
		X:         390,
		Y:         95,
		Width:     55,
		Height:    30,
		Thickness: 1,
	})

	data += zpl.Text(&element.Text{
		Text:     "Дайындалған күні:",
		X:        20,
		Y:        120,
		Width:    350,
		Lines:    1,
		FontSize: 16,
	})
	data += zpl.Text(&element.Text{
		Text:     "Дата изготовления:",
		X:        20,
		Y:        150,
		Width:    350,
		Lines:    1,
		FontSize: 16,
	})
	data += zpl.Text(&element.Text{
		Text:     date.Format("02/01/2006"),
		X:        200,
		Y:        140,
		Width:    350,
		Lines:    1,
		FontSize: 16,
	})

	data += zpl.CopyCount(l.CountCopy)
	data += zpl.End()

	return data

}
