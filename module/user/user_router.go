package user

import (
	"encoding/json"
	"errors"
	"gs2go/module"

	"github.com/gorilla/websocket"
)

type PbMessage struct {
	field1 string
	field2 int
}
type UserRouter struct {
	conn           *websocket.Conn
	actionHandlers map[module.ActionName]func(PbMessage) (PbMessage, error) // key: router name, value: output
	Name           module.RouterName
}

func (r *UserRouter) WsPbActionHandler(action module.ActionName) error {
	handler, ok := r.actionHandlers[action]
	if !ok {
		errMsg := errors.New("not found route: ")
		return errMsg
	}
	resultMsg, err := handler(PbMessage{"field1", 2})
	if err != nil {
		return errors.Join(errors.New("handle resultMsg"), err)
	}
	marshal, err := json.Marshal(resultMsg)
	if err != nil {
		return errors.Join(errors.New("marshal resultMsg"), err)
	}
	r.conn.WriteMessage(websocket.TextMessage, marshal)
	return nil
}

func (r *UserRouter) RouterName() module.RouterName {
	return r.Name
}

func NewUserRouter(conn *websocket.Conn) *UserRouter {
	m := make(map[module.ActionName]func(PbMessage) (PbMessage, error), 128)
	m[module.ActionName{"SignUp"}] = SignUp
	userRouter := &UserRouter{conn: conn, actionHandlers: m, Name: module.RouterName{"UserRouter"}}

	return userRouter
}

func SignUp(message PbMessage) (PbMessage, error) {
	return message, nil
}
