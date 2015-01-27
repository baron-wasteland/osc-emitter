package main

import (
  "time"
  "fmt"

  "strconv"
  "bytes"
)

type writeOp struct {
  key int
  val int
  resp chan bool
}

type readOp struct {
  key int
  resp chan int
}



func main() {

  tickChan1 := time.NewTicker(time.Millisecond * 10).C
  reads := make(chan *readOp)
  writes := make(chan *writeOp)
  
  readRespChan := make(chan int)


  go func() {
    sData := make(map[int]int,6)
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
      case <- tickChan1:
        var buffer bytes.Buffer
        for i := 0; i <= 5; i++ {
          read := &readOp{
            key: i,
            resp: readRespChan}
          reads <- read
          val := strconv.Itoa(<-read.resp)
          sensor := strconv.Itoa(i)
          buffer.WriteString("sensor" + sensor + ": " + val + " ")
        }
        fmt.Println(buffer.String())
    }
  }
}
