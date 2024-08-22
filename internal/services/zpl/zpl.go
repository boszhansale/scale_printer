package zpl

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"simonwaldherr.de/go/zplgfa"
	"strings"
	"test/internal/services/zpl/element"
)

func Start(ll, pw int) string {
	return fmt.Sprintf("^XA^CI28^LL%d^PW%d", ll, pw)
}
func End() string {
	return "^XZ"
}
func Text(text *element.Text) string {

	return fmt.Sprintf("^FO%d,%d^FB%d,%d,0^AEN,%d,%d^FD%s^FS", text.X, text.Y, text.Width, text.Lines, text.FontSize, text.FontSize, text.Text)
}
func BoxText(text *element.Text, box *element.BorderBox) string {
	return fmt.Sprintf("^FO%d,%d^GB%d,%d,%d^FS^FO%d,%d^AEN,%d,%d^FD%s^FS",
		box.X, box.Y,
		box.Width, box.Height, box.Thickness,
		text.X,
		text.Y,
		text.FontSize,
		text.FontSize,
		text.Text,
	)
}
func Barcode(barcode *element.Barcode) string {
	return fmt.Sprintf("^FO%d,%d^BEN,%d,Y,N,N^FD%s^FS", barcode.X, barcode.Y, barcode.Height, barcode.Code)
}
func CopyCount(count int) string {
	return fmt.Sprintf("^PQ%d,0,1,Y", count)
}
func StaticImage(o *element.Orientation) string {
	return fmt.Sprintf("^FO%d,%d%s", o.X, o.Y, getStaticImage())

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
func OldSetNumberFont() string {
	data := `
		^XA
		^CWE,E:9835202.TTF
		^XZ
	`
	return data
}

func SetNumberFont() string {
	data := `
	^XA
	^CWE,E:ARIALR.TTF
	^XZ
	`
	return data
}
func ListFonts() string {
	return `
		^XA
		^LL525
		^PW463
		^WDE:*.TTF
		^XZ
	`
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
