package database

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/go-xorm/xorm"
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// Orm orm引擎
	Orm *xorm.Engine
)

// Current 数据
type Current struct {
	DeviceID int       `xorm:"device_id" json:"deviceID"`
	DataTime time.Time `json:"dataTime"`
	E1       int64     `json:"e1"`
	E2       int64     `json:"e2"`
	E3       int64     `json:"e3"`
	E4       int64     `json:"e4"`
	E5       int64     `json:"e5"`
	E6       int64     `json:"e6"`
	E7       int64     `json:"e7"`
	E8       int64     `json:"e8"`
	E9       int64     `json:"e9"`
	E10      int64     `json:"e10"`
	E11      int64     `json:"e11"`
	E12      int64     `json:"e12"`
	E13      int64     `json:"e13"`
	E14      int64     `json:"e14"`
	E15      int64     `json:"e15"`
	E16      int64     `json:"e16"`
	J1       int64     `json:"j1"`
	J2       int64     `json:"j2"`
	J3       int64     `json:"j3"`
	J4       int64     `json:"j4"`
	J5       int64     `json:"j5"`
	J6       int64     `json:"j6"`
	J7       int64     `json:"j7"`
	J8       int64     `json:"j8"`
	J9       int64     `json:"j9"`
	J10      int64     `json:"j10"`
	J11      int64     `json:"j11"`
	J12      int64     `json:"j12"`
	J13      int64     `json:"j13"`
	J14      int64     `json:"j14"`
	J15      int64     `json:"j15"`
	J16      int64     `json:"j16"`
	J17      int64     `json:"j17"`
	J18      int64     `json:"j18"`
	J19      int64     `json:"j19"`
	J20      int64     `json:"j20"`
	J21      int64     `json:"j21"`
	J22      int64     `json:"j22"`
	J23      int64     `json:"j23"`
	J24      int64     `json:"j24"`
	J25      int64     `json:"j25"`
	J26      int64     `json:"j26"`
	J27      int64     `json:"j27"`
	J28      int64     `json:"j28"`
	J29      int64     `json:"j29"`
	J30      int64     `json:"j30"`
	J31      int64     `json:"j31"`
	J32      int64     `json:"j32"`
}

func (current Current) MarshalJSON() ([]byte, error) {
	type Alias Current
	return json.Marshal(&struct {
		DataTime string `json:"dataTime"`
		Alias
	}{
		DataTime: current.DataTime.Format("2006-01-02 15:04:05"),
		Alias:    (Alias)(current),
	})
}

func (current *Current) SetElement(key string, value int64) {
	v := reflect.ValueOf(current)
	v = v.Elem()
	t := v.FieldByName(key)
	t.SetInt(value)
}

type User struct {
	ID       int64  `xorm:"id pk" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Device struct {
	ID         int64     `xorm:"id pk" json:"id"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	DeviceID   int       `xorm:"device_id" json:"deviceID"`
	DeviceName string    `json:"deviceName"`
	EleNum     string    `json:"eleNum"`
	EleName    string    `json:"eleName"`
	RelayNum   string    `json:"relayNum"`
	RelayName  string    `json:"relayName"`
	CreatorID  int64     `xorm:"creator_id" json:"creatorID"`
	Longitude  float64   `json:"longitude"`
	Latitude   float64   `json:"latitude"`
}

func init() {
	// 数据库
	var err error
	//Orm, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/standalone",
	//	config.GlobalConfiguration.Database.Username,
	//	config.GlobalConfiguration.Database.Password,
	//	config.GlobalConfiguration.Database.Host,
	//	config.GlobalConfiguration.Database.Port,
	//))
	Orm, err = xorm.NewEngine("sqlite3", "main.db")

	if err != nil {
		log.Fatal(err)
	}

	//Orm.ShowSQL(true)
}
