package main

import (
	"fmt"
	"time"

	"bytes"
	"strconv"
)

type measurement struct {
  sensorType int
  value      int
}

type writeOp struct {
	key  int
	val  measurement
	resp chan bool
}

type readOp struct {
	key  int
	resp chan measurement
}

func main() {
	fmt.Println("started")

	tickChan1 := time.NewTicker(time.Millisecond * 10).C
	reads := make(chan *readOp)
	writes := make(chan *writeOp)

	readRespChan := make(chan measurement)

	go func() {
		sData := make(map[int]measurement, 6)
		for {
			select {
			case read := <-reads:
				read.resp <- sData[read.key]
			case write := <-writes:
				sData[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	go func() {
		StartSimulation(writes)
	}()

	for {
		select {
		case <-tickChan1:
			var buffer bytes.Buffer
			for i := 0; i <= 5; i++ {
				read := &readOp{
					key:  i,
					resp: readRespChan}
				reads <- read
				m          := <- read.resp
				val        := strconv.Itoa(m.value)
				sensorId   := strconv.Itoa(i)
				sensorType := strconv.Itoa(m.sensorType)
				buffer.WriteString("sensor" + sensorId + "[" + sensorType + "]: " + val + " ")
			}
			t := time.Now()
			fmt.Printf("\r" + t.Format(time.StampMilli) + ": " + buffer.String())
		}
	}
}
