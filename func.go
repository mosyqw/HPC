package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

//字节数(大端)组转成int(无符号的)
func bytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 8:
		var tmp uint64
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

//字节数(大端)组转成int(有符号)
func bytesToIntS(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 8:
		var tmp int64
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

// interface 转 string
func strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch reflect.TypeOf(value).Name() {
	case "float32":
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 32)
	case "float64":
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case "int":
		it := value.(int)
		key = strconv.Itoa(it)
	case "uint":
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case "int8":
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case "uint8":
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case "int16":
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case "uint16":
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case "int32":
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case "uint32":
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case "int64":
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case "uint64":
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case "string":
		key = value.(string)
	case "[]byte":
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
