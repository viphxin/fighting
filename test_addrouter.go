package main

import (
	"fmt"
	"fnet"
	"reflect"
	_ "runtime"
	"strconv"
	"strings"
)

type RouterA struct {
	f1 int
	f2 string
	f3 []int
}

func (this *RouterA) Func_1(name string, a *RouterA) {
	fmt.Println("exec Func_1!!!!" + name + strconv.Itoa(a.f1))
}

func (this *RouterA) Func_2() {
	fmt.Println("exec Func_2!!!!")
}

type FightingRouter struct {
}

/*
ping test
*/
func (this *FightingRouter) Api_0(request *fnet.PkgAll) {
	//request.Fconn.Send(0, nil)
	fmt.Println("exec api_0!!!!")
}

func AddRouter(router interface{}) {
	Apis := make(map[int]reflect.Value)
	value := reflect.ValueOf(router)
	tp := value.Type()
	for i := 0; i < value.NumMethod(); i += 1 {
		//fmt.Println(object.Method(i))
		//fmt.Println(runtime.FuncForPC(object.Method(i).Pointer()).Name())
		fmt.Printf("method[%d]%s\n", i, tp.Method(i).Name)
		name := tp.Method(i).Name
		k := strings.Split(name, "_")
		index, err := strconv.Atoi(k[1])
		if err != nil {
			panic("error api: " + name)
		}
		Apis[index] = value.Method(i)
	}

	//exec test
	for i := 0; i < 100; i += 1 {
		// Apis[1].Call([]reflect.Value{reflect.ValueOf("huangxin"), reflect.ValueOf(router)})
		// Apis[2].Call([]reflect.Value{})
		Apis[0].Call([]reflect.Value{reflect.ValueOf(&fnet.PkgAll{})})
	}
}

func main() {
	// AddRouter(&RouterA{
	// 	f1: 100099,
	// })
	AddRouter(&FightingRouter{})
}
