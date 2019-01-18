package send;

import (
	"net"
	"utils"
	"fmt"
)

func Send() {
	  //获取udpaddr
	  udpaddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:6001");
	  utils.ChkErr(err);
	  //连接，返回udpconn
	  udpconn, err1 := net.DialUDP("udp", nil, udpaddr);
	  utils.ChkErr(err1);
	  //写入数据
	  _, err2 := udpconn.Write([]byte("client\r\n"));
	  utils.ChkErr(err2);
	  buf := make([]byte, 256);
	  //读取服务端发送的数据
	  _, err3 := udpconn.Read(buf);
	  utils.ChkErr(err3);
	  fmt.Println(string(buf));
}