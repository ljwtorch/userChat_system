package main

import (
	"encoding/binary"
	"encoding/json"
	"example.com/m/v2/common/message"
	"fmt"
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

func writePkg(conn net.Conn, data []byte) (err error) {
	//先发送一个长度给对方
	var pkgLen uint32
	var buf [4]byte
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//此时发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) err ", err)
		return
	}
	//发送data本身
	n, err = conn.Write(buf[:4])
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) err ", err)
		return
	}
	return
}
