package utils

import (
	"fmt"
	"strconv"
)

func BarcodeGenerate(barcode string, weight string) (string, error) {
	if barcode == "" {
		return barcode, nil
	}
	if string(barcode[0]) != "2" {
		return barcode, nil
	}
	head := barcode[:7]
	w, err := strconv.Atoi(weight)
	if err != nil {
		return "", err
	}
	newBarcode := head + fmt.Sprintf("%05d", w)
	fmt.Println("noasd " + newBarcode)
	odd := 0
	even := 0
	for i := 0; i < len(newBarcode); i++ {
		if i%2 != 0 {
			num, _ := strconv.Atoi(string(barcode[i]))
			even += num
		} else {
			num, _ := strconv.Atoi(string(barcode[i]))
			odd += num
		}
	}
	total := odd + (even * 3)

	fmt.Println("total: " + strconv.Itoa(total))
	if total%10 == 0 {
		return newBarcode + "0", nil
	} else {
		up := (((total / 10) + 1) * 10) - total
		fmt.Println("up: " + strconv.Itoa(up))
		return newBarcode + strconv.Itoa(up), nil
	}

}
