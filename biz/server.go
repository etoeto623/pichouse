package biz

import (
	"bufio"
	"fmt"
	gmux "github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/user"
	"strings"
)

var houseBase = homeDir() + "/pichouse/"

func ServeMode(args []string){
	// 读取配置文件
	path, b := getCmdCfg(args, "c")
	if b {
		cfgPath = path
	}

	if len(strings.TrimSpace(GetCfg().PicHouse)) > 0 {
		houseBase = GetCfg().PicHouse
	}

	finish := make(chan bool)
	go func(){
		Log("server tcp port is: " + GetCfg().TcpPort)
		startTCPServer(GetCfg().TcpPort)
		finish<-true
	}() // 开启tcp服务器

	go func(){
		startHttpServer(GetCfg().HttpPort)
		finish<-true
	}()
	<-finish
}

func startHttpServer(port string){
	muxServer := gmux.NewRouter()
	muxServer.HandleFunc("/image/{token}", func(w http.ResponseWriter, r *http.Request) {
		vars := gmux.Vars(r)
		token := vars["token"]
		fileName, err := AesDecrypt(token, AesKey)
		if nil != err{
			fmt.Fprintln(w, "wrong token")
			return
		}
		// 读取文件并返回
		data, err := ioutil.ReadFile(houseBase + "/" + fileName)
		if nil != err{
			return
		}
		w.Write(data)
		Log("write image data")
	})
	http.ListenAndServe(":"+port, muxServer)
}

func startTCPServer(port string){
	Log("start tcp server at port: " + port)
	tcpAddr,_ := net.ResolveTCPAddr("tcp", ":"+port)
	tcpListener,_ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if nil != err {
			continue
		}
		Log("receive connect")
		go handleTCP(tcpConn)
	}
}

func handleTCP(conn *net.TCPConn){
	defer conn.Close()
	// 验证magic byte
	reader := bufio.NewReader(conn)
	if !verifyMagic(reader){
		Log("magic bytes error")
		closeConnAll(conn)
		return
	}

	fileName := GetRandomString(10) + ".png"
	filePath := houseBase + "/" + fileName
	file, err := os.Create(filePath)
	if nil != err{
		Log("file create error: " + err.Error())
		closeConnAll(conn)
		return
	}
	defer file.Close()

	var bytes = make([]byte, 1024)
	for{
		num, err := reader.Read(bytes)
		if nil != err && err != io.EOF{
			return
		}
		if num <= 0 || (nil != err && err == io.EOF) {
			break
		}
		file.Write(bytes[:num])
	}

	writer := bufio.NewWriter(conn)
	encrypted, err := AesEncrypt(fileName, AesKey)
	if nil != err {
		return
	}
	writer.WriteString(GetCfg().ViewImageUrl + "/image/" + encrypted)
	writer.Flush()
	conn.CloseWrite()
	Log("msg writed")
}

func closeConnAll(conn *net.TCPConn){
	conn.CloseRead()
	conn.CloseWrite()
}

func verifyMagic(reader *bufio.Reader) bool {
	mn := make([]byte, MagicLen)
	reader.Read(mn)
	return BytesEqual(mn, MagicBytes)
}

func homeDir() string {
	u, _ := user.Current()
	return u.HomeDir
}