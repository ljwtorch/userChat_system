# userChat_system
**客户端：**

1. 接收输入的id和pwd
2. 发送id和pwd
3. 接收到服务端返回的结果
4. 判断是成功还是失败，并显示对应的页面，怎样组织发送的数据

**服务器：**

1. 接收用户id、pwd（goroutine）
2. 比较
3. 返回结果

**发送的流程：**

1. 先创建一个Message的结构体
2. mes.Type=登录消息类型
3. mes.Data=登录消息的内容（序列化后）
4. 对mes进行序列化
5. 在网络传输中，最麻烦丢包

**接收数据的流程：**

1. 接收到客户端发送的长度
2. 根据接收到的长度len，在接收消息本身
3. 接收时要判断实际接收到的消息内容是否等于len
4. 如果不相等，就有纠错协议
5. 取到后再反序列化 message
6. 取出message.data(string) 反序列化 LoginMes
7. 取出loginMes.userid和loginMesPwd
8. 这时就可以比较
9. 根据比较结果，返回Message
10. 发送给客户端

**步骤：**

1. 完成客户端可以发送消息本身，服务器端可以正常接收到消息，并根据客户端发送的消息（LoginMes）判断用户的合法性，并返回相应的LoginResMes

   **思路分析：**

    - 让客户端发送消息本身
    - 服务器端接收到消息，然后反序列化成对应的消息结构体
    - 服务器端根据反序列化成对应的消息，判断是否登录用户是合法，返回LoginResMes
    - 客户端解析返回的LoginResMes，显示对应的界面
    - 这里我们需要做函数的封装
