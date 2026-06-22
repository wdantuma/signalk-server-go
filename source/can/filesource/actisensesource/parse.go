package actisensesource

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"go.einride.tech/can"
)

func Parse(s string) ([]can.Frame, error) {
	frames := make([]can.Frame, 0)
	parts := strings.Split(s, ",")
	if len(parts) < 14 {
		return frames, errors.New("Invalid actisense string")
	}
	var fid uint32
	prio, err := strconv.ParseUint(parts[1], 10, 3)
	if err != nil {
		return frames, err
	}
	fid = uint32(prio) << 29
	pgn, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return frames, err
	}
	fid = fid | (uint32(pgn) << 8)
	src, err := strconv.ParseUint(parts[3], 10, 8)
	if err != nil {
		return frames, err
	}
	fid = fid | uint32(src)
	ID := fid
	IsExtended := true
	len, err := strconv.ParseUint(parts[5], 10, 8)
	if err != nil {
		return frames, err
	}
	bytes, err := hex.DecodeString(strings.Join(parts[6:], ""))
	if err != nil {
		return frames, err
	}

	if len == 8 {
		frame := can.Frame{}
		frame.ID = ID
		frame.IsExtended = IsExtended
		frame.Data = can.Data(bytes)
		frames = append(frames, frame)
	} else {
		var seq uint8 = 0
		var n uint8 = 0x40
		var index uint64 = 0
		for len > index {
			newBytes := make([]uint8, 8)
			for i := 0; i < 8; i++ {
				newBytes[i] = 255
			}
			frame := can.Frame{}
			frame.ID = ID
			frame.IsExtended = IsExtended
			newBytes[0] = n | seq
			if index == 0 {
				newBytes[1] = uint8(len)
				copy(newBytes[2:], bytes[index:index+6])
				index += 6
			} else {
				if len-index > 7 {
					copy(newBytes[1:], bytes[index:index+7])
				} else {
					copy(newBytes[1:], bytes[index:])
				}
				index += 7
			}
			frame.Data = can.Data(newBytes)
			frames = append(frames, frame)
			seq++
		}

	}

	return frames, nil
}
