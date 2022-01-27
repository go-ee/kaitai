package kaitai

import (
	"encoding/binary"
	"fmt"
	"math"
	"regexp"
)

type Native struct {
	Type     string
	Length   uint8
	EndianBe *bool
}

func (o *Native) BuildReader(attr *Attr, spec *Spec) (ret AttrReader, err error) {
	if o.EndianBe == nil {
		o.EndianBe = spec.Meta.EndianBe
	}
	var endianConverter EndianReader
	if *o.EndianBe {
		endianConverter = BigEndianConverter
	} else {
		endianConverter = LittleEndianConverter
	}

	var fix ReadFix
	var dynamic ReadDynamic

	switch o.Type {
	case "str":
		fix = ReadFixAttr(attr, ToString)
		dynamic = ReadDynamicAttr(attr, ToString)
	case "strz":
		fix = ReadFixAttr(attr, ToString)
		dynamic = ReadDynamicAttr(attr, ToString)
	case "b":
		fix = ReadB(endianConverter, o.Length)
	case "u":
		fix = ReadU(endianConverter, o.Length)
	case "s":
		fix = ReadS(endianConverter, o.Length)
	case "f":
		fix = ReadF(endianConverter, o.Length)
	default:
		err = fmt.Errorf("not supported Native(%v,%v)", o.Type, o.Length)
	}

	if fix != nil {
		ret = &NativeReaderFix{
			AttrReaderBase: &AttrReaderBase{attr, o},
			fix:            fix,
		}
	} else if dynamic != nil {
		ret = &NativeReaderDynamic{
			AttrReaderBase: &AttrReaderBase{attr, o},
			dynamic:        dynamic}
	}
	return
}

type NativeReaderFix struct {
	*AttrReaderBase
	fix ReadFix
}

func (o *NativeReaderFix) ReadTo(fillItem *Item, reader Reader) (err error) {
	if value, currentErr := o.fix(reader); currentErr == nil {
		fillItem.Value = value
	} else {
		err = currentErr
	}
	return
}

type NativeReaderDynamic struct {
	*AttrReaderBase
	dynamic ReadDynamic
}

func (o *NativeReaderDynamic) ReadTo(fillItem *Item, reader Reader) (err error) {
	if value, currentErr := o.dynamic(reader, fillItem.Parent); currentErr == nil {
		fillItem.Value = value
	} else {
		err = currentErr
	}
	return
}

var buildInRegExp *regexp.Regexp

func init() {
	buildInRegExp = regexp.MustCompile(`([bfsu])([1-8])(be|le|)`)
}

type EndianReader interface {
	Uint(data []byte) uint
	Uint8(data []byte) uint8
	Uint16(data []byte) uint16
	Uint32(data []byte) uint32
	Uint64(data []byte) uint64
	Float32fromBits(data []byte) float32
	Float64fromBits(data []byte) float64
	ReadBitsInt(reader Reader, n uint8) (ret uint64, err error)
}

var BigEndianConverter *bigEndianConverter
var LittleEndianConverter *littleEndianConverter

type bigEndianConverter struct {
}

func (o *bigEndianConverter) Uint(data []byte) uint {
	return uint(data[0])
}

func (o *bigEndianConverter) Uint8(data []byte) uint8 {
	return data[0]
}

func (o *bigEndianConverter) Uint16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}

func (o *bigEndianConverter) Uint32(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}
func (o *bigEndianConverter) Uint64(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}
func (o *bigEndianConverter) Float32fromBits(data []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(data))
}
func (o *bigEndianConverter) Float64fromBits(data []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(data))
}

func (o *bigEndianConverter) ReadBitsInt(reader Reader, n uint8) (ret uint64, err error) {
	bitsNeeded := int(n) - int(reader.bitsLeft)
	if bitsNeeded > 0 {
		// 1 bit  => 1 byte
		// 8 bits => 1 byte
		// 9 bits => 2 bytes
		bytesNeeded := ((bitsNeeded - 1) / 8) + 1
		if bytesNeeded > 8 {
			err = fmt.Errorf("ReadBitsInt(%d): more than 8 bytes requested", n)
			return
		}
		_, err = reader.Read(reader.buf[:bytesNeeded])
		if err != nil {
			return ret, err
		}
		for i := 0; i < bytesNeeded; i++ {
			reader.bits <<= 8
			reader.bits |= uint64(reader.buf[i])
			reader.bitsLeft += 8
		}
	}

	// raw mask with required number of 1s, starting from the lowest bit
	var mask uint64 = (1 << n) - 1
	// shift "bits" to align the highest bits with the mask & derive the result
	shiftBits := reader.bitsLeft - n
	ret = (reader.bits >> shiftBits) & mask
	// clear top bits that we've just read => AND with 1s
	reader.bitsLeft -= n
	mask = (1 << reader.bitsLeft) - 1
	reader.bits &= mask

	return
}

type littleEndianConverter struct {
}

func (o *littleEndianConverter) Uint(data []byte) uint {
	return uint(data[0])
}

func (o *littleEndianConverter) Uint8(data []byte) uint8 {
	return data[0]
}

func (o *littleEndianConverter) Uint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

func (o *littleEndianConverter) Uint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}
func (o *littleEndianConverter) Uint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}
func (o *littleEndianConverter) Float32fromBits(data []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data))
}
func (o *littleEndianConverter) Float64fromBits(data []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(data))
}

func (o *littleEndianConverter) ReadBitsInt(reader Reader, length uint8) (ret uint64, err error) {
	bitsNeeded := int(length) - int(reader.bitsLeft)
	var bitsLeft = uint64(reader.bitsLeft)
	if bitsNeeded > 0 {
		// 1 bit  => 1 byte
		// 8 bits => 1 byte
		// 9 bits => 2 bytes
		bytesNeeded := ((bitsNeeded - 1) / 8) + 1
		if bytesNeeded > 8 {
			err = fmt.Errorf("ReadBitsIntLe(%d): more than 8 bytes requested", length)
			return
		}
		_, err = reader.Read(reader.buf[:bytesNeeded])
		if err != nil {
			return ret, err
		}
		for i := 0; i < bytesNeeded; i++ {
			reader.bits |= uint64(reader.buf[i]) << bitsLeft
			bitsLeft += 8
		}
	}

	// raw mask with required number of 1s, starting from the lowest bit
	var mask uint64 = (1 << length) - 1
	// derive reading result
	ret = reader.bits & mask
	// remove bottom bits that we've just read by shifting
	reader.bits >>= length
	reader.bitsLeft = uint8(bitsLeft) - length

	return
}
