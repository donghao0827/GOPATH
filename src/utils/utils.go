package utils;

import (
    "encoding/binary"
    "strings"
    "strconv"
    "bytes"
    "log"
)

//域名字符串转btye数组
func ParseDomainName(domain string) []byte {
	var (
		buffer		bytes.Buffer
		segments	[]string = strings.Split(domain, ".")
	)
	for _, seg := range segments {
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	binary.Write(&buffer, binary.BigEndian, byte(0x00))
	return buffer.Bytes()
}

//IPv4地址转整型
func ParseIPv4(ip string) uint32 {
    var (
        segments    []string = strings.Split(ip, ".")
    )
    var sum uint32
    seg0, _ := strconv.Atoi(segments[0])
    seg1, _ := strconv.Atoi(segments[1])
    seg2, _ := strconv.Atoi(segments[2])
    seg3, _ := strconv.Atoi(segments[3])

    sum += uint32(seg0) << 24
    sum += uint32(seg1) << 16
    sum += uint32(seg2) << 8
    sum += uint32(seg3)
    return sum
}

//检错
func ChkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}