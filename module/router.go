package module

import "google.golang.org/protobuf/proto"

const (
	USER_SERVICE  int16 = 1
	PGSQL_SERVICE int16 = 2
)

type ActionName struct {
	Name byte
}

type Router interface {
	WsPbActionHandler(action ActionName, pb []byte) (proto.Message, error)
	RouterName() RouterName
}
