package main

import (
	"fighting/api"
	"fighting/lockstep"
	"fmt"
	"github.com/viphxin/xingo/fserver"
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"github.com/viphxin/xingo/utils"

	_ "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	_ "runtime/pprof"
	_ "time"
)

func DoConnectionMade(fconn iface.Iconnection) {
	logger.Debug("111111111111111111111111")
}

func DoConnectionLost(fconn iface.Iconnection) {
	logger.Debug("222222222222222222222222")
	uid, _ := fconn.GetProperty("uid")
	roomId, _ := fconn.GetProperty("room_id")
	if roomIdStr, ok := roomId.(string); ok {
		room := lockstep.RoomManagerObj.GetRoom(roomIdStr)
		if room != nil {
			if uidStr, ok := uid.(string); ok {
				logger.Debug(fmt.Sprintf("user %s lostconnection", uidStr))
				if room.DoLostConnection(uidStr) {
					logger.Debug("delete room. reason: room empty!!!")
					lockstep.RoomManagerObj.LeaveRoom(roomIdStr, uidStr)
				}
			}

		}
	}
}

func main() {
	s := fserver.NewServer()

	//add api ---------------start
	FightingRouterObj := &api.FightingRouter{}
	s.AddRouter(FightingRouterObj)
	//add api ---------------end
	//regest callback
	utils.GlobalObject.OnConnectioned = DoConnectionMade
	utils.GlobalObject.OnClosed = DoConnectionLost

	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6061", nil))
	// 	// for {
	// 	// 	time.Sleep(time.Second * 10)
	// 	// 	fm, err := os.OpenFile("./memory.log", os.O_RDWR|os.O_CREATE, 0644)
	// 	// 	if err != nil {
	// 	// 		fmt.Println(err)
	// 	// 	}
	// 	// 	pprof.WriteHeapProfile(fm)
	// 	// 	fm.Close()
	// 	// }
	// }()

	s.Start()
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("=======", sig)
	s.Stop()
}
