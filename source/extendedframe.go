package source

import (
	"go.einride.tech/can"
)

type ExtendedFrame struct {
	ID   uint32
	Data []byte
}

func NewExtendedFrame(frame *can.Frame) *ExtendedFrame {
	newFrame := ExtendedFrame{ID: frame.ID, Data: make([]byte, 0)}
	newFrame.Data = append(newFrame.Data, frame.Data[0:]...)
	return &newFrame
}

func packLittleEndian(bytes []byte, startByte int) uint64 {
	var packed uint64
	pad := 8 - (len(bytes) - startByte)
	for i := 0; i < pad; i++ {
		bytes = append(bytes, 0)
	}
	packed |= uint64(bytes[startByte+0]) << (0 * 8)
	packed |= uint64(bytes[startByte+1]) << (1 * 8)
	packed |= uint64(bytes[startByte+2]) << (2 * 8)
	packed |= uint64(bytes[startByte+3]) << (3 * 8)
	packed |= uint64(bytes[startByte+4]) << (4 * 8)
	packed |= uint64(bytes[startByte+5]) << (5 * 8)
	packed |= uint64(bytes[startByte+6]) << (6 * 8)
	packed |= uint64(bytes[startByte+7]) << (7 * 8)
	return packed
}

func asSigned(unsigned uint64, bits int) int64 {
	switch bits {
	case 8:
		return int64(int8(uint8(unsigned)))
	case 16:
		return int64(int16(uint16(unsigned)))
	case 32:
		return int64(int32(uint32(unsigned)))
	case 64:
		return int64(unsigned)
	default:
		// calculate bit mask for sign bit
		signBitMask := uint64(1 << (bits - 1))
		// check if sign bit is set
		isNegative := unsigned&signBitMask > 0
		if !isNegative {
			// sign bit not set means we can reinterpret the value as-is
			return int64(unsigned)
		}
		// calculate bit mask for extracting value bits (all bits except the sign bit)
		valueBitMask := signBitMask - 1
		// calculate two's complement of the value bits
		value := ((^unsigned) & valueBitMask) + 1
		// result is the negative value of the two's complement
		return -1 * int64(value)
	}
}

func (frame *ExtendedFrame) UnsignedBitsLittleEndian(start int, length int) uint64 {
	startByte := 0
	bitOffset := 0

	if start+length > 64 {
		startByte = start / 8
		bitOffset = startByte * 8
	}

	bytes := frame.Data
	// pack bits into one continuous value
	packed := packLittleEndian(bytes, startByte)
	// lsb index in the packed value is the start bit
	lsbIndex := start - bitOffset
	// shift away lower bits
	shifted := packed >> lsbIndex
	// mask away higher bits
	masked := shifted & ((1 << length) - 1)
	// done
	return masked
}

func (frame *ExtendedFrame) SignedBitsLittleEndian(start, length int) int64 {
	unsigned := frame.UnsignedBitsLittleEndian(start, length)
	return asSigned(unsigned, length)
}
