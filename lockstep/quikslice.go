package lockstep

import (
	"fighting/pb"
	_ "fmt"
	"sync"
)

type QuickSlice struct {
	Len  int
	Cap  int
	data []*pb.UserInputData
	sync.Mutex
}

func NewQuickSlice(l int, c int) *QuickSlice {
	return &QuickSlice{
		Len:  l,
		Cap:  c,
		data: make([]*pb.UserInputData, l, c),
	}
}

func (this *QuickSlice) Append(v *pb.UserInputData) {
	this.Lock()
	defer this.Unlock()

	this.data = append(this.data, v)
}

func (this *QuickSlice) GetAll() (d []*pb.UserInputData) {
	this.Lock()
	defer this.Unlock()
	if len(this.data) > 0 {
		d, this.data = this.data[:len(this.data)], make([]*pb.UserInputData, this.Len, this.Cap)
	}

	return
}
