package lockstep

import (
	"fmt"
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"sync"
	"time"
)

type RoomManager struct {
	Rooms        map[string]*Room
	CanJoinRooms []*Room
	ManagerLock  sync.RWMutex
	RoomGen      uint32
}

var RoomManagerObj *RoomManager

func init() {
	RoomManagerObj = &RoomManager{
		Rooms:        make(map[string]*Room),
		CanJoinRooms: make([]*Room, 0, 10),
		RoomGen:      0,
	}

	//statistics room
	go func() {
		for {
			time.Sleep(time.Second * 20)
			logger.Info(fmt.Sprintf("total rooms: %d.total CanJoinRooms: %d", RoomManagerObj.Count(), len(RoomManagerObj.CanJoinRooms)))
		}

	}()
}

func (this *RoomManager) Count() int {
	return len(this.Rooms)
}

func (this *RoomManager) JoinRandomRoom(fconn iface.Iconnection, cc int32) (string, error) {
	this.ManagerLock.Lock()
	defer this.ManagerLock.Unlock()

	var room *Room
	if len(this.CanJoinRooms) > 0 {
		index := len(this.CanJoinRooms) - 1
		room = this.CanJoinRooms[index]
		this.CanJoinRooms = this.CanJoinRooms[:index]
	} else {
		//create a room
		this.RoomGen += 1
		roomId := fmt.Sprintf("room_id_%d", this.RoomGen)
		room = NewRoom(roomId, cc)
		this.Rooms[roomId] = room
	}
	err := room.JoinRoom(fconn)
	if err == nil {
		if !room.IsFull() {
			this.CanJoinRooms = append(this.CanJoinRooms, room)
		}
		logger.Debug(room.RoomId)
		return room.RoomId, nil
	} else {
		return "", err
	}

}

func (this *RoomManager) GetRoom(roomId string) *Room {
	this.ManagerLock.RLock()
	defer this.ManagerLock.RUnlock()

	room, ok := this.Rooms[roomId]
	if ok {
		return room
	} else {
		return nil
	}
}

func (this *RoomManager) removeRoom(roomId string) {
	this.ManagerLock.Lock()
	defer this.ManagerLock.Unlock()
	delete(this.Rooms, roomId)
	canJoinRoomLen := len(this.CanJoinRooms)
	for i := 0; i < canJoinRoomLen; i += 1 {
		logger.Debug(fmt.Sprintf("%s------%s", this.CanJoinRooms[i].RoomId, roomId))
		if this.CanJoinRooms[i].RoomId == roomId {
			if i == canJoinRoomLen-1 {
				this.CanJoinRooms = this.CanJoinRooms[:i]
			} else {
				this.CanJoinRooms = append(this.CanJoinRooms[:i], this.CanJoinRooms[i+1:]...)
			}

			return
		}
	}
}

func (this *RoomManager) LeaveRoom(roomId string, uid string) {
	logger.Debug(fmt.Sprintf("roomId: %s,uid: %s.", roomId, uid))
	room := this.GetRoom(roomId)
	logger.Debug(room)
	if room != nil {
		room.LeaveRoom(uid)
		//如果房间内没有玩家了，需要删除房间
		if room.GetPlayerCount() == room.GetOfflineCount() {
			room.StopLoopPush()
			this.removeRoom(roomId)
		}
	}

	logger.Debug(this.Rooms)
	logger.Debug("this.CanJoinRooms this.CanJoinRooms this.CanJoinRooms")
	logger.Debug(this.CanJoinRooms)
}
