package communication

import (
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"time"
)

var (
	port io.ReadWriteCloser
)

// Start 定时读取数据
func Start() {
	// 串口配置
	options := serial.OpenOptions{
		PortName: "COM9",
		//PortName:              "/dev/ttyUSB0",
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		ParityMode:            serial.PARITY_NONE,
		RTSCTSFlowControl:     false,
		InterCharacterTimeout: 100,
		MinimumReadSize:       0,
	}

	var err error
	port, err = serial.Open(options)
	if err != nil {
		log.Fatalln(err)
	}

	go read()

}

func read() {
	buf := make([]byte, 128, 128)
	for {
		c, err := port.Read(buf)
		if err != nil {
			log.Println(err)
		}
		timeout := time.After(time.Millisecond * 800)
		recvBuf := make([]byte, 0, 64)
		recvBuf = append(recvBuf, buf[0:c]...)
		for len(recvBuf) < 5 || recvBuf[len(recvBuf)-1] != '\n' {
			c, err := port.Read(buf)
			if err != nil {
				log.Println(err)
				break
			}
			recvBuf = append(recvBuf, buf[0:c]...)
			select {
			case <-timeout:
				break
			default:
				continue
			}
		}
		log.Println(string(recvBuf))
	}
}
