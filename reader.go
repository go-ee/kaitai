package kaitai

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type ReadFix func(reader ReaderIO) (ret interface{}, err error)
type ReadDynamic func(reader ReaderIO, parent *Item, root *Item) (ret interface{}, err error)

type ItemReader interface {
	Read(reader ReaderIO, parent *Item, root *Item) (ret *Item, err error)
}

type Convert func(data []byte) (ret interface{}, err error)

type Item struct {
	Attr     *Attr
	Accessor interface{}
	Value    interface{}
}

func (o *Item) Expr(expr string) (ret *Item, err error) {
	ret = o
	for _, part := range strings.Split(expr, ".") {
		if value, ok := ret.Value.(map[string]*Item); ok {
			if ret = value[part]; ret == nil {
				err = fmt.Errorf("can't resolve '%v' of expression '%v'", part, expr)
				break
			}
		} else {
			err = fmt.Errorf("can't resolve '%v' of expression '%v'", part, expr)
			break
		}
	}
	return
}

type ReaderIO struct {
	io.ReadSeeker
	buf [8]byte

	// Number of bits remaining in "bits" for sequential calls to ReadBitsInt
	bitsLeft uint8
	bits     uint64
}

func (o *ReaderIO) ReadBytes(n uint8) (ret []byte, err error) {
	if n < 0 {
		err = fmt.Errorf("ReadBytes(%d): negative number of bytes to read", n)
		return
	}

	ret = make([]byte, n)
	_, err = io.ReadFull(o, ret)
	return
}

func (o *ReaderIO) ReadBytesFull() ([]byte, error) {
	return ioutil.ReadAll(o)
}

func ReadFixAttr(attr *Attr, convert Convert) (ret ReadFix) {
	if attr.SizeEos == "true" {
		ret = ReadFixFull(convert)
	} else if attr.Size != "" {
		if length, err := strconv.Atoi(attr.Size); err == nil {
			ret = ReadFixLength(uint8(length), convert)
		}
	}
	return
}

func ReadDynamicAttr(attr *Attr, convert Convert) (ret ReadDynamic) {
	if attr.Size != "" {
		if _, err := strconv.Atoi(attr.Size); err != nil {
			ret = ReadDynamicLengthExpr(attr.Size, convert)
		}
	}
	return
}

func ReadFixFull(convert Convert) (ret ReadFix) {
	return func(reader ReaderIO) (ret interface{}, err error) {
		return ReadFull(reader, convert)
	}
}

func ReadFull(reader ReaderIO, convert Convert) (ret interface{}, err error) {
	var data []byte
	if data, err = reader.ReadBytesFull(); err == nil {
		ret, err = convert(data)
	}
	return
}

func ReadFixLength(length uint8, convert Convert) (ret ReadFix) {
	return func(reader ReaderIO) (ret interface{}, err error) {
		return ReadLength(reader, length, convert)
	}
}

func ReadLength(reader ReaderIO, length uint8, convert Convert) (ret interface{}, err error) {
	var data []byte
	if length > 0 {
		data, err = reader.ReadBytes(length)
	} else {
		data, err = reader.ReadBytesFull()
	}

	if err == nil {
		ret, err = convert(data)
	}
	return
}

func ReadDynamicLengthExpr(expr string, convert Convert) (ret ReadDynamic) {
	return func(reader ReaderIO, parent *Item, root *Item) (ret interface{}, err error) {
		var sizeItem *Item
		if sizeItem, err = parent.Expr(expr); err != nil {
			if length, ok := sizeItem.Value.(uint8); ok {
				ret, err = ReadLength(reader, length, convert)
			} else {
				err = fmt.Errorf("cant parse Size to uint8, expr=%v, valiue=%v", expr, sizeItem.Value)
			}
		}
		return
	}
}

func ToString(data []byte) (ret interface{}, _ error) {
	ret = string(data)
	return
}

func ToSame(data []byte) (ret interface{}, _ error) {
	ret = data
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
	ReadBitsInt(reader ReaderIO, n uint8) (ret uint64, err error)
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

func (o *bigEndianConverter) ReadBitsInt(reader ReaderIO, n uint8) (ret uint64, err error) {
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

func (o *littleEndianConverter) ReadBitsInt(reader ReaderIO, length uint8) (ret uint64, err error) {
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
