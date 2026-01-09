package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
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
	Type    string `json:"type"`
	Index   int    `json:"index"`
	Degrees int    `json:"degrees"`
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load .env file")
	}

	serialPort := os.Getenv("SERIAL_PORT")
	baudRateStr := os.Getenv("BAUD_RATE")
	baudRate, err := strconv.Atoi(baudRateStr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse BAUD_RATE")
	}

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

	flag.Parse()

	http.HandleFunc("/", server)

	log.Info().Msgf("Starting server on %s", *addr)

	log.Fatal().Err(http.ListenAndServe(*addr, nil))
}

func handleWebSocketMessage(msg WebsocketMessage) {
	switch msg.Type {
	case "setArmAngle":
		// Map index: 0->1, 1->2
		servoNum := msg.Index + 1

		// Map degrees: 0->90, 90->180, -90->0
		mappedDegrees := msg.Degrees + 90

		log.Info().Msgf("Rotating servo %d to %d degrees (original: %d)", servoNum, mappedDegrees, msg.Degrees)

		// Send as text: "servo1 90\n"
		command := fmt.Sprintf("servo%d %d\n", servoNum, mappedDegrees)
		_, err := port.Write([]byte(command))
		if err != nil {
			log.Error().Err(err).Msg("failed to write to serial port")
			return
		}
	default:
		log.Warn().Msgf("Unknown message type: %s", msg.Type)
	}
}
