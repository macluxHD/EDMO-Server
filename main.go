package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.bug.st/serial"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type SetArmAngleData struct {
	Index   int `json:"index"`
	Degrees int `json:"degrees"`
}

type OscillatorUpdateData struct {
	Index      int     `json:"index"`
	Frequency  float32 `json:"frequency"`
	Amplitude  float32 `json:"amplitude"`
	Offset     float32 `json:"offset"`
	PhaseShift float32 `json:"phaseShift"`
}

func server(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade connection")
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("failed to read message")
			break
		}
		log.Debug().Msgf("received: %s", message)

		var msg WebsocketMessage
		err = json.Unmarshal(message, &msg)

		handleWebSocketMessage(msg)
	}
}

var port serial.Port
var mock bool

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("failed to load .env file")
	}

	serialPort := os.Getenv("SERIAL_PORT")
	baudRateStr := os.Getenv("BAUD_RATE")
	mock = os.Getenv("MOCK") == "true"
	baudRate, err := strconv.Atoi(baudRateStr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse BAUD_RATE")
	}

	if !mock {
		mode := &serial.Mode{
			BaudRate: baudRate,
			Parity:   serial.NoParity,
			DataBits: 8,
			StopBits: serial.OneStopBit,
		}

		port, err = serial.Open(serialPort, mode)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to open port %s", serialPort)
		}

		// send startup commands to EDMO robot and start listening for incoming packets
		startEDMO()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info().Msgf("received signal: %s, shutting down…", sig)
		if !mock {
			shutdownEDMO()
		}
		os.Exit(0)
	}()

	flag.Parse()
	http.HandleFunc("/", server)
	log.Info().Msgf("Starting server on %s", *addr)
	log.Fatal().Err(http.ListenAndServe(*addr, nil))
}

func handleWebSocketMessage(msg WebsocketMessage) {
	switch msg.Type {
	case "setArmAngle":
		var data SetArmAngleData
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal setArmAngle data")
			return
		}

		servoNum := data.Index

		// Map degrees: 0->90, 90->180, -90->0
		mappedDegrees := data.Degrees + 90

		log.Info().Msgf("Rotating servo %d to %d degrees (original: %d)", servoNum, mappedDegrees, data.Degrees)

		if !mock {
			setAngle(byte(servoNum), float32(mappedDegrees))
		}
	case "setOscillator":
		var data OscillatorUpdateData
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal setOscillator data")
			return
		}

		log.Info().Msgf("Updating oscillator %d: freq=%.2f, amp=%.2f, offset=%.2f, phaseShift=%.2f",
			data.Index, data.Frequency, data.Amplitude, data.Offset, data.PhaseShift)

		if !mock {
			oscillatorUpdate(byte(data.Index), data.Frequency, data.Amplitude, data.Offset, data.PhaseShift)
		}
	default:
		log.Warn().Msgf("Unknown message type: %s", msg.Type)
	}
}
