package BLC

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"log"
	"math/big"
)

/*
	转化为字节数组
*/
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,num)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}

/*
	随机生成字符串
*/
func CreateRandomString(len int) string  {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0;i < len ;i++  {
		randomInt,_ := rand.Int(rand.Reader,bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}


/*
	Json转数组[]string
*/
func JSONToArray (jsonString string) [] string{
	var sArr [] string
	if err := json.Unmarshal([]byte(jsonString),&sArr);err != nil{
		log.Panic(err)
	}
	return sArr
}