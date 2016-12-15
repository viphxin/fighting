package api

import (
	"fighting/lockstep"
	"fighting/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/fnet"
	"github.com/viphxin/xingo/logger"
	_ "time"
)

type FightingRouter struct {
}

/*
ping test
*/
func (this *FightingRouter) Api_0(request *fnet.PkgAll) {
	logger.Debug("call Api_0")
	// request.Fconn.SendBuff(0, nil)
	request.Fconn.Send(0, nil)
}

func (this *FightingRouter) Api_1(request *fnet.PkgAll) {
	msg := &pb.UserLogin{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		userId := msg.UserId
		tocken := msg.Accesstocken
		logger.Debug(fmt.Sprintf("user login uid: %s.tocken: %s", userId, tocken))
		if userId != "" {
			request.Fconn.SetProperty("uid", userId)
			resp := &pb.CommonResponse{
				State: 1,
			}
			request.Fconn.SendBuff(1, resp)
		} else {
			logger.Error("no userid found")
			request.Fconn.LostConnection()
		}
	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}

func (this *FightingRouter) Api_2(request *fnet.PkgAll) {
	uid, _ := request.Fconn.GetProperty("uid")
	if uidStr, ok := uid.(string); ok {
		roomId, err := lockstep.RoomManagerObj.JoinRandomRoom(request.Fconn, 2)
		if err == nil {
			if roomId != "" {
				request.Fconn.SetProperty("room_id", roomId)
				room := lockstep.RoomManagerObj.GetRoom(roomId)
				p, err := room.GetPlayer(uidStr)
				if err == nil {
					p.SetState(lockstep.CHECKED)
					resp := &pb.CommonResponse{
						State: 1,
					}
					request.Fconn.SendBuff(2, resp)
				}
			} else {
				logger.Error("add room error")
				request.Fconn.LostConnection()
			}
		} else {
			logger.Error(err)
		}

	}
}

func (this *FightingRouter) Api_4(request *fnet.PkgAll) {
	uid, _ := request.Fconn.GetProperty("uid")
	roomId, _ := request.Fconn.GetProperty("room_id")
	if roomIdStr, ok := roomId.(string); ok {
		room := lockstep.RoomManagerObj.GetRoom(roomIdStr)
		if room != nil {
			if uidStr, ok := uid.(string); ok {
				p, err := room.GetPlayer(uidStr)
				if err == nil {
					p.SetState(lockstep.ALREADY)
					resp := &pb.CommonResponse{
						State: 1,
					}
					request.Fconn.SendBuff(4, resp)
					//检测是否所有客户端资源加载完成
					if room.IsAllReady() {
						room.StartLoopPush()
					}
				}
			}

		}
	}
}

func (this *FightingRouter) Api_10(request *fnet.PkgAll) {
	uid, _ := request.Fconn.GetProperty("uid")
	roomId, _ := request.Fconn.GetProperty("room_id")
	if roomIdStr, ok := roomId.(string); ok {
		room := lockstep.RoomManagerObj.GetRoom(roomIdStr)
		if room != nil {
			if uidStr, ok := uid.(string); ok {
				p, err := room.GetPlayer(uidStr)
				if err == nil {
					msg := &pb.UserInput{}
					err := proto.Unmarshal(request.Pdata.Data, msg)
					if err == nil {
						room.AddUserInput(p.Rid, msg.Data)
					}
				}
			}

		}
	}
}
