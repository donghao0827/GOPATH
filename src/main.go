package main;

import (
	"bcserver"
	"middleware"
	// "log"
	// "net"
	// "time"
	// "utils"
	"sync"
	"fmt"
)

func main() {
	fmt.Println("service start")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		middleware.Middleware("127.0.0.1:53", "127.0.0.1:4700")
		wg.Done()
	}()
	wg.Add(1) 
	go func() {  
		bcserver.BCServer("127.0.0.1:4700")
		wg.Done()
	}()
	wg.Wait()
	// if addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:4700"); err != nil {
	// 	log.Fatalf("Unable to resolve address: %s", err)
	// } else {
	// 	go server(addr)
	// 	time.Sleep(200 * time.Millisecond)
	// 	go client(addr)

	// 	time.Sleep(5 * time.Second)
	// }
}

// func server(addr *net.UDPAddr) {
// 	conn, err := net.ListenUDP("udp", addr)
// 	utils.ChkErr(err)
// 	buf := make([]byte, 256)
// 	_, rAddr, errRecv := conn.ReadFromUDP(buf)
// 	fmt.Println(addr, rAddr)
// 	utils.ChkErr(errRecv)
// 	_, errSend := conn.WriteToUDP([]byte{3, 4}, rAddr)
// 	utils.ChkErr(errSend)
// }

// func client(addr *net.UDPAddr) {
// 	conn, errDial := net.DialUDP("udp", nil, addr); 
// 	utils.ChkErr(errDial)
// 	_, errSend := conn.Write([]byte{1})
// 	utils.ChkErr(errSend)
// 	fmt.Println("*********")
// 	buf := make([]byte, 256)
// 	_, errRecv := conn.Read(buf)
// 	utils.ChkErr(errRecv)
// }