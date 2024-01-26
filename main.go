package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	ADDRESS = "0.0.0.0:5000"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler(logger))

	logger.Info("starting server", "address", ADDRESS)
	http.ListenAndServe(ADDRESS, mux)
}

func wsHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			logger.Warn("failed to open connection", "ip", r.RemoteAddr)
			return
		}
		defer c.CloseNow()

		// set timeout
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()

		v := HelloMessage{}
		if err := wsjson.Read(ctx, c, &v); err != nil {
			logger.Warn("failed to read socket message", "ip", r.RemoteAddr)
			return
		}

		logger.Info("incoming message", "message", v)

		message := []byte("welcome to our server")
		if err := c.Write(ctx, websocket.MessageText, message); err != nil {
			logger.Warn("failed to write into socket", "error", err.Error())
			return
		}

		// c.Close(websocket.StatusNormalClosure, "")
	}
}

type HelloMessage struct {
	Message string `json:"message"`
}
