package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"github.com/shopspring/decimal"
)

func start() {
	// 初始化
	Init()
	// 计划任务
	c := cron.New()
	c.AddFunc(conf.Cron, run)
	c.Start()
}

func run() {
	// 设置日志存储位置
	logFile, err := os.OpenFile(conf.PathLog+"/"+time.Now().Format("2006-01-02")+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log_error.Println(err)
	}
	log.SetOutput(logFile)

	// 数据结果
	data := new(Data)

	// 数据输入
	if conf_in.ModbusTCP.Ip != "" {
		//主动连接服务器
		var addr strings.Builder
		addr.WriteString(conf_in.ModbusTCP.Ip)
		addr.WriteString(":")
		addr.WriteString(strconv.Itoa(conf_in.ModbusTCP.Port))
		conn, err := net.Dial("tcp", addr.String())
		if err != nil {
			log_error.Println(err)
			return
		}

		defer conn.Close()
		// 遍历表计
		for _, v := range conf_in.ModbusTCP.Meter {
			meter := Meter{Name: v.Name, Data: make(map[string]interface{})}
			// 消息体
			head := []byte{0, 0, 0, 0, 0, 6, byte(v.Address), 3}
			for _, v_1 := range v.Par {
				buf := new(bytes.Buffer)
				buf.Write(head)
				binary.Write(buf, binary.BigEndian, v_1.Start)
				binary.Write(buf, binary.BigEndian, v_1.Length)

				// 发送数据
				// fmt.Println(hex.EncodeToString(buf.Bytes()))
				conn.Write(buf.Bytes())

				// 接收数据
				receive := make([]byte, 9+v_1.Length*2)
				_, err := conn.Read(receive)
				if err != nil {
					log_error.Println(err)
					return
				}
				// fmt.Println(hex.EncodeToString(receive))
				// 数据赋值
				for _, v_2 := range v_1.Data {
					start := 9 + v_2.Offset
					val := receive[start : start+v_2.Length]
					switch v_2.Length {
					case 1:
					default:

					}
					// 处理浮点数，符合IEEE754格式
					if v_2.Float {
						meter.Data[v_2.Name] = math.Float32frombits(binary.BigEndian.Uint32(val))
					} else {
						// 是否有符号
						if v_2.Symbol {
							meter.Data[v_2.Name], _ = bytesToIntS(val)
						} else {
							meter.Data[v_2.Name], _ = bytesToIntU(val)
						}
					}
					// 倍率处理
					value, _ := decimal.NewFromString(strval(meter.Data[v_2.Name]))
					rate, _ := decimal.NewFromString(strval(v_2.Rate))
					meter.Data[v_2.Name] = value.Mul(rate)
				}
			}
			data.Data = append(data.Data, meter)
		}
		data.Time = time.Now().Unix()
	}

	// 数据输出
	if conf_out.MQTT.Broker != "" {
		// 连接MQTT
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", conf_out.MQTT.Broker, conf_out.MQTT.Port))
		opts.SetUsername(conf_out.MQTT.Username)
		opts.SetPassword(conf_out.MQTT.Password)
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log_error.Println(token.Error())
			return
		}
		payload, _ := json.Marshal(data)
		log.Println(string(payload))
		client.Publish(conf_out.MQTT.Theme, 0, false, payload)
		client.Disconnect(250)
	}
	// MySQL
	// if conf_out.MySQL.Ip != "" {
	// 	db, err := sql.Open("mysql", conf_out.MySQL.Username+":"+conf_out.MySQL.Password+"@("+conf_out.MySQL.Ip+":"+strconv.Itoa(conf_out.MySQL.Port)+")/"+conf_out.MySQL.Dbname)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(db)
	// }
}

// 定义数据结构体
type Data struct {
	Time int64   `json:"time"`
	Data []Meter `json:"data"`
}

// 数据项
type Meter struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}
