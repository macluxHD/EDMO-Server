package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func readPackets() {
	var buf bytes.Buffer
	temp := make([]byte, 1)

	for {
		n, err := port.Read(temp)
		if err != nil {
			log.Error().Err(err).Msg("serial read error")
			return
		}
		if n == 0 {
			continue
		}

		b := temp[0]
		buf.WriteByte(b)

		data := buf.Bytes()

		if len(data) >= 2 && string(data[:2]) == "ED" {
			if len(data) >= 4 && string(data[len(data)-2:]) == "MO" {
				packet := make([]byte, len(data)-4)
				copy(packet, data[2:len(data)-2]) // strip "ED" and "MO"
				buf.Reset()

				go handlePacket(packet)
			}
		} else if buf.Len() > 1024 {
			buf.Reset()
		}
	}
}

func unescape(data []byte) []byte {
	out := bytes.Buffer{}
	skip := false

	for i := 0; i < len(data); i++ {
		if skip {
			skip = false
			continue
		}

		if data[i] == '\\' {
			if i+1 < len(data) {
				out.WriteByte(data[i+1])
				skip = true
			}
		} else {
			out.WriteByte(data[i])
		}
	}
	return out.Bytes()
}

func handlePacket(packet []byte) {
	// Unescape (EDMO escapes header/footer bytes in the payload)
	raw := unescape(packet)

	if len(raw) < 1 {
		return
	}

	cmd := raw[0]
	payload := raw[1:]

	switch cmd {
	case 69:
		log.Debug().Msgf("Received command 69 with %d payload bytes", len(payload))

		if len(payload) >= 4+4+4 {
			a := binary.LittleEndian.Uint32(payload[0:4])
			b := math.Float32frombits(binary.LittleEndian.Uint32(payload[4:8]))
			c := math.Float32frombits(binary.LittleEndian.Uint32(payload[8:12]))
			d := math.Float32frombits(binary.LittleEndian.Uint32(payload[12:16]))
			log.Debug().Msgf("Cmd69 parsed: Time=%d, Frequency=%.2f, Amplitude=%.2f, Offset=%.2f", a, b, c, d)
		} else {
			log.Debug().Msgf("Cmd69 raw payload: %x", payload)
		}

	default:
		log.Warn().Msgf("Unknown command %d (payload %d bytes)", cmd, len(payload))
	}
}

func shutdownEDMO() {
	log.Info().Msg("ending session…")
	if err := endSession(); err != nil {
		log.Error().Err(err).Msg("failed to send end session")
	}
	if port != nil {
		port.Close()
	}
}

func writeEDMOPacket(cmd byte, payload []byte) error {
	buf := bytes.NewBuffer(nil)

	buf.Write([]byte("ED"))

	buf.WriteByte(cmd)

	if payload != nil {
		buf.Write(payload)
	}

	buf.Write([]byte("MO"))

	_, err := port.Write(buf.Bytes())
	if err != nil {
		log.Error().Err(err).Msgf("failed write packet cmd %d", cmd)
	}
	return err
}

func identificationCommand(uuidStr string) error {
	uuidBytes, err := parseUUID(uuidStr)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	return writeEDMOPacket(0, uuidBytes)
}

func parseUUID(s string) ([]byte, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return uid[:], nil
}

func startSession(refTime uint32) error {
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload, refTime)
	return writeEDMOPacket(1, payload)
}

func endSession() error {
	log.Debug().Msg("sending end session command")
	return writeEDMOPacket(6, nil)
}

func oscillatorUpdate(idx byte, freq, amp, offset, phaseShift float32) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(idx)
	binary.Write(buf, binary.LittleEndian, freq)
	binary.Write(buf, binary.LittleEndian, amp)
	binary.Write(buf, binary.LittleEndian, offset)
	binary.Write(buf, binary.LittleEndian, phaseShift)

	return writeEDMOPacket(3, buf.Bytes())
}

func setAngle(idx byte, angle float32) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(idx)
	binary.Write(buf, binary.LittleEndian, angle)

	return writeEDMOPacket(7, buf.Bytes())
}

func startEDMO() {
	id := uuid.New().String()
	if err := identificationCommand(id); err != nil {
		log.Error().Err(err).Msg("failed to send identification")
	}

	refTime := uint32(time.Now().Unix())
	if err := startSession(refTime); err != nil {
		log.Error().Err(err).Msg("failed to start session")
	}

	go readPackets()
}
