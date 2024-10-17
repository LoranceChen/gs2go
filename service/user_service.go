package service

import (
	"errors"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	patch_data "gs2go/ck-patch-data-protocol/protobuf"
	"gs2go/proto_define"
	myprotobuf "gs2go/service/myprotobuf"
)

const (
	SIGN_UP       ActionName = 1
	MULTIPLE_CALL ActionName = 2
	HELLO         ActionName = 3
	HELLO3        ActionName = 4
)

type UserRouter struct {
	conn           *websocket.Conn
	actionHandlers map[ActionName]func(message []byte) (proto.Message, error) // key: router name, value: output
	Name           RouterName
}

func (r *UserRouter) WsPbActionHandler(action ActionName, pb []byte) (proto.Message, error) {
	handler, ok := r.actionHandlers[action]
	_ = patch_data.TString{
		Key:          "",
		OriginalText: "",
	}
	_ = myprotobuf.AbilityIcon{
		Id:                   0,
		UsageTypeString:      "",
		IconResourceKey:      "",
		IconEffectTypeString: "",
		FrameTypeString:      "",
		ElementTypeString:    "",
		UnitTypeString:       "",
	}

	if !ok {
		errMsg := errors.New("not found route: " + string(action))
		return nil, errMsg
	}
	resultMsg, err := handler(pb)
	if err != nil {
		return nil, errors.Join(errors.New("handle resultMsg"), err)
	}

	return resultMsg, nil
}

func (r *UserRouter) RouterName() RouterName {
	return r.Name
}

func hello(message *proto_define.HelloRequest) (*proto_define.HelloResponse, error) {
	return &proto_define.HelloResponse{
		Echo:     "echo from gs2go: " + message.Msg,
		Sequence: message.Sequence,
	}, nil
}

func NewUserRouter(conn *websocket.Conn) *UserRouter {
	m := make(map[ActionName]func(message []byte) (proto.Message, error), 128)

	parsedHelloMsg := func(message []byte) (proto.Message, error) {
		request := &proto_define.HelloRequest{}
		proto.Unmarshal(message, request)
		response, err := hello(request)
		return response, err
	}

	m[HELLO] = parsedHelloMsg
	userRouter := &UserRouter{conn: conn, actionHandlers: m, Name: RouterName(USER_SERVICE)}

	return userRouter
}
