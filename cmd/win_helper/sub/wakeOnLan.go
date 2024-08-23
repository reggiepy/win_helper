package sub

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"win_helper/pkg/util/tools"
)

func newWakeOnLan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wake-on-lan",
		Short: "wake-on-lan",
		Long:  `wake-on-lan`,
		RunE: func(cmd *cobra.Command, args []string) error {
			macAddress := "1C:83:41:78:28:F5"    // 要唤醒的主机的MAC地址
			broadcastAddress := "10.201.127.255" // 广播地址
			port := 9                            // UDP端口

			success := tools.WakeOnLan(macAddress, broadcastAddress, port)
			if success {
				fmt.Println("Magic Packet sent successfully.")
			} else {
				fmt.Println("Failed to send Magic Packet.")
			}
			return nil
		},
	}

	return cmd
}

func newWakeOnLanScan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan-mac",
		Short: "scan-mac",
		Long:  `scan-mac`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 获取所有网卡信息
			interfaces, err := net.Interfaces()
			if err != nil {
				fmt.Println("Error getting interfaces:", err)
				return nil
			}

			// 存储要扫描的网卡IP地址和CIDR范围
			var networks []string

			// 遍历每个网卡
			for _, iface := range interfaces {
				// 只考虑IPv4地址
				if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
					continue // 跳过未启用的和回环接口
				}

				addrs, err := iface.Addrs()
				if err != nil {
					fmt.Printf("Error getting addresses for interface %s: %v\n", iface.Name, err)
					continue
				}

				// 遍历网卡的每个地址
				for _, addr := range addrs {
					ipNet, ok := addr.(*net.IPNet)
					if !ok || ipNet.IP.IsLoopback() {
						continue // 跳过非IPNet类型和回环地址
					}

					// 构造CIDR网络范围，例如：192.168.1.0/24
					network := ipNet.IP.String() + "/24"
					networks = append(networks, network)
				}
			}

			// 扫描网络范围
			fmt.Println("Scanning networks:", networks)
			result := tools.ScanNetwork(networks)

			// 输出扫描结果
			fmt.Println("Scan results:")
			for ip, mac := range result {
				fmt.Printf("IP: %s, MAC: %s\n", ip, mac)
			}
			return nil
		},
	}

	return cmd
}
