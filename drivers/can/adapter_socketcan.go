//go:build linux

package can

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)


type socketcanAdapter struct {

	mu sync.RWMutex

	fd int

	iface string

	connected bool

	frameCh chan Frame

	stop chan struct{}

}



func NewSocketCANAdapter() Adapter {

	return &socketcanAdapter{}

}



func (a *socketcanAdapter) Connect(
	ctx context.Context,
	cfg ConnectionConfig,
) error {


	a.mu.Lock()
	defer a.mu.Unlock()


	if a.connected {
		return ErrAlreadyConnected
	}


	fd, err := unix.Socket(
		unix.AF_CAN,
		unix.SOCK_RAW,
		unix.CAN_RAW,
	)

	if err != nil {
		return fmt.Errorf(
			"create socketcan socket failed:%w",
			err,
		)
	}



	iface, err := net.InterfaceByName(cfg.Interface)
	if err != nil {
		unix.Close(fd)
		return fmt.Errorf("get interface %s failed: %w", cfg.Interface, err)
	}

	addr := &unix.SockaddrCAN{
		Ifindex: iface.Index,
	}


	err = unix.Bind(fd, addr)

	if err != nil {

		unix.Close(fd)

		return fmt.Errorf(
			"bind socketcan failed:%w",
			err,
		)
	}



	a.fd = fd

	a.iface = cfg.Interface

	a.frameCh = make(chan Frame,100)

	a.stop = make(chan struct{})

	a.connected=true



	go a.readLoop()



	return nil

}



func (a *socketcanAdapter) Disconnect(
	ctx context.Context,
) error {


	a.mu.Lock()
	defer a.mu.Unlock()


	if !a.connected {
		return nil
	}



	close(a.stop)


	unix.Close(a.fd)



	close(a.frameCh)


	a.connected=false


	return nil
}



func (a *socketcanAdapter) IsConnected() bool {

	a.mu.RLock()
	defer a.mu.RUnlock()


	return a.connected

}



func (a *socketcanAdapter) Receive() <-chan Frame {

	return a.frameCh

}



func (a *socketcanAdapter) Send(
	ctx context.Context,
	frame Frame,
) error {


	if !a.IsConnected(){

		return ErrNotConnected

	}



	raw:=make([]byte,16)


	id:=frame.ID


	if frame.Extended{

		id |= unix.CAN_EFF_FLAG

	}


	if frame.RTR{

		id |= unix.CAN_RTR_FLAG

	}



	raw[0]=byte(id)
	raw[1]=byte(id>>8)
	raw[2]=byte(id>>16)
	raw[3]=byte(id>>24)


	raw[4]=frame.DLC


	copy(
		raw[8:],
		frame.Data,
	)



	a.mu.RLock()

	defer a.mu.RUnlock()



	_,err:=unix.Write(
		a.fd,
		raw,
	)


	return err

}




func (a *socketcanAdapter) readLoop(){

	buf:=make([]byte,16)


	for{


		select{

		case <-a.stop:
			return

		default:

		}



		n,err:=unix.Read(
			a.fd,
			buf,
		)


		if err!=nil{

			continue

		}


		if n<16{

			continue

		}



		canID:=uint32(buf[0]) |
			uint32(buf[1])<<8 |
			uint32(buf[2])<<16 |
			uint32(buf[3])<<24



		frame:=Frame{


			ID:
			canID & unix.CAN_EFF_MASK,


			DLC:
			buf[4],


			Timestamp:
			time.Now(),


			Extended:
			canID&unix.CAN_EFF_FLAG!=0,


			RTR:
			canID&unix.CAN_RTR_FLAG!=0,


			Data:
			append(
				[]byte{},
				buf[8:8+buf[4]]...,
			),

		}



		select{


		case a.frameCh<-frame:


		default:

			// 丢弃堵塞数据

		}


	}


}