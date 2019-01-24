package bcserver;

import (
	//"fmt"
	"net"
	"utils"
	"io/ioutil"
	"bytes"
	"strings"
	"encoding/binary"
	"github.com/json-iterator/go"
)

//接受json数据结构体
type Response struct {
    Code  uint32
    Msg   string
    Data  ResBody
}

type ResBody struct {
    Auth_count     uint16
    Extra_count    uint16
    Auth_details   []DNSRecord
    Extra_details  []DNSRecord
}

type DNSRecord struct {
    Name      string
    Type      uint16
    Class     uint16
    TTL       uint32
    RDLength  uint16
    RData     RData
}

type RData struct {
    Addr	  string
    Name_server   string
}

//DNS消息头部
type DnsHeader struct {
	Id									uint16  //Message ID 2字节
    Flag								uint16  //头部内容 2字节
	Qdcount, Ancount, Nscount, Arcount	uint16  //标识计数 各2字节 Qdcount请求部分的条目数、Ancount相应部分资源记录数、Nscount权威部分域名资源记录数、Arcount额外部分资源记录数
}

func (header *DnsHeader) SetFlag(QR uint16, OperationCode uint16, AuthoritativeAnswer uint16, Truncation uint16, RecursionDesired uint16, RecursionAvailable uint16, ResponseCode uint16) {
    header.Flag = QR<<15 + OperationCode<<11 + AuthoritativeAnswer<<10 + Truncation<<9 + RecursionDesired<<8 + RecursionAvailable<<7 + ResponseCode
}
//Bits 包括QR 1bit 标识请求/应答、OPCODE 4bit 标识请求类型、 AA 1bit 只在响应中有效，标识是否为权威、TC 1bit 是否截断、RD 1bit 是否递归、RA 1bit 是否支持递归、 RCODE 4bit 只在响应中标注，标识响应消息类型

//请求部分
type QuerySection struct {
    QueryType   uint16
    QueryClass  uint16
}

//响应部分  (不含RDATA)
type ResponseSection struct {
    ResponseName        []byte
    ResponseType        uint16
    ResponseClass       uint16
    ResponseTTL         uint32
    ResponseRDLength    uint16
    ResponseRData       ResponseRData
}

//RDATA部分
type ResponseRData struct {
    ResponseAddr        uint32
    ResponseNameServer  []byte
}

//构建应答数据包
func BuildPacket(id uint16, queryName string, blockResponse ResBody) ([]byte){
    //响应头部
	responseHeader := DnsHeader {
		Id	   :	id,
		Qdcount:	1,
		Ancount:	0,
		Nscount:	blockResponse.Auth_count,
		Arcount:	blockResponse.Extra_count,
	}
    responseHeader.SetFlag(1, 0, 0, 0, 0, 0, 0)
    requestQuery := QuerySection {
        QueryType:  1,
        QueryClass:  1,
    }

    var buffer bytes.Buffer
    binary.Write(&buffer, binary.BigEndian, responseHeader)
    binary.Write(&buffer, binary.BigEndian, utils.ParseDomainName(queryName))
    binary.Write(&buffer, binary.BigEndian, requestQuery)
    binary.Write(&buffer, binary.BigEndian, BuildDNSRecord(blockResponse.Auth_details, queryName))
    binary.Write(&buffer, binary.BigEndian, BuildDNSRecord(blockResponse.Extra_details, queryName))
    //fmt.Println("resolve success!")
	return buffer.Bytes()
}

func BuildDNSRecord(details []DNSRecord, queryName string) []byte {
    var buffer bytes.Buffer
    for _, value := range details {
        rdata := value.RData
        responseName := value.Name
        var responseNameBytes []byte
        if strings.Contains(responseName, queryName) {
            responseNameBytes = utils.ParseDomainName(responseName);  //用指针代替，后续优化
        } else {
            responseNameBytes = utils.ParseDomainName(responseName);
        }

        var (
            ResponseRDLength        uint16
        )
        binary.Write(&buffer, binary.BigEndian, responseNameBytes)
        binary.Write(&buffer, binary.BigEndian, value.Type)
        binary.Write(&buffer, binary.BigEndian, value.Class)
        binary.Write(&buffer, binary.BigEndian, value.TTL)
        if value.Type == 2 {
            ResponseNameServerBytes := utils.ParseDomainName(rdata.Name_server)
            ResponseRDLength = uint16(len(ResponseNameServerBytes))
            binary.Write(&buffer, binary.BigEndian, ResponseRDLength)
            binary.Write(&buffer, binary.BigEndian, ResponseNameServerBytes)
        } else if  value.Type == 1 {
            ResponseAddrBytes := utils.ParseIPv4(rdata.Addr)
            ResponseRDLength = 4
            binary.Write(&buffer, binary.BigEndian, ResponseRDLength)
            binary.Write(&buffer, binary.BigEndian, ResponseAddrBytes)
        }
    }
    return buffer.Bytes()
}

func BCServerHandle(queryName string) Response {
	bytes, err := ioutil.ReadFile("bcserver/data.json")
	utils.ChkErr(err)
	var res Response
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    errJson := json.Unmarshal(bytes, &res)
	utils.ChkErr(errJson)
	return res
}

func QueryHandle(queryData []byte) (uint16, string) {
	var id uint16
	var queryName string
	binary.Read(bytes.NewReader(queryData[0:2]), binary.BigEndian, &id)
	queryName = utils.ParseBytesToDomainName(queryData[2:])
	//fmt.Println("queryName is ", queryName)
	return id, queryName
}

func BCServer(addr string) {
	//fmt.Println("bcserver start")
	udpAddr, err1 := net.ResolveUDPAddr("udp4", addr)
    utils.ChkErr(err1)
    //监听端口
    udpConn, err2 := net.ListenUDP("udp", udpAddr)
	utils.ChkErr(err2)
	for {
        buf := make([]byte, 512)
		_, clientAddr, err3 := udpConn.ReadFromUDP(buf)
		utils.ChkErr(err3)
		go func() {
            id, queryName := QueryHandle(buf)
			res := BCServerHandle(queryName)
            bytesBuffer := BuildPacket(id, queryName, res.Data)
            //fmt.Println(bytesBuffer)
			udpConn.WriteToUDP(bytesBuffer, clientAddr)
		}()
	}
}
