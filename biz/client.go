package biz

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func ClientMode(args []string){
	uploadUrl, b := getCmdCfg(args, "uu")

	if !b {
		Log("param uu(aka:uploadUrl) not specified")
		os.Exit(1)
	}

	args = filterArgs(args)
	if len(args) > 1{
		// 目前只支持上传一张图片
		return
	}
	fileUrl := strings.TrimSpace(args[0])
	// 如果文件直接就是一个网络地址，则直接返回
	if strings.HasPrefix(fileUrl, "http://") ||
			strings.HasPrefix(fileUrl, "https://") {
		backMsg(fileUrl)
		return
	}

	// 调用服务器接口上传文件
	content, err := readFile(fileUrl)
	if nil != err{
		backMsg(fileUrl)
		return
	}

	// tcp连接服务器
	tcpAddr, err := net.ResolveTCPAddr("tcp", uploadUrl)
	if nil != err{
		backMsg(fileUrl)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err{
		backMsg(fileUrl)
		return
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	writer.Write(MagicBytes)
	writer.Write(content)
	writer.Flush()
	conn.CloseWrite()  // 需要关闭写，否则服务端还会一直等待

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if nil != err && err != io.EOF{
		backMsg(fileUrl)
		return
	}
	backMsg(msg)
}

func readFile(path string) ([]byte, error){
	f, err := os.Open(path)
	if nil != err{
		Log("read local file error")
		// 读取文件异常，直接返回
		return nil, err
	}
	defer f.Close() // 延迟关闭文件
	return ioutil.ReadAll(f)
}

func backMsg(msg string){
	fmt.Println(msg)
}

func filterArgs(args []string) []string {
	var params []string
	for i := range args {
		if strings.HasPrefix(args[i], "-") {
			continue
		}
		params = append(params, args[i])
	}
	return params
}