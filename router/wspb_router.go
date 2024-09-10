package router

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gs2go/module"
	userModule "gs2go/module/user"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type bufferwarp struct {
	Buffer *bytes.Buffer
	err    error
}

func (bw *bufferwarp) WriteByte(bt byte) {
	if bw.err != nil {
		return
	}

	err := bw.Buffer.WriteByte(bt)
	if err != nil {
		bw.err = errors.Join(errors.New("WsPbRouter"), err)
		return
	}
}

func (bw *bufferwarp) Write(bs []byte) {
	if bw.err != nil {
		return
	}

	_, err := bw.Buffer.Write(bs)
	if err != nil {
		bw.err = errors.Join(errors.New("WsPbRouter"), err)
		return
	}
}

func WsPbRouter(w http.ResponseWriter, r *http.Request, upgrader websocket.Upgrader) error {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Join(errors.New("upgrade"), err)
	}

	defer c.Close()

	routes := setUpWsPb(c)

	for {
		var err error

		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("readMessage: %v", err)
		}

		//log.Printf("recv: %s", message)

		buffer := bytes.NewBuffer(message)
		_, err = buffer.ReadByte()
		if err != nil {
			continue
		}
		serviceBytes := make([]byte, 2)
		_, err = buffer.Read(serviceBytes)
		if err != nil {
			continue
		}
		serviceIndex := int16(binary.LittleEndian.Uint16(serviceBytes))

		actionByte, err := buffer.ReadByte()

		sequence := make([]byte, 4)
		_, err = buffer.Read(sequence)
		if err != nil {
			return errors.Join(errors.New("WsPbRouter"), err)
		}

		available := buffer.Available()
		pb := make([]byte, available)
		_, err = buffer.Read(pb)
		if err != nil {
			return errors.Join(errors.New("WsPbRouter"), err)
		}

		router := routes.Value[module.RouterName{Name: serviceIndex}]
		protoRsp, err := router.WsPbActionHandler(module.ActionName{Name: actionByte}, pb)
		if err != nil {
			log.Printf("wspbActionHandler: %v", err)
			continue
		}
		protoRspBytes, err := proto.Marshal(protoRsp)
		if err != nil {
			log.Printf("WsPbActionHandler: %v", err)

			continue
		}

		// todo response fail message

		newBuffer := bytes.NewBuffer(make([]byte, 1+2+1+4+len(protoRspBytes)))
		newBuffer.Reset()
		bw := bufferwarp{newBuffer, err}
		bw.WriteByte(byte(1))
		bw.Write(serviceBytes)
		bw.WriteByte(actionByte)
		bw.Write(sequence)
		bw.Write(protoRspBytes)
		//marshal, err := protojson.Marshal(protoRsp)
		// _, err = protojson.Marshal(protoRsp)
		// if err != nil {
		// 	return err
		// }

		//log.Printf("rsp: %v", string(marshal))
		c.WriteMessage(websocket.BinaryMessage, newBuffer.Bytes())
	}

	return err
}

func setUpWsPb(conn *websocket.Conn) *module.Routes {
	userRouter := userModule.NewUserRouter(conn)
	routes := module.NewRoutes()
	routes.RegisterRoute(userRouter)
	return routes
}
