package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/drivers/can"
)

func main() {
	log.Println("Starting CAN example...")

	// 在 Linux 上, 你需要使用一个真实的 CAN 接口名称, 例如 "can0" 或 "vcan0".
	// 在这个例子中, 我们使用 "vcan0".
	const canInterface = "vcan0"

	// 创建一个新的 SocketCAN 适配器.
	// 在非 Linux 系统上, 这会返回一个存根(stub)适配器.
	adapter := can.NewSocketCANAdapter()

	// 定义连接配置.
	config := can.ConnectionConfig{
		Interface: canInterface,
	}

	// 连接到 CAN 总线.
	log.Printf("Connecting to CAN interface: %s", canInterface)
	err := adapter.Connect(context.Background(), config)
	if err != nil {
		// 当在 macOS 或其他非 Linux 系统上运行时, 存根适配器
		// 会返回这个特定的错误. 我们优雅地处理它.
		log.Printf("Failed to connect to CAN bus: %v", err)
		log.Println("This is expected if you are not running on a Linux system with a configured CAN interface.")
		log.Println("Example finished.")
		return
	}
	defer adapter.Disconnect(context.Background())

	log.Println("Successfully connected to CAN bus.")

	// 创建一个可以被取消的上下文.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动一个 goroutine 来读取和打印接收到的 CAN 报文.
	go func() {
		log.Println("Listening for incoming CAN frames...")
		frameChan := adapter.Receive()
		for {
			select {
			case frame, ok := <-frameChan:
				if !ok {
					log.Println("Frame channel closed.")
					return
				}
				log.Printf("Received CAN frame: ID=0x%X, DLC=%d, Data=%X, Extended=%v, RTR=%v",
					frame.ID, frame.DLC, frame.Data, frame.Extended, frame.RTR)
			case <-ctx.Done():
				log.Println("Stopping frame listener.")
				return
			}
		}
	}()

	// 启动一个 goroutine 每隔 2 秒发送一个样本 CAN 报文.
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sampleFrame := can.Frame{
					ID:       0x123,
					DLC:      8,
					Data:     []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
					Extended: false,
					RTR:      false,
				}
				log.Printf("Sending sample CAN frame: ID=0x%X", sampleFrame.ID)
				if err := adapter.Send(context.Background(), sampleFrame); err != nil {
					log.Printf("Error sending frame: %v", err)
				}
			case <-ctx.Done():
				log.Println("Stopping frame sender.")
				return
			}
		}
	}()

	// 等待一个信号来优雅地关闭程序.
	log.Println("Example is running. Press Ctrl+C to exit.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received, disconnecting...")
}