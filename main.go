package main

import (
	"encoding/json"

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

type SensorConfig struct {
	Types []SensorType `json:"types"`
}

type InstrumentConfig struct {
	Id         int          `json:"id"`
	Threshold  int          `json:"threshold"`
	SensorType int          `json:"sensorType"`
	Controls   []OscControl `json:"controls"`
	Notes      []Note       `json:"notes"`
}

type Instruments struct {
	Instruments []InstrumentConfig `json:"instruments"`
}

type OscConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type OscConfigs struct {
	OscConfig []OscConfig `json:"oscConfig"`
}

func loadConfig(oscFile string, sensorFile string, instrumentsFile string) ([]*Instrument, error) {
	// setup osc hosts
	file, _ := os.Open(oscFile)
	decoder := json.NewDecoder(file)
	oscConfig := OscConfigs{}
	err := decoder.Decode(&oscConfig)
	if err != nil {
		return nil, err
	}

	// setup sensor defs
	file, _ = os.Open(sensorFile)
	decoder = json.NewDecoder(file)
	sensorConfig := SensorConfig{}
	err = decoder.Decode(&sensorConfig)
	if err != nil {
		return nil, err
	}

	// setup instruments
	file, _ = os.Open(instrumentsFile)
	decoder = json.NewDecoder(file)
	instrumentConfig := Instruments{}
	err = decoder.Decode(&instrumentConfig)
	if err != nil {
		return nil, err
	}

	instruments := make([]*Instrument, 0)
	for _, instrument := range instrumentConfig.Instruments {
		ins := CreateInstrument(
			instrument.Id,
			instrument.Notes,
			instrument.Controls,
			sensorConfig.Types[instrument.SensorType],
			oscConfig.OscConfig)

		instruments = append(instruments, ins)
	}

	return instruments, nil
}

func main() {

	// TODO: setup some file watches for auto config reload
	instruments, err := loadConfig("config/osc.json", "config/sensors.json", "config/instruments.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
