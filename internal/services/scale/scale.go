package scale

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"
)

type Scale struct {
	conn net.Conn
}

func Connect(address string) (*Scale, error) {
	log.Println("Connecting to ", address)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)

	if err != nil {
		return nil, errors.New("Не удалось подключиться к весам ")
	}
	//defer conn.Close()

	return &Scale{conn}, nil

}

func (s *Scale) GetWeight() (int64, bool, error) {
	header := []byte{0xF8, 0x55, 0xCE, 0x01, 0x00}
	command := append(header, 0xA0, 0xA0, 0x00)
	_, err := s.conn.Write(command)
	if err != nil {
		return 0, false, err
	}

	reader := bufio.NewReader(s.conn)
	response := make([]byte, 14)
	_, err = reader.Read(response)
	if err != nil {
		return 0, false, err
	}
	weightData := response[6:8]
	sta := response[11] == 0x01

	value := int64(binary.LittleEndian.Uint16(weightData))
	return value, sta, nil
}
