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

func escape(data []byte) []byte {
	out := bytes.Buffer{}

	for _, b := range data {
		// Escape header/footer bytes ('E', 'D', 'M', 'O') and the escape character itself ('\\')
		if b == 'E' || b == 'D' || b == 'M' || b == 'O' || b == '\\' {
			out.WriteByte('\\')
		}
		out.WriteByte(b)
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
			time := binary.LittleEndian.Uint32(payload[0:4])
			frequency_0 := math.Float32frombits(binary.LittleEndian.Uint32(payload[4:8]))
			amplitude_0 := math.Float32frombits(binary.LittleEndian.Uint32(payload[8:12]))
			offset_0 := math.Float32frombits(binary.LittleEndian.Uint32(payload[12:16]))
			phaseShift_0 := math.Float32frombits(binary.LittleEndian.Uint32(payload[16:20]))
			phase_0 := math.Float32frombits(binary.LittleEndian.Uint32(payload[20:24]))
			frequency_1 := math.Float32frombits(binary.LittleEndian.Uint32(payload[24:28]))
			amplitude_1 := math.Float32frombits(binary.LittleEndian.Uint32(payload[28:32]))
			offset_1 := math.Float32frombits(binary.LittleEndian.Uint32(payload[32:36]))
			phaseShift_1 := math.Float32frombits(binary.LittleEndian.Uint32(payload[36:40]))
			phase_1 := math.Float32frombits(binary.LittleEndian.Uint32(payload[40:44]))
			frequency_2 := math.Float32frombits(binary.LittleEndian.Uint32(payload[44:48]))
			amplitude_2 := math.Float32frombits(binary.LittleEndian.Uint32(payload[48:52]))
			offset_2 := math.Float32frombits(binary.LittleEndian.Uint32(payload[52:56]))
			phaseShift_2 := math.Float32frombits(binary.LittleEndian.Uint32(payload[56:60]))
			phase_2 := math.Float32frombits(binary.LittleEndian.Uint32(payload[60:64]))
			frequency_3 := math.Float32frombits(binary.LittleEndian.Uint32(payload[64:68]))
			amplitude_3 := math.Float32frombits(binary.LittleEndian.Uint32(payload[68:72]))
			offset_3 := math.Float32frombits(binary.LittleEndian.Uint32(payload[72:76]))
			phaseShift_3 := math.Float32frombits(binary.LittleEndian.Uint32(payload[76:80]))
			phase_3 := math.Float32frombits(binary.LittleEndian.Uint32(payload[80:84]))
			a := math.Float32frombits(binary.LittleEndian.Uint32(payload[84:88]))
			gyro_x := math.Float32frombits(binary.LittleEndian.Uint32(payload[92:96]))
			gyro_y := math.Float32frombits(binary.LittleEndian.Uint32(payload[96:100]))
			gyro_z := math.Float32frombits(binary.LittleEndian.Uint32(payload[100:104]))
			f := math.Float32frombits(binary.LittleEndian.Uint32(payload[104:108]))
			accel_x := math.Float32frombits(binary.LittleEndian.Uint32(payload[112:116]))
			accel_y := math.Float32frombits(binary.LittleEndian.Uint32(payload[116:120]))
			accel_z := math.Float32frombits(binary.LittleEndian.Uint32(payload[120:124]))
			k := math.Float32frombits(binary.LittleEndian.Uint32(payload[124:128]))
			mag_x := math.Float32frombits(binary.LittleEndian.Uint32(payload[132:136]))
			mag_y := math.Float32frombits(binary.LittleEndian.Uint32(payload[136:140]))
			mag_z := math.Float32frombits(binary.LittleEndian.Uint32(payload[140:144]))
			p := math.Float32frombits(binary.LittleEndian.Uint32(payload[144:148]))
			grav_x := math.Float32frombits(binary.LittleEndian.Uint32(payload[152:156]))
			grav_y := math.Float32frombits(binary.LittleEndian.Uint32(payload[156:160]))
			grav_z := math.Float32frombits(binary.LittleEndian.Uint32(payload[160:164]))
			u := math.Float32frombits(binary.LittleEndian.Uint32(payload[164:168]))
			rot_x := math.Float32frombits(binary.LittleEndian.Uint32(payload[172:176]))
			rot_y := math.Float32frombits(binary.LittleEndian.Uint32(payload[176:180]))
			rot_z := math.Float32frombits(binary.LittleEndian.Uint32(payload[180:184]))
			rot_w := math.Float32frombits(binary.LittleEndian.Uint32(payload[184:188]))
			log.Debug().Msgf("Cmd69 parsed: Time=%d, Frequency(0)=%.2f, Amplitude(0)=%.2f, Offset(0)=%.2f, PhaseShift(0)=%.2f, Phase(0)=%.2f, Frequency(1)=%.2f, Amplitude(1)=%.2f, Offset(1)=%.2f, PhaseShift(1)=%.2f, Phase(1)=%.2f, Frequency(2)=%.2f, Amplitude(2)=%.2f, Offset(2)=%.2f, PhaseShift(2)=%.2f, Phase(2)=%.2f, Frequency(3)=%.2f, Amplitude(3)=%.2f, Offset(3)=%.2f, PhaseShift(3)=%.2f, Phase(3)=%.2f, a=%.2f, Gyro=(%.2f, %.2f, %.2f), f=%.2f, Accel=(%.2f, %.2f, %.2f), k=%.2f, Mag=(%.2f, %.2f, %.2f), p=%.2f, Grav=(%.2f, %.2f, %.2f), u=%.2f, Rot=(%.2f, %.2f, %.2f, %.2f)",
				time,
				frequency_0, amplitude_0, offset_0, phaseShift_0, phase_0,
				frequency_1, amplitude_1, offset_1, phaseShift_1, phase_1,
				frequency_2, amplitude_2, offset_2, phaseShift_2, phase_2,
				frequency_3, amplitude_3, offset_3, phaseShift_3, phase_3,
				a, gyro_x, gyro_y, gyro_z,
				f, accel_x, accel_y, accel_z,
				k, mag_x, mag_y, mag_z,
				p, grav_x, grav_y, grav_z,
				u, rot_x, rot_y, rot_z, rot_w,
			)
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

	// Escape the command and payload to avoid header/footer bytes corrupting the packet
	escapedPayload := escape(append([]byte{cmd}, payload...))
	buf.Write(escapedPayload)

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
