package main

import (
	"encoding/binary"
	"encoding/json"
	"example.com/m/v2/common/message"
	"fmt"
	"io"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	//这里我们将读取数据包，直接封装成遗憾函数readPkg()，返回Message，Err
	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据……")
	_, err = conn.Read(buf[0:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//根据buf[:4]转成一个unit32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])

	//根据pkgLen读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}

	//再将pkgLen反序列化成 --- message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes) //这里要有& 否则返回是空的
	if err != nil {
		fmt.Println("json.Unmarshal err =", err)
		return
	}
	return
}

//处理客户端的通讯
func process(conn net.Conn) {
	//这里需要延时关闭conn
	defer conn.Close()
	for {
		//循环读客户端发送的信息
		mes, err := readPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务端退出……")
				return
			} else {
				fmt.Println("readPkg(conn) err =", err)
				return
			}
		}
		fmt.Println("mes =", mes)
	}
}

func main() {
	//提示信息
	fmt.Println("服务器再8889端口监听……")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err =", err)
		return
	}
	//一旦监听成功就等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来链接服务器……")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept() err =", err)
		}
		//一旦链接成功，则启动一个协程和客户端保持通讯
		go process(conn)
	}

}
