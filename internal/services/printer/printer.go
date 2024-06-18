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
		^GFA,1353,656,16,
000000000007FFFFFFDDDC0000F80000
000000000007FFFFFFDDDC0003FE0000
000000000007000001DDDC0007FF0000
3FE0FFF83FE7000001DDDC000F9F0000
3FE0FFF83FE7000001DDDC000F078000
3FE0FFF83FE7000001DDDC001E03C000
3FE0F8F83FE7000001DDDC001E03C000
3E00F0783FC7000001DDDC003C01E800
3E00F0783C07000001DDDC003801F800
3E00F0783C07000001DDDC007800F800
3E00F0783C07000001CCDC00F001F000
3E00F0783C07000001C01C003000F000
3E00F0783C07000001C01C0000003000
3E00F0783C07000001C01C0003FE0000
3E00F0783C07000001C01C0383FE0200
3E00F0783C07000001C03C0F839C0A00
3E00F0783C07000001C1F00F801C0F00
3E00F0783C07000001C1E00F801C0700
3E00FAF83C07000001C1C00F80380700
3FE0FFF83C0700000181C01E00380300
3FE0FFF83C0380000381C01C00300300
3FE0FFF83C03C0000701C03C00700100
3E00F0783C01E0000F01C07800600000
3E00F0783C00F0001E01C07800E00000
3E00F0783C007E01F801C07000E00000
3E00F0783C001FCFE001C07000000000
3E00F0783C0003CF0001C07000000000
3E00F0783C0000CE0001C07800000000
3E00F0783C0000CE0001C07C000E0100
3E00F0783C0000CE0001C03FFF1FFF00
3E00F0783C0000CE0001C01FFE3FFF00
3E00F0783C0000CE0001C00FFE3FFF00
3E00F0783C0001CE0001C000000E0000
3FE0F0783FE001CE0001C00000040000
3FE0F0783FE0018E0001C00000000000
3FE0F0783FE003870001C00000000000
3FE0F0783FE007078001C001EEDFBC00
1FE0F0783FE00E03E001C00124DE2400
0000000000007C01FC01C00124FF3C00
000000000003F0007F81C00124DE2800
00000000000380000F81C000E4DBAC00

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
		return data

	}
	if label.Paper == "70" {
		data += "^XA^CI28^LL900^PW500"
		data += "^FO375,690" + getStaticImage()
		//2800000014556
		data += fmt.Sprintf("^FO5,5^FB520,3,0^AEN,20,20^FD%s^FS", label.Name)
		data += fmt.Sprintf("^FO5,65^FB520,35,0^AEN,16,16^FD%s^FS", label.Description)

		data += fmt.Sprintf("^FO455,599^GB55,30,1^FS ^FO464,610^AEN,16,16^FD%s^FS", label.DateCode)

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

		return data
	}

	if label.Paper == "30" {
		data += "^XA^CI28^LL360^PW463"
		data += "^FO300,200" + getStaticImage()

		data += fmt.Sprintf("^FO390,95^GB55,30,1^FS ^FO399,100^AEN,16,16^FD%s^FS", label.DateCode)

		data += fmt.Sprintf("^FO20,130^AENб16,16^FDДайындалған күні:^FS")
		data += fmt.Sprintf("^FO20,150^AENб16,16^FDДата изготовления:^FS")
		data += fmt.Sprintf("^FO200,140^AENб16,16^FD%s^FS", date.Format("02/01/2006"))

		data += "^PQ" + countPrint + ",0,1,Y"
		data += "^XZ"

		return data
	}
	return ""
}

func GetImage() string {

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
