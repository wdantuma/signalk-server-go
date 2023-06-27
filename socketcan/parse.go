package socketcan

import (
	"errors"
	"strconv"
	"strings"

	"go.einride.tech/can"
)

func Parse(s string) (can.Frame, error) {
	frame := can.Frame{}
	parts := strings.Split(s, "#")
	if len(parts) != 2 {
		return frame, errors.New("Invalid candump log string")
	}

	fid, err := strconv.ParseUint(parts[0], 16, 32)
	if err != nil {
		return frame, err
	}
	frame.ID = uint32(fid)
	frame.IsExtended = true

	data, err := strconv.ParseUint(parts[1], 16, 64)
	if err != nil {
		return frame, err
	}
	frame.Data.UnpackBigEndian(data)
	frame.Length = 8

	return frame, nil
}
