package tool

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"path"
	"strings"

	"middle/library/tool/region"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	//ErrorEmptyInterfaceAddrs 定义错误，无法找到网卡信息时
	ErrorEmptyInterfaceAddrs = fmt.Errorf("empty found in InterfaceAddrs")
)

// IP 单例ip，并且导出
var IP = &ip{}
var ipServ = &region.IP2Region{}
var ipDataPath string

func init() {
	dir, err := Path.ExecRootPath()
	if err != nil {
		panic(err)
	}
	ipDataPath = path.Join(dir, "manifest/config", "ip2region.db")
}

type ip struct {
}

// for test
func (*ip) setIPDbDir(dir string) {
	ipDataPath = path.Join(dir, "manifest/config", "ip2region.db")
}

// GetAddr 根据IP地址
func (*ip) GetAddr(ctx context.Context, ipv4 string) (region.IPInfo, error) {
	// 线程安全的
	if !ipServ.Initialized() {
		err := ipServ.Init(ipDataPath)
		if err != nil {
			g.Log().Warningf(ctx, "ip.GetAddr init path=%v, err=%v", ipDataPath, err.Error())
			return region.IPInfo{}, err
		}
		g.Log().Info(ctx, "ip.GetAddr init success")
	}
	return ipServ.MemorySearch(ipv4)
}

// LocalIPv4s 获取本机局域网ipv4地址
func (*ip) LocalIPv4s() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", ErrorEmptyInterfaceAddrs
}

// IsIPV4 判断字符串是否是ipv4
func (*ip) IsIPV4(ipv4 string) bool {
	address := net.ParseIP(ipv4)
	return address == nil
}

// IsLanIP 判断字符串是否是内网IP
func (*ip) IsLanIP(ipv4 string) bool {
	ip := net.ParseIP(ipv4)
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return false
}

// IsChineseMainland 是否是中国大陆
func (receiver *ip) IsChineseMainland(ctx context.Context, ipStr string) bool {
	tmpRegion, _ := receiver.GetAddr(ctx, ipStr)
	if strings.Contains(tmpRegion.Country, "中国") {
		if !strings.Contains(tmpRegion.Province, "香港") && !strings.Contains(tmpRegion.Province, "澳门") && !strings.Contains(tmpRegion.Province, "台湾") {
			return true
		}
	}
	return false
}

// IsSingapore 是否是新加坡
func (receiver *ip) IsSingapore(ctx context.Context, ipStr string) bool {
	tmpRegion, _ := receiver.GetAddr(ctx, ipStr)
	return strings.Contains(tmpRegion.Country, "新加坡")
}

// IsConfineCountry 是否限制注册国家
func (receiver *ip) IsConfineCountry(ctx context.Context, ipStr string) bool {
	tmpRegion, _ := receiver.GetAddr(ctx, ipStr)
	//if strings.Contains(tmpRegion.Country, "中国") {
	//	if !strings.Contains(tmpRegion.Province, "香港") && !strings.Contains(tmpRegion.Province, "澳门") && !strings.Contains(tmpRegion.Province, "台湾") {
	//		return true
	//	}
	//} else if strings.Contains(tmpRegion.Country, "新加坡") {
	//	return true
	//}
	if strings.Contains(tmpRegion.String(), "美国") ||
		strings.Contains(tmpRegion.String(), "尼日利亚") ||
		strings.Contains(tmpRegion.String(), "肯尼亚") ||
		strings.Contains(tmpRegion.String(), "加纳") ||
		strings.Contains(tmpRegion.String(), "乌干达") ||
		strings.Contains(tmpRegion.String(), "南非") ||
		strings.Contains(tmpRegion.String(), "坦桑尼亚") ||
		strings.Contains(tmpRegion.String(), "赞比亚") ||
		strings.Contains(tmpRegion.String(), "埃塞俄比亚") ||
		strings.Contains(tmpRegion.String(), "纳米比亚") ||
		strings.Contains(tmpRegion.String(), "牙买加") ||
		strings.Contains(tmpRegion.String(), "喀麦隆") ||
		strings.Contains(tmpRegion.String(), "津巴布韦") ||
		strings.Contains(tmpRegion.String(), "博茨瓦纳") ||
		strings.Contains(tmpRegion.String(), "冈比亚") {
		return false
	}
	return true
}

func ParseIP(addr string) uint32 {
	var ip uint32 = 0
	ipv := net.ParseIP(addr)
	if ipv != nil {
		if strings.Contains(addr, ".") { //ipv4
			ip = binary.BigEndian.Uint32(ipv.To4())
		} else if strings.Contains(addr, ":") { // ipv6
			ip = binary.BigEndian.Uint32(gconv.Bytes(addr))
		}
	}
	return ip
}

// LocalIPv4s 获取本机局域网ipv4地址
func LocalIPv4s() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", ErrorEmptyInterfaceAddrs
}

func PublicIPv4s() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && IsPublicIP(ipnet.IP) {
			return ipnet.IP.String(), nil
		}
	}

	return "", ErrorEmptyInterfaceAddrs
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
