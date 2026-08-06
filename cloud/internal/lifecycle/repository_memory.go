package lifecycle

import (
	"context"
	"sync"
)



type MemoryRepository struct {


	mu sync.RWMutex


	events map[string][]Event


}



func NewMemoryRepository()*MemoryRepository{


	return &MemoryRepository{


		events:
		make(
			map[string][]Event,
		),


	}

}




func (r *MemoryRepository) SaveEvent(
	ctx context.Context,
	event *Event,
)error{


	r.mu.Lock()
	defer r.mu.Unlock()



	r.events[
		event.ProductInstanceID,
	]=append(
		r.events[
			event.ProductInstanceID,
		],
		*event,
	)


	return nil

}




func (r *MemoryRepository) GetEvents(
	ctx context.Context,
	productInstanceID string,
)(
[]Event,error){


	r.mu.RLock()
	defer r.mu.RUnlock()



	return r.events[
		productInstanceID,
	],nil

}