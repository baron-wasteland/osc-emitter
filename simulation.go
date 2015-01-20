package main

import ( 
  "fmt"
  "net/http"
  "math/rand"
  "time"
)


func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func StartSimulation(writes chan *writeOp) {
    tickChan2 := time.NewTicker(time.Millisecond * 20).C
    writeRespChan := make(chan bool)

    go func() {
      http.HandleFunc("/", handler)
      http.ListenAndServe(":8080", nil)      
      }()


    for {
    select {
      case <- tickChan2:
        write := &writeOp{
          key:  rand.Intn(5),
          val:  rand.Intn(100),
          resp: writeRespChan}
        writes <- write
        <-write.resp
    }
   }

}