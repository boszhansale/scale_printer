package printer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/alexbrainman/printer"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"
	"simonwaldherr.de/go/zplgfa"
	"strings"
	"time"
)

type Label struct {
	Name, Description, CreateDate, DateCode, Weight, Cert, Barcode, Paper, Lang, Measure string
}

func (l Label) Print(printerName string, countPrint string) error {

	//printerName := "ZDesigner ZD888-203dpi ZPL"
	//
	p, err := printer.Open(printerName)
	if err != nil {
		log.Println(printerName + ": error print not found")
		return errors.New("не найден принтер: " + printerName)
	}
	defer p.Close()
	err = p.StartRawDocument("scale")
	if err != nil {
		log.Print(err)
		return err
	}

	//_, err = p.Write([]byte(getFont()))
	//if err != nil {
	//	log.Print(err)
	//	return errors.New("ошибка при установке шрифта ")
	//}

	if printerName == "ZDesigner ZD888-203dpi ZPL" {
		_, err = p.Write([]byte(oldSetNumberFont()))
		if err != nil {
			log.Print(err)
			return errors.New("ошибка при установке шрифта ")
		}
	} else {
		_, err = p.Write([]byte(setNumberFont()))
		if err != nil {
			log.Print(err)
			return errors.New("ошибка при установке шрифта ")
		}
	}
	_, err = p.Write([]byte(getData(l, countPrint)))
	if err != nil {
		log.Print(err)
		return errors.New("ошибка при записи на принтер ")
	}

	return nil

}
func getFont() string {
	data := `
	^XA
	^LL525
	^PW463
	^WDE:*.TTF
	^XZ
`
	return data
}
func oldSetNumberFont() string {
	data := `
	^XA
	^CWE,E:9835202.TTF
	^XZ
`
	return data
}

func setNumberFont() string {
	data := `
	^XA
	^CWE,E:ARIALR.TTF
	^XZ
	`
	return data
}

func getStaticImage() string {

	return `
		^GFA,800,384,12,
		000000001FFFFFDB60007800
		000000001FFFFFDB6001FE00
		3F8FFE3F980000DB6003FE00
		3F8FFE3F980000DB60038700
		3F8FFE3F980000DB60070380
		3F8F1E3F980000DB60070380
		380F1E3C180000DB600E01E0
		380F1E3C180000DB600E01E0
		380F1E3C180000C0601C01E0
		380F1E3C180000C060000040
		380F1E3C180000C06000FC00
		380F1E3C180000C0E060FC08
		380F1E3C180000C1C1F08C38
		380F1E3C180001C780F0181C
		3D0FFE3C1800018701F0181C
		3F8FFE3C1C00018701C0180E
		3F8FFE3C0C00038703803007
		380F1E3C0600070703803007
		380F1E3C03800E0707006003
		380F1E3C01F9F80707006003
		380F1E3C0079E00707000003
		380F1E3C0019800707000003
		380F1E3C0019800703800C07
		380F1E3C0019800703FF9FFE
		380F1E3C0019800700FF9FFC
		380F1E3C0019800700000C00
		3F8F1E3F8019800700000000
		3F8F1E3F8031C00700000000
		3F8F1E3F8070E007003B4E70
		0000000000E07807002D7E50
		000000000FC03F0700297E70
		000000000E000F8300394F50
		`
}

