package lightboard

import (
	"io"

	"fmt"

	"github.com/golang/glog"
	"github.com/jacobsa/go-serial/serial"
)

func NewLightBoard(port string, leds int) (*LightBoard, error) {
	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        19200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	p, err := serial.Open(options)
	if err != nil {
		glog.Fatalf("serial.Open: %v", err)
		return nil, err
	}

	// Wait for the board to be ready.
	{
		buf := make([]byte, 5)
		n, err := p.Read(buf)
		if err != nil {
			if err != io.EOF {
				glog.Error("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			glog.Info("LightBoard: ", string(buf))
		}
	}

	lightBoard := &LightBoard{
		Port:   p,
		Lights: make([]RGB, leds),
	}
	return lightBoard, nil
}

type LightBoard struct {
	Port   io.ReadWriteCloser
	Lights []RGB
}

type RGB struct {
	Red   int
	Green int
	Blue  int
}

func (lb *LightBoard) SetIntensity(light int, intensity float32) {
	led := light / 3
	color := light % 3
	value := intensity * 255

	switch color {
	case 0:
		lb.Lights[led].Red = int(value)
	case 1:
		lb.Lights[led].Green = int(value)
	case 2:
		lb.Lights[led].Blue = int(value)
	}

	n, err := lb.Port.Write([]byte(
		fmt.Sprintf("LED%02x%02x%02x%02x\n",
			led, lb.Lights[led].Red, lb.Lights[led].Green, lb.Lights[led].Blue),
	))
	if err != nil {
		glog.Errorf("LightBoard: port.Write: %v", err)
	}

	glog.Info("LightBoard: Wrote", n, "bytes.")
}

func (lb *LightBoard) Clear() {
	n, err := lb.Port.Write([]byte("CLEAR\n"))
	if err != nil {
		glog.Errorf("LightBoard: port.Write: %v", err)
	}

	glog.Info("LightBoard: Wrote", n, "bytes.")
}
