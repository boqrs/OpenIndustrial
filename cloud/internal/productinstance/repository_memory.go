package productinstance

import (
	"context"
	"sync"
)


type MemoryRepository struct {

	mu sync.RWMutex

	items map[string]*ProductInstance

}


func NewMemoryRepository() *MemoryRepository {

	return &MemoryRepository{

		items:make(
			map[string]*ProductInstance,
		),

	}

}



func (r *MemoryRepository) Create(
	ctx context.Context,
	instance *ProductInstance,
) error {


	r.mu.Lock()
	defer r.mu.Unlock()


	r.items[instance.ID]=instance


	return nil

}



func (r *MemoryRepository) GetBySN(
	ctx context.Context,
	sn string,
)(
	*ProductInstance,
	error,
){


	r.mu.RLock()
	defer r.mu.RUnlock()



	for _,item:=range r.items{


		if item.SN==sn{

			return item,nil

		}

	}



	return nil,ErrNotFound

}



func (r *MemoryRepository) Update(
	ctx context.Context,
	instance *ProductInstance,
) error {


	r.mu.Lock()
	defer r.mu.Unlock()


	r.items[instance.ID]=instance


	return nil

}