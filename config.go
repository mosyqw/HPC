package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var conf = &Config{}
var conf_in = &ConfigIO{}
var conf_out = &ConfigIO{}
var log_error *log.Logger

type Config struct {
	Cron    string
	Path    string
	PathLog string
}

type ConfigIO struct {
	ModbusTCP ModbusTCP
	MQTT      MQTT
	MySQL     MySQL
}

type ModbusTCP struct {
	Ip    string
	Port  int
	Retry int
	Cycle int
	Meter []ModbusTCPMeter
}

type ModbusTCPMeter struct {
	Name     string
	Address  int
	Timeout  int
	Interval int
	Par      []ModbusTCPPar
}

type ModbusTCPPar struct {
	Start  uint16
	Length uint16
	Data   []ModbusTCPData
}

type ModbusTCPData struct {
	Name   string
	Offset int
	Length int
	Rate   interface{}
	Symbol bool
	Before bool
	Flip   bool
	Float  bool
}

type MQTT struct {
	Broker   string
	Port     int
	Username string
	Password string
	Theme    string
}

type MySQL struct {
	Ip       string
	Port     int
	Username string
	Password string
	Dbname   string
}

func Init() (err error) {
	// 获取程序安装路径
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	conf.Path = path
	// 创建文件夹
	conf.PathLog = path + "/log"
	exist, _ := PathExists(conf.PathLog)
	if !exist {
		err := os.Mkdir(conf.PathLog, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	// 错误日志输出
	log_error_file, err := os.OpenFile(conf.PathLog+"/error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log_error = log.New(log_error_file, "", log.LstdFlags|log.LUTC) // 将文件设置为loger作为输出
	if _, err := toml.DecodeFile(path+"/config.toml", &conf); err != nil {
		log_error.Println(err)
	}
	if _, err := toml.DecodeFile(path+"/in.toml", &conf_in); err != nil {
		log_error.Println(err)
	}
	if _, err := toml.DecodeFile(path+"/out.toml", &conf_out); err != nil {
		log_error.Println(err)
	}
	return
}
