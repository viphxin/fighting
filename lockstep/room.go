package lockstep

import (
	"errors"
	"fighting/conf"
	"fighting/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"github.com/viphxin/xingo/utils"
	"sync"
	"time"
)

type RoomState uint8

const (
	WAITING RoomState = iota
	RUNNING
	ROOMEND
)

type Room struct {
	PlayerNumGen int32
	RoomId       string
	PropertyBag  map[string]interface{}
	Climit       int32
	Players      map[string]*Player
	RoomLock     sync.RWMutex
	//step queue
	IsStartLoopPush bool
	Seed            int32
	//StepQueue        chan *pb.UserInputData
	StepQueue *QuickSlice
	Stepnum   int32
	State     RoomState
	//step 补偿
	AvgPerTTL int64
	MSendTime time.Time
	//long(1.0/self.frameSpeed*10000000)
	FrameTickLength int64
}

func NewRoom(roomId string, cc int32) *Room {
	return &Room{
		PlayerNumGen:    0,
		RoomId:          roomId,
		PropertyBag:     make(map[string]interface{}),
		Climit:          cc,
		Players:         make(map[string]*Player),
		IsStartLoopPush: false,
		Seed:            1000,
		//StepQueue:       make(chan *pb.UserInputData, 128),
		StepQueue:       NewQuickSlice(0, 4),
		Stepnum:         -1,
		State:           WAITING,
		AvgPerTTL:       30000000,
		FrameTickLength: int64(1.0 / float64(utils.GlobalObject.FrameSpeed) * 1000000000),
	}
}

func (this *Room) GetPlayerCount() int {
	return len(this.Players)
}

func (this *Room) IsFull() bool {
	if int32(this.GetPlayerCount()) == this.Climit {
		return true
	} else {
		return false
	}
}

func (this *Room) IsAllInStateN(n PLAYERSTATE) bool {
	for _, p := range this.Players {
		if p.GetState() != n {
			return false
		}
	}
	return true
}

func (this *Room) IsEmpty() bool {
	if this.GetPlayerCount() == 0 {
		return true
	} else {
		return false
	}
}

func (this *Room) GetOfflineCount() int {
	var count int = 0
	for _, p := range this.Players {
		if !p.IsOnline() {
			count += 1
		}
	}
	return count
}

func (this *Room) IsAllReady() bool {
	if this.IsFull() && this.IsAllInStateN(ALREADY) {
		return true
	} else {
		return false
	}
}

func (this *Room) GetRid(uid string) (int32, error) {
	p, ok := this.Players[uid]
	if ok {
		return p.Rid, nil
	} else {
		return 0, errors.New("no player in room")
	}
}

func (this *Room) GetPlayer(uid string) (*Player, error) {
	p, ok := this.Players[uid]
	if ok {
		return p, nil
	} else {
		return p, errors.New("no player in room")
	}
}

func (this *Room) Broadcast(msgId uint32, data proto.Message) {
	for _, p := range this.Players {
		p.SendMsg(msgId, data)
	}
}

func (this *Room) Multicast(msgId uint32, data proto.Message, excludes map[string]bool) {
	for _, p := range this.Players {
		uid, _ := p.GetUid()
		_, ok := excludes[uid]
		if !ok {
			p.SendMsg(msgId, data)
		}

	}
}

func (this *Room) SendReady() {
	this.Broadcast(3, nil)
}

func (this *Room) SendMsg(uid string, msgId uint32, data proto.Message) {
	p, err := this.GetPlayer(uid)
	if err == nil {
		p.SendMsg(msgId, data)
	}
}

func (this *Room) JoinRoom(fconn iface.Iconnection) error {
	// this.RoomLock.Lock()
	// defer this.RoomLock.Unlock()

	uid, err := fconn.GetProperty("uid")
	if err == nil {
		if uidString, ok := uid.(string); ok {
			if p, err := this.GetPlayer(uidString); err == nil {
				//断线重连
				p.ReConnection(fconn)
				logger.Debug(fmt.Sprintf("reconnectioned ======== %s", uidString))
			} else {
				if int32(this.GetPlayerCount()) < this.Climit {
					this.PlayerNumGen += 1
					p := NewPlayer(fconn, this.PlayerNumGen)
					this.Players[uidString] = p
					//test发送PlayerJionMsg, InitRoomMsg
					playerInitRoomMsg := &pb.InitRoomMsg{
						Seed:      this.Seed,
						LocalID:   p.Rid,
						MaxPlayer: this.Climit,
					}

					this.SendMsg(uidString, 5, playerInitRoomMsg)

					playerJoinMsg := &pb.PlayerJionMsg{
						PlayerLocalID: p.Rid,
					}
					this.Multicast(6, playerJoinMsg, map[string]bool{uidString: true})

					for _, temp := range this.Players {
						playerJoinMsg := &pb.PlayerJionMsg{
							PlayerLocalID: temp.Rid,
						}
						this.SendMsg(uidString, 6, playerJoinMsg)
					}

				} else {
					//lostconnection
					fconn.LostConnection()
					return errors.New("player cc limit!!!")
				}
			}
		}
	} else {
		//lostconnection
		fconn.LostConnection()
		return errors.New("no user login")
	}
	return nil
}

