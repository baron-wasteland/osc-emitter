package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	NUM_SENSORS       = 6
	SENSOR_CONTINUOUS = 0
	SENSOR_IMPULSE    = 1
	OSC_SEND_FREQ_MS  = 100
)

type readInstrument struct {
	key  int
	resp chan *Instrument
}

type measurement struct {
	sensorId   int
	sensorType int
	value      int
}

type SensorType struct {
	MinVal int    `json:"minVal"`
	MaxVal int    `json:"maxVal"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
}

type InstrumentConfig struct {
	Id         int          `json:"id"`
	Threshold  int          `json:"threshold"`
	SensorType int          `json:"sensorType"`
	Controls   []OscControl `json:"controls"`
	Notes      []Note       `json:"notes"`
}

type OscConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Config struct {
	Sensors     []SensorType       `json:"sensors"`
	OscConfig   []OscConfig        `json:"oscConfig"`
	Instruments []InstrumentConfig `json:"instruments"`
}

func main() {

	// parse some flags
	configFile := flag.String("c", "", "path to config file")
	flag.Parse()
	if *configFile == "" {
		fmt.Println("You must pass a config file, run with -h to see how")
		os.Exit(1)
	}

	// read config
	file, _ := os.Open(*configFile)
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// make instruments from config
	instruments := make([]*Instrument, 0)
	for _, instrument := range config.Instruments {
		ins := CreateInstrument(
			instrument.Id,
			instrument.Notes,
			instrument.Controls,
			config.Sensors[instrument.SensorType],
			config.OscConfig)

		instruments = append(instruments, ins)
	}

	fmt.Println("started main")

	// channel to receive sensor updates on
	updates := make(chan measurement, 5000)

	// channel to request instrument details on
	reads := make(chan *readInstrument)

	// setup signal handling
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// TODO: add some flags to enable/disable simulation or something
	go func() {
		StartServer(updates, reads)
	}()

	for {
		select {
		// update instrument from sensor/simulation
		case m := <-updates:
			instruments[m.sensorId].update(m)
		// manager page state reads
		case read := <-reads:
			// fmt.Println(instruments[read.key])
			read.resp <- instruments[read.key]
		case <-signals:
			fmt.Println("signal received, cleaning up")
			// placeholder for when we have actual cleanup work
			fmt.Println("cleanup complete, exiting")
			os.Exit(0)
		}
	}

}
