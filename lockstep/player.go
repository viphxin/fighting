package lockstep

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/iface"
)

type PLAYERSTATE uint8

const (
	QUEUEED PLAYERSTATE = iota
	CHECKED
	ALREADY
	FIGHTING
	END
)

type Player struct {
	Fconn       iface.Iconnection
	Rid         int32
	PropertyBag map[string]interface{}
	MsgQueue    [][]byte //消息队列缓存，用于断线后游戏回放

	PlayerState PLAYERSTATE
}

func NewPlayer(fconn iface.Iconnection, rid int32) *Player {
	p := &Player{
		Fconn:       fconn,
		Rid:         rid,
		PropertyBag: make(map[string]interface{}),
		MsgQueue:    make([][]byte, 10),

		PlayerState: QUEUEED,
	}
	return p
}

func (this *Player) GetSid() (uint32, error) {
	if this.Fconn != nil {
		return this.Fconn.GetSessionId(), nil
	} else {
		return 0, errors.New("the player is offline")
	}
}

func (this *Player) GetUid() (string, error) {
	if this.Fconn != nil {
		uid, err := this.Fconn.GetProperty("uid")
		if err == nil {
			return uid.(string), nil
		} else {
			return "", errors.New("the player not have attr uid")
		}

	} else {
		return "", errors.New("the player is offline")
	}
}

func (this *Player) IsOnline() bool {
	if _, err := this.GetSid(); err == nil {
		return true
	} else {
		return false
	}
}

func (this *Player) SetState(s PLAYERSTATE) {
	this.PlayerState = s
}

func (this *Player) GetState() PLAYERSTATE {
	return this.PlayerState
}

func (this *Player) SendMsg(msgId uint32, data proto.Message) {
	if this.Fconn != nil {
		this.Fconn.Send(msgId, data)
	}
}

func (this *Player) GetProperty(key string) (interface{}, error) {
	value, ok := this.PropertyBag[key]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("no property in connection")
	}
}

func (this *Player) SetProperty(key string, value interface{}) {
	this.PropertyBag[key] = value
}

func (this *Player) RemoveProperty(key string) {
	delete(this.PropertyBag, key)
}

func (this *Player) LostConnection() {
	this.Fconn = nil
}

func (this *Player) ReConnection(fconn iface.Iconnection) {
	this.Fconn = fconn
}
