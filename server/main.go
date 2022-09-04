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

//severProcessLogin 专门处理登录请求
func severProcessLogin(conn net.Conn, mes *message.Message) (err error) {
	//核心代码
	//1.先从mes中取出mes.Data，并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err =", err)
		return
	}

	//声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//2.声明一个LoginResMes
	var loginResMes message.LoginResMes

	//如果用户id=100 密码=123456 认为合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		//合法
		loginResMes.Code = 200
	} else {
		//不合法 500表示不存在
		loginResMes.Code = 500
		loginResMes.Error = "该用户不存在，注册后再使用"
	}
	//3.将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) err =", err)
		return
	}
	//4.将data赋值给resMes
	resMes.Data = string(data)
	//5.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) err =", err)
		return
	}
	//6.发送data 封装到writePkg函数
	err = writePkg(conn, data)
	return
}

//编写一个ServerProcessMes函数；根据客户端发送消息种类不同，决定调用哪个函数来处理
func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		err = severProcessLogin(conn, mes)
	case message.RegisterMesType:
	//处理注册
	default:
		fmt.Println("消息类型不存在，无法处理……")
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
		//fmt.Println("mes =", mes)
		err = serverProcessMes(conn, &mes)
		if err != nil {
			return
		}
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
