package tray

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type MethodTray struct {
	request  *protogen.Message
	response *protogen.Message
}

func NewMethodTray() *MethodTray {
	return &MethodTray{}
}

func (mt *MethodTray) Request() *protogen.Message {
	return mt.request
}

func (mt *MethodTray) Response() *protogen.Message {
	return mt.response
}

func (mt *MethodTray) SetResponse(response *protogen.Message) {
	mt.response = response
}

func (mt *MethodTray) SetRequest(request *protogen.Message) {
	mt.request = request
}
