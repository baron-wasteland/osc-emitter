package main

import ( 
  "fmt"
  "net/http"
  "math/rand"
  "time"
)


func indexHandler() http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request)  {
    w.Write([]byte("Hey, world."))
  }

  return http.HandlerFunc(fn)
}

func StartSimulation(writes chan *writeOp) {
    tickChan2 := time.NewTicker(time.Millisecond * 20).C
    writeRespChan := make(chan bool)

    go func() {
      mux := http.NewServeMux()

      idx = 
      mux.Handle("/",)
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