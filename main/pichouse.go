package main

import (
	"pichouse/biz"
	"fmt"
	"os"
)

/**
有两种使用方式：客户端模式和服务器端模式
客户端模式：
	pichouse client -tcp=80 -http=81
服务器端模式：
	pichosue server file_url
 */

func main() {
	args := os.Args
	if len(args) <= 1 {
		// 没有指定图片信息
		return
	}

	dispatch := map[string]func([]string){
		"help":   helpMode,
		"server": biz.ServeMode,
		"client": biz.ClientMode,
	}

	fn := dispatch[args[1]]
	if nil == fn{
		fn = helpMode
	}
	fn(args[2:])
}

// 打印帮助信息
func helpMode(args []string) {
	fmt.Println(`Usage is as bellow:
  server [-port=xx]
  client file_url
  help`)
}