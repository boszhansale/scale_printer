package printer

import (
	"errors"
	"fmt"
	"github.com/alexbrainman/printer"
	_ "image/png"
	"log"
	"test/internal/services/zpl"
)

type Printer struct {
	*printer.Printer
}

func NewPrinter(name string) (*Printer, error) {
	p, err := printer.Open(name)
	if err != nil {
		return nil, errors.New("невозможно открыть принтер: " + name)
	}
	err = p.StartRawDocument("scale")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if name == "ZDesigner ZD888-203dpi ZPL" {
		_, err = p.Write([]byte(zpl.OldSetNumberFont()))
		if err != nil {
			log.Print(err)
			return nil, errors.New("ошибка при установке шрифта ")
		}
	} else {
		_, err = p.Write([]byte(zpl.SetNumberFont()))
		if err != nil {
			log.Print(err)
			return nil, errors.New("ошибка при установке шрифта ")
		}
	}
	return &Printer{p}, nil
}

func (p *Printer) Print(command string) error {

	_, err := p.Write([]byte(command))
	if err != nil {
		log.Print(err)
		return errors.New("ошибка при записи на принтер ")
	}
	fmt.Println("Записано на принтер: ", command)
	return nil
}
func P(command string) {
	name := "ZDesigner ZD888-203dpi ZPL"
	p, err := printer.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	err = p.StartRawDocument("scale")
	if err != nil {
		log.Fatal(err)
	}
	_, err = p.Write([]byte(zpl.OldSetNumberFont()))
	if err != nil {
		log.Fatal(err)
	}

	_, err = p.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
}
