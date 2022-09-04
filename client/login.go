package main

import (
	"encoding/binary"
	"encoding/json"
	"example.com/m/v2/common/message"
	"fmt"
	"net"
	"time"
)

//登录函数
func login(userId int, userPwd string) (err error) {
	//定协议
	//fmt.Printf(" userId = %d userPwd=%s\n", userId, userPwd)
	//return err

	//1,连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dail err =", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.LoginMesType

	//3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4.将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}
	//5.把data赋给mes.Data
	mes.Data = string(data) //data需要转成string才能赋

	//6.将mes进行序列化 此时的data是要发送的数据
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
		return
	}

	//7.先把data的长度发给服务器
	//先获取到data的长度，然后转换成一个表示长度的byte切片
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
	//fmt.Println("客户端，发送消息的长度成功,长度为 =", len(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) err ", err)
		return
	}
	//休眠10
	time.Sleep(time.Second * 10)
	fmt.Println("休眠了10s")

	//这里还需要处理服务器端返回的消息
	return

}
