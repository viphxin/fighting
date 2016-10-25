package main

import (
	"github.com/golang/protobuf/proto"
	"log"
	"pb"
)

func main() {
	test := &pb.UserLogin{ // 使用辅助函数设置域的值
		UserId:       "hello",
		Accesstocken: "hello",
	} // 进行编码

	data, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	} // 进行解码
	newTest := &pb.UserLogin{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	} // 测试结果
	if test.UserId != newTest.UserId {
		log.Fatalf("data mismatch %q != %q", test.UserId, newTest.UserId)
	}
}