func (this *Room) LeaveRoom(uid string) {
	delete(this.Players, uid)
}

func (this *Room) GetProperty(key string) (interface{}, error) {
	value, ok := this.PropertyBag[key]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("no property in connection")
	}
}

func (this *Room) SetProperty(key string, value interface{}) {
	this.PropertyBag[key] = value
}

func (this *Room) RemoveProperty(key string) {
	delete(this.PropertyBag, key)
}

func (this *Room) DoLostConnection(uid string) bool {
	p, err := this.GetPlayer(uid)
	if err == nil {
		logger.Debug("room doLostConnection " + uid)
		p.LostConnection()
		logger.Debug(fmt.Sprintf("%d====%d", this.GetPlayerCount(), this.GetOfflineCount()))

		if this.GetPlayerCount() == this.GetOfflineCount() {
			//所有玩家都掉线了, 玩个屁啊, 直接删除房间得了
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

/*
删除房间前回收房间用过的资源
*/
func (this *Room) DelRoomHandle() {
	this.StopLoopPush()
}

func (this *Room) StopLoopPush() {
	if this.IsStartLoopPush {
		this.IsStartLoopPush = false
		logger.Debug("canceled loopPush of room " + this.RoomId)
	}
}

func (this *Room) StartLoopPush() {
	this.RoomLock.Lock()
	defer this.RoomLock.Unlock()

	if !this.IsStartLoopPush {
		logger.Debug("start loopPush!!!!")
		//发送ready的消息
		this.SendReady()
		this.MSendTime = time.Now()
		this.IsStartLoopPush = true //每个房间只启动一个goroutine
		this.LoopPush()
	}
}

func (this *Room) LoopPush() {
	go func() {
		for {
			//没有玩家后需要清除定时器
			if this.IsStartLoopPush && (this.GetPlayerCount() > this.GetOfflineCount()) {
				// st := time.Now()
				// //do task
				// this.Step()
				// //auto sleep time!!!!!!!!!!!!!!!!!!!
				// // logger.Debug(fmt.Sprintf("11111step run time: %d=======%d", td.Nanoseconds(), this.FrameTickLength))
				// logger.Debug(fmt.Sprintf("step cost total time: %f ms", time.Now().Sub(st).Seconds()*1000))
				// needSleepTime := this.FrameTickLength - time.Now().Sub(st).Nanoseconds() - this.AvgPerTTL
				// // logger.Debug(fmt.Sprintf("22222step run time: %d", needSleepTime))
				// if needSleepTime > 0 {
				// 	time.Sleep(time.Duration(needSleepTime))
				// 	// time.Sleep(time.Millisecond * 2)
				// }
				// // time.Sleep(time.Second * 2)

				//do task
				this.Step()
				time.Sleep(time.Millisecond * time.Duration(conf.ServerConfObj.StepPerMs))
			} else {
				//close(this.StepQueue)
				logger.Debug("LoopPush stoped successful!!!")
				return
			}
		}
	}()
}

func (this *Room) AddUserInput(rid int32, input int32) {
	// this.StepQueue <- &pb.UserInputData{
	// 	ID:   rid,
	// 	Data: input,
	// }
	this.StepQueue.Append(&pb.UserInputData{
		ID:   rid,
		Data: input,
	})
}

func (this *Room) Step() {
	this.Stepnum += 1
	var inputs []*pb.UserInputData = this.StepQueue.GetAll()

	// logger.Debug(inputs)
	collectUsersInput := &pb.CollectUsersInput{
		Step:           this.Stepnum,
		UsersInputData: inputs,
	}

	//send
	this.Broadcast(11, collectUsersInput)
	//auto adjust
	if len(inputs) == 0 {
		c_count := time.Now().Sub(this.MSendTime).Nanoseconds() / this.FrameTickLength

		for {
			// logger.Debug(fmt.Sprintf("%d >>>>>>%d", c_count, this.Stepnum))
			// logger.Debug(fmt.Sprintf("%d !!!!!!!!!!!!>>>>>>%d", c_count, int64(this.Stepnum)))
			if c_count >= int64(this.Stepnum) {
				this.Stepnum += 1
				collectUsersInput.Step = this.Stepnum
				this.Broadcast(11, collectUsersInput)
			} else {
				break
			}

		}
	}
}

// func (this *Room) Step() {
// 	inputCount := len(this.StepQueue)
// 	this.Stepnum += 1
// 	var inputs []*pb.UserInputData

// 	if inputCount > 0 {
// 		inputs = make([]*pb.UserInputData, 0, inputCount)
// 		for i := 1; i < inputCount; i += 1 {
// 			inputs = append(inputs, <-this.StepQueue)
// 		}
// 	}
// 	// logger.Debug(inputs)
// 	collectUsersInput := &pb.CollectUsersInput{
// 		Step:           this.Stepnum,
// 		UsersInputData: inputs,
// 	}
// 	//send
// 	this.Broadcast(11, collectUsersInput)
// }
