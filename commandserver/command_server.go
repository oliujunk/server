package commandserver

import (
	"encoding/hex"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"net"
	"oliujunk/server/config"
	"oliujunk/server/database"
	"oliujunk/server/utils"
	"time"
)

var (
	linkMap map[net.Conn]LinkData
	job     *cron.Cron
)

const IdleTime = 120

// LinkData 已连接socket
type LinkData struct {
	idleTimer *time.Timer  // 超时定时器
	closed    chan bool    // 是否已关闭
	effective bool         // 是否有效
	deviceID  int          // 设备ID
	jobID     cron.EntryID // 任务ID
}

func Start() {
	log.Println("命令服务启动")

	job = cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))

	job.Start()

	linkMap = make(map[net.Conn]LinkData)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GlobalConfiguration.CommandServer.Port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening and serving HTTP on : %d", config.GlobalConfiguration.CommandServer.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go connectHandler(conn)
	}
}

func connectHandler(conn net.Conn) {
	log.Println("连接已建立", conn.RemoteAddr().String())
	linkData := LinkData{
		idleTimer: time.NewTimer(time.Second * IdleTime),
		closed:    make(chan bool),
		effective: false,
		deviceID:  0,
	}
	linkMap[conn] = linkData
	recvBuf := make([]byte, 1024)
	readID(conn)

	go func(conn net.Conn) {
		for {
			select {
			case <-linkMap[conn].closed:
				return
			case <-linkMap[conn].idleTimer.C:
				idleHandler(conn)
				return
			}
		}
	}(conn)

	for {
		recvLen, err := conn.Read(recvBuf)
		if recvLen == 0 || err != nil {
			closeHandler(conn)
			return
		}

		log.Printf(conn.RemoteAddr().String(), hex.EncodeToString(recvBuf[0:recvLen]))

		if utils.TestFrame(recvBuf[0:recvLen], recvLen) {
			processCommand(recvBuf[0:recvLen], conn)
		}

		linkData.idleTimer.Reset(time.Second * IdleTime)
	}
}

func idleHandler(conn net.Conn) {
	log.Println("连接已超时", conn.RemoteAddr().String())
	_ = conn.Close()
}

func closeHandler(conn net.Conn) {
	log.Println("连接已关闭", conn.RemoteAddr().String())
	_ = conn.Close()

	linkMap[conn].closed <- true
	close(linkMap[conn].closed)
	linkMap[conn].idleTimer.Stop()
	job.Remove(linkMap[conn].jobID)

	delete(linkMap, conn)
}

func readID(conn net.Conn) {
	sendBuf := []byte{0x00, 0x03, 0x00, 0x60, 0x00, 0x04}
	crc := utils.Crc16(sendBuf, len(sendBuf))
	sendBuf = append(sendBuf, (byte)(crc))
	sendBuf = append(sendBuf, (byte)(crc>>8))
	_, _ = conn.Write(sendBuf)
}

func readData(conn net.Conn) {
	sendBuf := []byte{0x00, 0x03, 0x00, 0x00}
	crc := utils.Crc16(sendBuf, len(sendBuf))
	sendBuf = append(sendBuf, (byte)(crc))
	sendBuf = append(sendBuf, (byte)(crc>>8))
	_, _ = conn.Write(sendBuf)
}

type ReadJob struct {
	conn net.Conn
}

func (readJob *ReadJob) Run() {
	readData(readJob.conn)
}

func processCommand(data []byte, conn net.Conn) {
	switch data[1] {
	case 0x03:
		{
			switch data[3] {
			case 0x40: // 数据
				saveData(linkMap[conn].deviceID, data)
				break
			case 0x60: // ID
				deviceID := int(data[6])*1000000 + int(data[7])*10000 + int(data[8])*100 + int(data[9])
				linkData := linkMap[conn]
				linkData.effective = true
				linkData.deviceID = deviceID
				//jobID, _ := job.AddJob("0 */1 * * * *", &ReadJob{conn: conn})
				jobID, _ := job.AddJob("*/10 * * * * *", &ReadJob{conn: conn})
				linkData.jobID = jobID
				linkMap[conn] = linkData
				break
			default:
				break
			}
			break
		}
	case 0x10:
		break
	default:
		break
	}
}

func saveData(deviceID int, data []byte) {
	element := make([]int16, 16)
	relay := make([]byte, 32)
	current := database.Current{
		DeviceID: deviceID,
		DataTime: time.Now(),
	}
	for i := 0; i < 16; i++ {
		element[i] = ((int16)(data[4+i*2]) << 8) + (int16)(data[5+i*2])
		current.SetElement(fmt.Sprintf("E%d", i+1), int64(element[i]))
	}
	for i := 0; i < 32; i++ {
		relay[i] = data[36+i]
		current.SetElement(fmt.Sprintf("J%d", i+1), int64(relay[i]))
	}
	_, err := database.Orm.Table("current").Insert(current)
	if err != nil {
		log.Println(err)
	}
}
