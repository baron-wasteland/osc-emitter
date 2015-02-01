package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	// "math/rand"
	"net/http"
	// "strconv"
	// "time"
)

const (
	SENSOR_CONTINUOUS = 0
	SENSOR_IMPULSE = 1
)

type Message struct {
	Value int
	Id    int
	SensorType int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var m Message

func readWrapper(writes chan *writeOp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeRespChan := make(chan bool)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			//log.Println(err)
			return
		}

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// s := string(p[:len(p)])
			// fmt.Printf(s)
			// fmt.Printf("\n")

			_ = json.Unmarshal(p, &m)
			fmt.Printf("%d: %d, %d\n", m.Id, m.Value, m.SensorType)

			write := &writeOp{
				key:  m.Id,
				val:  m.Value,
				resp: writeRespChan}
			writes <- write
			<-write.resp

		}
	}
}

func StartSimulation(writes chan *writeOp) {
	go func() {
		fmt.Println("Listening")
		http.HandleFunc("/echo", readWrapper(writes))
		http.Handle("/", http.FileServer(http.Dir(".")))
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic("Error: " + err.Error())
		}

	}()
}
