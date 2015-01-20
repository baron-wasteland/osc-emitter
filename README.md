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