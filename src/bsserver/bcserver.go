package bcserver;

import (
	"fmt"
	"net"
	"utils"
	"io/ioutil"
)

func BCServerHandle(conn *net.UDPConn) {
	buf := make([]byte, 256)
	_, serverAddr, err := conn.ReadFromUDP(buf)
	utils.ChkErr(err)
	bytes, err := ioutil.ReadFile("middleware/data.json")
	utils.ChkErr(err)
	fmt.Println(bytes)
	conn.WriteToUDP(bytes, serverAddr)
}

func BCServer() {
	fmt.Println("*******************")
	udpaddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:8081")
    utils.ChkErr(err)
    //监听端口
    udpconn, err := net.ListenUDP("udp", udpaddr)
	utils.ChkErr(err)
	for {
		BCServerHandle(udpconn)
	}
}