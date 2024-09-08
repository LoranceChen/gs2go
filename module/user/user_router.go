package user

import (
	"errors"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"gs2go/module"
	"gs2go/proto_define"
)

const (
	SIGN_UP       byte = 1
	MULTIPLE_CALL byte = 2
	HELLO         byte = 3
)

type PbMessage struct {
	field1 string
	field2 int
}
type UserRouter struct {
	conn           *websocket.Conn
	actionHandlers map[module.ActionName]func(message []byte) (proto.Message, error) // key: router name, value: output
	Name           module.RouterName
}

func (r *UserRouter) WsPbActionHandler(action module.ActionName, pb []byte) (proto.Message, error) {
	handler, ok := r.actionHandlers[action]
	if !ok {
		errMsg := errors.New("not found route: " + string(action.Name))
		return nil, errMsg
	}
	resultMsg, err := handler(pb)
	if err != nil {
		return nil, errors.Join(errors.New("handle resultMsg"), err)
	}

	return resultMsg, nil
}

func (r *UserRouter) RouterName() module.RouterName {
	return r.Name
}

func hello(message *proto_define.HelloRequest) (*proto_define.HelloResponse, error) {
	return &proto_define.HelloResponse{
		Echo:     "echo: " + message.Msg,
		Sequence: message.Sequence,
	}, nil
}

func NewUserRouter(conn *websocket.Conn) *UserRouter {
	m := make(map[module.ActionName]func(message []byte) (proto.Message, error), 128)

	parsedHelloMsg := func(message []byte) (proto.Message, error) {
		request := &proto_define.HelloRequest{}
		proto.Unmarshal(message, request)
		response, err := hello(request)
		return response, err
	}

	m[module.ActionName{HELLO}] = parsedHelloMsg
	userRouter := &UserRouter{conn: conn, actionHandlers: m, Name: module.RouterName{module.USER_SERVICE}}

	return userRouter
}