func getData(label Label, countPrint string) string {
	data := ""
	date, _ := time.Parse("2006-01-02", label.CreateDate)

	if label.Paper == "58" {
		data += "^XA^CI28^LL725^PW463"
		data += "^FO300,590" + getStaticImage()

		data += fmt.Sprintf("^FO5,5^FB445,3,0^AEN,20,20^FD%s^FS", label.Name)
		data += fmt.Sprintf("^FO5,65^FB445,35,0^AEN,16,16^FD%s^FS", label.Description)

		data += fmt.Sprintf("^FO390,495^GB55,30,1^FS ^FO399,508^AEN,16,16^FD%s^FS", label.DateCode)

		data += fmt.Sprintf("^FO10,525^AENб16,16^FD%s^FS", label.Cert)
		data += fmt.Sprintf("^FO10,540^AENб16,16^FDДайындалған күні/Дата изготовления %s^FS", date.Format("02/01/2006"))
		if label.Measure == "2" && label.Weight != "0" && label.Weight != "" {
			data += "^FO10,555^AENб16,16^FDтаза салмағы/масса нетто: " + label.Weight + " гр +/-3%^FS"
		}
		if label.Barcode != "" {
			data += fmt.Sprintf("^FO10,570^BEN,70,Y,N,N^FD%s^FS", label.Barcode)
		}

		if label.Lang == "kz" {
			data += fmt.Sprintf("^FB445,6,0^FO9,667^AENб15,15^FD%s^FS", "Өндіруші: «Первомайские деликатесы» ЖШС, Қазақстан Республикасы, Алматы облысы, Іле ауданы, Қоянқұс ауылы, Абай көшесі, №200")
			data += fmt.Sprintf("^FB455,6,0^FO9,710^AENб15,15^FD%s^FS", "Изготовитель: ТОО«Первомайские Деликатесы», Республика Казахстан, Алматинская область, Илийский район, село Коянкус,улица Абай, №200. т:+7(727)260-36-48")
		} else {
			data += fmt.Sprintf("^FB445,6,0^FO9,667^AENб15,15^FD%s^FS", "Manufacturer: Pervomayskie Delikatesy LLP, Republic of Kazakhstan, Almaty region, Ili district,Koyankus village,Abay Street, No. 200 tel: +7(727)260-36-48")
		}

		data += "^PQ" + countPrint + ",0,1,Y"
		data += "^XZ"
	} else {
		data += "^XA^CI28^LL900^PW500"
		data += "^FO375,690" + getStaticImage()
		//2800000014556
		data += fmt.Sprintf("^FO5,5^FB520,3,0^AEN,20,20^FD%s^FS", label.Name)
		data += fmt.Sprintf("^FO5,65^FB520,35,0^AEN,16,16^FD%s^FS", label.Description)

		data += fmt.Sprintf("^FO465,599^GB55,30,1^FS ^FO474,610^AEN,16,16^FD%s^FS", label.DateCode)

		data += fmt.Sprintf("^FO10,625^AENб16,16^FD%s^FS", label.Cert)
		data += fmt.Sprintf("^FO10,640^AENб16,16^FDДайындалған күні/Дата изготовления %s^FS", date.Format("02/01/2006"))
		if label.Measure == "2" && label.Weight != "0" && label.Weight != "" {

			data += "^FO10,655^AENб16,16^FDтаза салмағы/масса нетто: " + label.Weight + " гр +/-3%^FS"

		}
		if label.Barcode != "" {
			data += fmt.Sprintf("^FO10,670^BEN,70,Y,N,N^FD%s^FS", label.Barcode)
		}
		if label.Lang == "kz" {
			data += fmt.Sprintf("^FB520,6,0^FO5,765^AENб15,15^FD%s^FS", "Өндіруші: «Первомайские деликатесы» ЖШС, Қазақстан Республикасы, Алматы облысы, Іле ауданы, Қоянқұс ауылы, Абай көшесі, №200")
			data += fmt.Sprintf("^FB520,6,0^FO5,808^AENб15,15^FD%s^FS", "Изготовитель: ТОО«Первомайские Деликатесы», Республика Казахстан, Алматинская область, Илийский район, село Коянкус,улица Абай, №200. т:+7(727)260-36-48")
		} else {
			data += fmt.Sprintf("^FB520,6,0^FO5,765^AENб15,15^FD%s^FS", "Manufacturer: Pervomayskie Delikatesy LLP, Republic of Kazakhstan, Almaty region, Ili district,Koyankus village,Abay Street, No. 200 tel: +7(727)260-36-48")
		}
		data += "^PQ" + countPrint + ",0,1,Y"
		data += "^XZ"
	}

	return data
}

func getImage() string {

	file, err := os.Open("image.png")
	if err != nil {
		log.Printf("Warning: could not open the file: %s\n", err)
	}

	defer file.Close()

	// load image head information
	config, format, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("Warning: image not compatible, format: %s, config: %v, error: %s\n", format, config, err)
	}

	// reset file pointer to the beginning of the file
	file.Seek(0, 0)

	// load and decode image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("Warning: could not decode the file, %s\n", err)
	}

	flat := zplgfa.FlattenImage(img)

	str := convertToGraphicField(flat)
	fmt.Println(str)

	return fmt.Sprintf("^XA,^FS ^FO20,20 %s^FS,^XZ", str)
}
func getGraphicData(img image.Image) []byte {
	var data []byte
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			data = append(data, byte(r>>8))
		}
	}
	return data
}
func convertToGraphicField(source image.Image) string {
	var gfType string
	size := source.Bounds().Size()
	width := size.X / 8
	height := size.Y
	if size.Y%8 != 0 {
		width = width + 1
	}

	var GraphicFieldData string

	for y := 0; y < size.Y; y++ {
		line := make([]uint8, width)
		lineIndex := 0
		index := uint8(0)
		currentByte := line[lineIndex]
		for x := 0; x < size.X; x++ {
			index = index + 1
			p := source.At(x, y)
			lum := color.Gray16Model.Convert(p).(color.Gray16)
			if lum.Y < math.MaxUint16/2 {
				currentByte = currentByte | (1 << (8 - index))
			}
			if index >= 8 {
				line[lineIndex] = currentByte
				lineIndex++
				if lineIndex < len(line) {
					currentByte = line[lineIndex]
				}
				index = 0
			}
		}

		hexstr := strings.ToUpper(hex.EncodeToString(line))

		GraphicFieldData += fmt.Sprintln(hexstr)

	}
	gfType = "A"

	return fmt.Sprintf("^GF%s,%d,%d,%d,\n%s", gfType, len(GraphicFieldData), width*height, width, GraphicFieldData)
}
