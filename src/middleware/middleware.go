package middleware;

import (
    "net"
    //"net/http"
    //"fmt"
    "encoding/binary"
    //"strings"
    "bytes"
    //"encoding/json"
    //"io/ioutil"
    "utils"
)

//解析区块链返回数据
// func ParseBlockResponse(domainName string) (interface {}) {
// 	/*
//     //发送http post请求
//     client := &http.Client{}
//     //request, _ := http.NewRequest("POST", "http://192.168.27.99:5002/query?url=www.qq.com&timestamp=178542318", nil)
//     request, _ := http.NewRequest("POST", "http://localhost/server.php", nil)
//     request.Header.Set("Content-type", "application/json")
//     response, err := client.Do(request)
//     chkError(err)
//     var res Response
//     var json = jsoniter.ConfigCompatibleWithStandardLibrary
//     if response.StatusCode == 200 {
//         body, _ := ioutil.ReadAll(response.Body)
//         err := json.Unmarshal(body, &res)
//         if err != nil {
//             return false
//         }
//         return res
//     }
//     return false
//     */
//     /*
//     bytes, errRead := ioutil.ReadFile("./data.json")
//     chkError(errRead)
//     var res Response
//     var json = jsoniter.ConfigCompatibleWithStandardLibrary
//     err := json.Unmarshal(bytes, &res)
//     chkError(err)
//     return res
//     */
//     return {}
// }

func UdpHandle(data []byte) []byte {
    //defer conn.Close()
    // var id uint16
    // bytesBuffer := (data[0:2])
    // binary.Read(bytes.NewReader(bytesBuffer), binary.BigEndian, &id)
    // slice := data[12:]
    // i, j := 0, 0
    // var domainArr []string
    // var domainName string
    // for slice[i] != 0 {
    //     length :=  utils.BytesToInt(slice[i])
    //     domainArr = append(domainArr, string(slice[i + 1: i + 1 + length]))
    //     i = i + length + 1
    //     j++
    // }
    // domainName = strings.Join(domainArr, ".")
    var buffer bytes.Buffer 
    binary.Write(&buffer, binary.BigEndian, data[0:2])
    binary.Write(&buffer, binary.BigEndian, data[12:64])
    //fmt.Println(domainName, "resolve success!")
    return buffer.Bytes()
}

func Middleware(fromAddr string, toAddr string) {
    //fmt.Println("middleware start")
    udpAddr, err := net.ResolveUDPAddr("udp4", fromAddr)
    utils.ChkErr(err)
    udpConn, err2 := net.ListenUDP("udp", udpAddr)
    utils.ChkErr(err2)
    //fmt.Println("from", fromAddr)
    //udp没有对客户端连接的Accept函数
    for {
		buf := make([]byte, 256)
		_, clientAddr, errFromClient := udpConn.ReadFromUDP(buf)
		utils.ChkErr(errFromClient)
		go func() {
            BCConn, errSend2 := net.Dial("udp", toAddr)
            utils.ChkErr(errSend2)
            _, errSend3 := BCConn.Write(UdpHandle(buf))
            utils.ChkErr(errSend3)
            //fmt.Println("Package send to BlockChain!")

            msg := make([]byte, 512)
            _, errRecv := BCConn.Read(msg)
            //fmt.Println(msg)
            utils.ChkErr(errRecv)
            _, errToClient := udpConn.WriteToUDP(msg, clientAddr)
            utils.ChkErr(errToClient)
            //fmt.Println()
            //fmt.Println("Package reply Client!")
        }()
    }
}
