# osc-emitter

Takes various sensor inputs and outputs a stream of OSC messages.


Overview
=====

The main loop will manage an object with the value for each sensor, with reads and writes passing through a channel that can be handed to other functions, a la https://gobyexample.com/stateful-goroutines

Code that needs to write to the sensor object can do so at will.

Another routine will check the values of that object every Xms and dispatch OSC messages.

We will need to support the following types of events:

 - Constant stream of OSC messages, repeating the last value if unchanged.
 - Note on/off events
   - timed duration
   - data driven duration (like a threshold)

Web Simulation
====

The web simulation serves up a basic html page with a script that upgrades the connection to a websocket, using http://www.gorillatoolkit.org/pkg/websocket

To use it:

    go get github.com/gorilla/websocket
    go run main.go simulation.go

Open your browser and navigate to localhost:8080/simulation.html
Hit the 'send' button, and every 20 ms a random int will be written to a sensor object, and sensor values will be written to the terminal every 10ms.

Reload the page to stop sending.