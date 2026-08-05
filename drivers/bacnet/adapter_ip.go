package bacnet

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/alexbeltran/gobacnet"
	gb "github.com/alexbeltran/gobacnet/types"
)

// connectIPWithDiscovery 处理 BACnet/IP 连接和设备发现逻辑。
// 它将连接信息设置到传入的 GoBACnetAdapter 实例中。
func connectIPWithDiscovery(ctx context.Context, a *GoBACnetAdapter, cfg ConnectionConfig) error {
	client, err := gobacnet.NewClient(
		"", // 本地地址，留空表示绑定到所有可用接口
		int(cfg.Port),
	)
	if err != nil {
		return fmt.Errorf("create bacnet client: %w", err)
	}

	addr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%d", cfg.DeviceAddress, cfg.Port),
	)
	if err != nil {
		client.Close()
		return err
	}

	a.client = client
	a.target = gb.UDPToAddress(addr)

	// 使用 -1, -1 进行广播 WhoIs，以发现网络中的所有设备
	// 修复点：添加 ctx 参数
	devices, err := client.WhoIs(-1, -1)
	if err != nil {
		client.Close()
		return err
	}

	if len(devices) == 0 {
		client.Close()
		return ErrDeviceNotFound // 未发现任何设备
	}

	var found bool
	for _, dev := range devices {
		if cfg.DeviceID != 0 {
			// 如果配置中指定了 DeviceID，则查找匹配的设备
			if dev.ID.Instance == gb.ObjectInstance(cfg.DeviceID) {
				a.device = dev
				found = true
				break
			}
		} else {
			// 如果没有指定 DeviceID，则使用第一个发现的设备
			a.device = dev
			found = true
			break
		}
	}

	if !found {
		client.Close()
		return ErrDeviceNotFound // 未找到指定设备
	}

	// 设置连接信息
	a.setConnected(ConnectionInfo{
		Mode:        cfg.Mode,
		Address:     cfg.DeviceAddress,
		Port:        cfg.Port,
		DeviceID:    cfg.DeviceID,
		Description: "BACnet/IP",
	})

	log.Printf(
		"BACnet connected device=%d address=%s",
		a.device.ID.Instance,
		cfg.DeviceAddress,
	)

	return nil
}