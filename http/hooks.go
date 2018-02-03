package main

import (
	"github.com/michaelbironneau/analyst/engine"
	"golang.org/x/net/websocket"
	"encoding/json"
	"github.com/michaelbironneau/analyst"
)

type websocketWriter struct {
	ws *websocket.Conn
	msgType string
}

func (w *websocketWriter) Write(p []byte) (n int, err error){
	msg := Message{
		Type: w.msgType,
		Data: json.RawMessage(p),
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}
	return w.ws.Write(b)
}

func redirectOutputHook(ws *websocket.Conn, msgType string) engine.DestinationHook {
	return func(_ string, d engine.Destination) error {
		cd, ok := d.(*engine.ConsoleDestination)
		if !ok {
			return nil
		}
		cd.Writer = &websocketWriter{ws, msgType}
		return nil
	}
}

func outputHooks(ws *websocket.Conn) analyst.RuntimeOptions {
	var opts analyst.RuntimeOptions
	opts.Logger = engine.NewGenericLogger(engine.Trace, &websocketWriter{ws, MsgLog})
	opts.Hooks = []interface{}{redirectOutputHook}
	return opts
}