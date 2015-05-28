# osc-emitter

Takes various sensor inputs and outputs a stream of OSC messages.

Install Dependencies
======
```
brew install portmidi
go get github.com/gorilla/websocket
go get github.com/hypebeast/go-osc/osc
```

Overview
===
- main loop manages slice of Instruments, with reads and writes flowing through channels.
- separate loops will receive input form sensors (or the web simulation) and update Instruments.
- each Instrument will send a constant stream of OSC messages to N receivers.


Run it
===
```
go run main.go serve.go instrument.go -c config/test.json
```

- Load the web manager ui to view instrument state: http://localhost:8080/manager.html
- Load the web simulation ui to change state: http://localhost:8080/simulation.html (works on mobile, just change localhost to the proper ip)
- start things that listen for OSC and map stuff

Create other config files for other situations.

TODO
===
- instrument interface to create other types of instruments
- refactor current instrument implemenation to utilize a general interface
- create other types of instruments (linear/rotary encoder, always on/off, others)
- web ui using config files with config endpoint for js
- how to tweak config settings during performance?
- better web sim physics