package printer

import (
	"errors"
	"github.com/alexbrainman/printer"
	_ "image/png"
	"test/internal/logger"
	"test/internal/services/zpl"
)

type Printer struct {
	*printer.Printer
}

func NewPrinter(name string) (*Printer, error) {
	logger.Info("открываем принтер: " + name)
	p, err := printer.Open(name)
	if err != nil {
		logger.Error("ошибка при открытии принтера: " + err.Error())
		return nil, errors.New("невозможно открыть принтер: " + name)
	}

	return &Printer{p}, nil
}

func (p *Printer) Start(name string) error {
	err := p.StartDocument("scale", "RAW")
	if err != nil {
		logger.Error("ошибка при запуске документа: " + err.Error())
		return err
	}
	if name == "ZDesigner ZD888-203dpi ZPL" {
		_, err = p.Write([]byte(zpl.OldSetNumberFont()))
		if err != nil {
			logger.Error("ошибка при установке старого шрифта: " + err.Error())
			return errors.New("ошибка при установке шрифта ")
		}
	} else {
		_, err = p.Write([]byte(zpl.SetNumberFont()))
		if err != nil {
			logger.Error("ошибка при установке нового шрифта: " + err.Error())
			return errors.New("ошибка при установке шрифта ")
		}
	}
	return nil
}

func (p *Printer) Print(command string) error {

	_, err := p.Write([]byte(command))
	if err != nil {
		logger.Error("ошибка при записи на принтер: " + err.Error())
		return errors.New("ошибка при записи на принтер ")
	}

	return nil
}
func RawPrint(name, command string) error {
	p, err := printer.Open(name)
	if err != nil {
		logger.Error("ошибка при открытии принтера: " + err.Error())
		return err
	}
	defer func(p *printer.Printer) {
		err := p.Close()
		if err != nil {
			logger.Error("ошибка при закрытии принтера: " + err.Error())
		}
	}(p)

	err = p.StartDocument("scale", "RAW")
	if err != nil {
		logger.Error("ошибка при запуске документа: " + err.Error())
		return err
	}
	_, err = p.Write([]byte(zpl.SetNumberFont()))
	if err != nil {
		logger.Error("ошибка при установке нового шрифта: " + err.Error())
		return err
	}
	_, err = p.Write([]byte(command))
	if err != nil {
		logger.Error("ошибка при записи на принтер: " + err.Error())
		return err
	}
	err = p.EndDocument()
	if err != nil {
		logger.Error("ошибка при завершении документа: " + err.Error())
		return err
	}
	return nil
}
