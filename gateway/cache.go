package gateway


import (
	"sync"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/driver"
)



// Cache stores realtime industrial states.
//
// It is the realtime memory database of gateway.
//
type Cache struct {


	mu sync.RWMutex


	devices map[string]*DeviceState


}




// DeviceState represents one device realtime state.
//
type DeviceState struct {


	DeviceID string


	Points map[string]*PointState


}




// PointState represents one measurement.
//
type PointState struct {


	PointID string


	Value any


	Quality uint8


	Timestamp time.Time



	Metadata map[string]any


}





// NewCache creates cache.
//
func NewCache()*Cache{


	return &Cache{


		devices:make(
			map[string]*DeviceState,
		),


	}

}






// Update updates cache from driver event.
//
func (c *Cache) Update(
	event driver.Event,
){


	if event.DeviceID=="" ||
	   event.PointID=="" {


		return

	}



	c.mu.Lock()

	defer c.mu.Unlock()



	device,ok:=c.devices[event.DeviceID]


	if !ok {


		device=&DeviceState{


			DeviceID:event.DeviceID,


			Points:make(
				map[string]*PointState,
			),

		}


		c.devices[event.DeviceID]=device


	}



	point,ok:=device.Points[event.PointID]


	if !ok {


		point=&PointState{


			PointID:event.PointID,


			Metadata:make(
				map[string]any,
			),


		}


		device.Points[event.PointID]=point

	}




	point.Value=event.Value


	point.Quality=uint8(event.Quality)


	point.Timestamp=event.Timestamp



	for k,v:=range event.Metadata{


		point.Metadata[k]=v


	}


}





// GetDevice returns device state.
//
func (c *Cache) GetDevice(
	deviceID string,
)(*DeviceState,bool){



	c.mu.RLock()

	defer c.mu.RUnlock()



	device,ok:=c.devices[deviceID]


	if !ok {

		return nil,false

	}



	return cloneDevice(device),true

}







// GetPoint returns one point.
//
func (c *Cache) GetPoint(
	deviceID string,
	pointID string,
)(*PointState,bool){


	c.mu.RLock()

	defer c.mu.RUnlock()



	device,ok:=c.devices[deviceID]


	if !ok {

		return nil,false

	}



	point,ok:=device.Points[pointID]


	if !ok {


		return nil,false

	}



	return clonePoint(point),true

}







// Snapshot returns all realtime states.
//
func (c *Cache) Snapshot()map[string]*DeviceState{


	c.mu.RLock()

	defer c.mu.RUnlock()



	result:=make(
		map[string]*DeviceState,
	)


	for id,state:=range c.devices{


		result[id]=cloneDevice(state)

	}



	return result

}






func cloneDevice(
	src *DeviceState,
)*DeviceState{


	dst:=&DeviceState{


		DeviceID:src.DeviceID,


		Points:make(
			map[string]*PointState,
		),

	}



	for id,p:=range src.Points{


		dst.Points[id]=clonePoint(p)


	}


	return dst

}





func clonePoint(
	src *PointState,
)*PointState{


	dst:=*src


	if src.Metadata!=nil{


		dst.Metadata=make(
			map[string]any,
		)


		for k,v:=range src.Metadata{


			dst.Metadata[k]=v


		}

	}


	return &dst

}
