package demo

import (
	"fmt"

	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

func usage() {
	fmt.Println("usage:")
	flag.PrintDefaults()
	os.Exit(-1)
}

func main1() {
	//	portName := flag.String("port", "", "serial port to test (/dev/ttyUSB0, etc)")

	// if *portName == "" {
	// 	fmt.Println("Must specify port")
	// 	usage()
	// }

	// Set up options.
	options := serial.OpenOptions{
		PortName:        "/dev/cu.usbmodem1441",
		BaudRate:        19200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	{
		buf := make([]byte, 5)
		n, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			fmt.Println("Rx: ", buf)
		}
	}

	// Make sure to close it later.
	defer port.Close()

	for i := 0; i < 10; i++ {

		r := i * 10
		g := i * 20
		b := (i * 40) % 255

		n, err := port.Write([]byte(
			fmt.Sprintf("LED%02x%02x%02x%02x\n", i, r, g, b),
		))
		if err != nil {
			log.Fatalf("port.Write: %v", err)
		}

		fmt.Println("Wrote", n, "bytes.")
	}

	time.Sleep(time.Second * 5)

	{
		n, err := port.Write([]byte("CLEAR\n"))
		if err != nil {
			log.Fatalf("port.Write: %v", err)
		}

		fmt.Println("Wrote", n, "bytes.")
	}
}
