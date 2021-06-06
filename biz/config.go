package biz

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var cfgPath = homeDir() + "/.pichouse"

// 配置信息
type Cfg struct {
	UploadUrl string // 服务端图片上传地址 127.0.0.1:8008
	PicHouse string // 图片保存地址
	ViewImageUrl string // 图片查看地址 http://127.0.0.1:8118/image/xxxxx
	HttpPort string // http服务的端口  8118
	TcpPort string // tcp服务的端口  8008
}

func (cfg Cfg) toString() string {
	data, e := json.Marshal(cfg)
	if nil != e {
		return ""
	}
	return string(data)
}

var config Cfg
var inited = false
var lock sync.Mutex

func GetCfg() Cfg {
	lock.Lock()
	defer lock.Unlock()
	if inited {
		return config
	}
	file, err := os.Open(cfgPath)
	if nil != err {
		Log("config file open error: " + err.Error())
		os.Exit(1)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if nil!=err {
		Log("config file read error")
		os.Exit(1)
	}
	e := json.Unmarshal(data, &config)
	if nil != e {
		Log("config file parse error")
		os.Exit(1)
	}
	inited = true
	return config
}

func getCmdCfg(args []string, typePort string) (string, bool){
	if nil == args || len(args) == 0 {
		return "",false
	}
	prefix := "-" + typePort + "="
	for i := range args {
		cfg := args[i]
		if strings.HasPrefix(cfg, prefix) {
			return strings.Replace(cfg, prefix, "", 1), true
		}
	}
	return "", false
}