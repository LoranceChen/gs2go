package router

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"gs2go/module"
	userModule "gs2go/module/user"
	"log"
	"net/http"
)

func WsPbRouter(w http.ResponseWriter, r *http.Request, upgrader websocket.Upgrader) error {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Join(errors.New("upgrade"), err)
	}

	defer c.Close()

	routes := setUpWsPb(c)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("readMessage: %v", err)
		}

		log.Printf("recv: %s", message)

		jsonMessage := &JsonMessage{}
		err = json.Unmarshal(message, jsonMessage)
		if err != nil {
			log.Printf("jsonErr: %v", err)
			break
		}
		router := routes.Value[module.RouterName{Name: jsonMessage.Service}]
		err = router.WsPbActionHandler(module.ActionName{Name: jsonMessage.Action})
		if err != nil {
			log.Printf("WsPbActionHandler: %v", err)
			break
		}

	}

	return err
}

func setUpWsPb(conn *websocket.Conn) *module.Routes {
	userRouter := userModule.NewUserRouter(conn)
	routes := module.NewRoutes()
	routes.RegisterRoute(userRouter)
	return routes
}
