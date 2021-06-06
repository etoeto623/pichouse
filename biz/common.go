package biz

import (
	"math/rand"
	"time"
)

const magicStr = "maytheforcebewithyou"
var MagicBytes = []byte(magicStr)
const MagicLen = len(magicStr)
const AesKey = "thisisrandaeskey"

func BytesEqual(b1, b2 []byte) bool {
	if nil == b1 || nil == b2 ||
		len(b1) != len(b2){
		return false
	}
	for i:=0 ; i < len(b1); i++{
		if b1[i] != b2[i]{
			return false
		}
	}
	return true
}

func Log(msg string){
	//stamp := time.Now().Format("2006-01-02 15:04:05")
	//fmt.Println(stamp + ": " + msg)
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

type ModeHandler interface {
	handle ([]string)
}