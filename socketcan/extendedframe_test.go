package socketcan

import (
	"testing"
)

func TestUnsignedBitsLittleEndian0(t *testing.T) {
	frame := ExtendedFrame{}
	//frame.Data = []byte{0x01, 0x90, 0xB3, 0x99, 0x0E, 0xB5, 0x2C, 0x22, 0x03, 0x1E, 0x8D, 0xC6, 0x1F, 0x0B, 0x9D, 0x1B, 0x95, 0x00, 0x41, 0x00, 0x09, 0xFF, 0xFF, 0xFF, 0x7F, 0xC5, 0xF8, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	frame.Data = []byte{0xB5, 0x2C, 0x22, 0x03}

	val := frame.UnsignedBitsLittleEndian(0, 32)
	if val != 52571317 {
		t.Error()
	}

}

func TestUnsignedBitsLittleEndian1(t *testing.T) {
	frame := ExtendedFrame{}
	frame.Data = []byte{0x5A, 0x96, 0x11, 0x01, 0x80}

	val := frame.UnsignedBitsLittleEndian(1, 32)
	if val != 52571317 {
		t.Error()
	}

}

func TestUnsignedBitsLittleEndian2(t *testing.T) {
	frame := ExtendedFrame{}
	frame.Data = []byte{0x6D, 0x4B, 0x08, 0x80, 0xC0}

	val := frame.UnsignedBitsLittleEndian(2, 32)
	if val != 52571317 {
		t.Error()
	}
}

func TestUnsignedBitsLittleEndian3(t *testing.T) {
	frame := ExtendedFrame{}
	frame.Data = []byte{0x16, 0xA5, 0x84, 0x40, 0x60}

	val := frame.UnsignedBitsLittleEndian(3, 32)
	if val != 52571317 {
		t.Error()
	}
}

func TestUnsignedBitsLittleEndian7(t *testing.T) {
	frame := ExtendedFrame{}
	frame.Data = []byte{0x01, 0x6A, 0x58, 0x44, 0x06}

	val := frame.UnsignedBitsLittleEndian(7, 32)
	if val != 52571317 {
		t.Error()
	}
}

func TestUnsignedBitsLittleEndian40(t *testing.T) {
	frame := ExtendedFrame{}
	frame.Data = []byte{0x01, 0x90, 0xB3, 0x99, 0x0E, 0xB5, 0x2C, 0x22, 0x03, 0x1E, 0x8D, 0xC6, 0x1F, 0x0B, 0x9D, 0x1B, 0x95, 0x00, 0x41, 0x00, 0x09, 0xFF, 0xFF, 0xFF, 0x7F, 0xC5, 0xF8, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	val := frame.UnsignedBitsLittleEndian(40, 32)
	if val != 52571317 {
		t.Error()
	}
}

func TestUnsignedBitsLittleEndian33(t *testing.T) {
	frame := ExtendedFrame{}
	frame.Data = []byte{0, 18}

	val := frame.UnsignedBitsLittleEndian(3, 3)
	if val != 3 {
		t.Error()
	}
}
