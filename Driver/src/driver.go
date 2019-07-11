package driver

import (
	"fmt"
	"log"

	// "github.com/argandas/serial"
	"go.bug.st/serial.v1"
	"go.bug.st/serial.v1/enumerator"
)

type Reader struct {
	serialNo string
	name     string
	buffer   string
	port     serial.Port
}

func (r Reader) Init(SerialNo string) {
	r.serialNo = SerialNo
	Name, err := findThePort(SerialNo)
	if err != nil {
		log.Fatal(err)
	}
	r.name = Name
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	r.port.SetMode(mode)
	serial.Open(r.name, mode)
}
func (r Reader) close() {
	defer r.port.Close()
}
func (r Reader) read() {
	buff := make([]byte, 27)
	for {
		n, err := r.port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			break
		}
		r.buffer = r.buffer + string(buff)
	}
}
func (r Reader) Buffer() (buffer string, err error) {
	r.read()
	return r.buffer, nil
}
func (r Reader) ClearBuffer() (err error) {
	r.buffer = ""
	r.port.ResetOutputBuffer()
	r.port.ResetInputBuffer()
	return nil
}
func findThePort(SerialNo string) (port string, err error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}
	for _, port := range ports {
		if port.IsUSB {
			if port.SerialNumber == SerialNo {
				return port.Name, nil
			}
		}
	}
	return
}
