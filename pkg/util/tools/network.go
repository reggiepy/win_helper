package tools

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

// WakeOnLan 发送Magic Packet唤醒指定MAC地址的主机
func WakeOnLan(macAddress string, broadcastAddress string, port int) bool {
	// 解析MAC地址
	macAddress = strings.ReplaceAll(macAddress, ":", "")
	macAddress = strings.ReplaceAll(macAddress, "-", "")
	if len(macAddress) != 12 {
		fmt.Printf("Error parsing MAC address: invalid length\n")
		return false
	}

	macBytes := make([]byte, 6)
	for i := 0; i < 6; i++ {
		byteValue, err := fmt.Sscanf(macAddress[2*i:2*i+2], "%X", &macBytes[i])
		if err != nil || byteValue != 1 {
			fmt.Printf("Error parsing MAC address\n")
			return false
		}
	}

	// 构造Magic Packet
	magicPacket := make([]byte, 102)
	for i := 0; i < 6; i++ {
		magicPacket[i] = 0xff
	}
	for i := 1; i <= 16; i++ {
		copy(magicPacket[i*6:], macBytes)
	}

	// 创建UDP套接字
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.ParseIP(broadcastAddress),
		Port: port,
	})
	if err != nil {
		fmt.Printf("Error creating UDP connection: %v\n", err)
		return false
	}
	defer conn.Close()

	// 发送Magic Packet
	_, err = conn.Write(magicPacket)
	if err != nil {
		fmt.Printf("Error sending Magic Packet: %v\n", err)
		return false
	}

	return true
}

// GetLocalIP 获取本机的局域网IP地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("local IP not found")
}

// pingAndArp Ping指定IP地址，并使用gopsutil获取其MAC地址
func pingAndArp(ip string, wg *sync.WaitGroup, mu *sync.Mutex, result *map[string]string) {
	defer wg.Done()

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		fmt.Printf("Failed to create pinger: %v\n", err)
		return
	}

	pinger.Count = 1
	pinger.Timeout = time.Second
	pinger.SetPrivileged(true) // 设置为privileged模式，允许发送ICMP请求

	err = pinger.Run()
	if err != nil {
		fmt.Printf("Failed to ping %s: %v\n", ip, err)
		return
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		interfaces, err := net.Interfaces()
		if err != nil {
			fmt.Printf("Failed to get interfaces: %v\n", err)
			return
		}

		for _, iface := range interfaces {
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Printf("Failed to get addresses for interface %s: %v\n", iface.Name, err)
				continue
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && ipNet.Contains(net.ParseIP(ip)) {
					mu.Lock()
					(*result)[ip] = iface.HardwareAddr.String()
					mu.Unlock()
					return
				}
			}
		}
	}
}


// ScanNetwork 扫描指定CIDR网络范围内的所有IP地址，并获取其MAC地址
func ScanNetwork(networks []string) map[string]string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	result := make(map[string]string)

	for _, network := range networks {
		ip, ipnet, err := net.ParseCIDR(network)
		if err != nil {
			fmt.Printf("Error parsing CIDR %s: %v\n", network, err)
			continue
		}

		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
			wg.Add(1)
			go pingAndArp(ip.String(), &wg, &mu, &result)
		}
	}

	wg.Wait()
	return result
}

// incrementIP 增加一个IP地址
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
