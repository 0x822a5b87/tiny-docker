package network

import (
	"errors"
	"fmt"
	"math/bits"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/util"
)

type Bitmap struct {
	Size uint64   `json:"size"`
	Bits []uint64 `json:"bits"`
}

func NewBitmap(size uint64) (*Bitmap, error) {
	if size == 0 {
		return nil, constant.ErrInvalidSize
	}

	size = util.NextPowerOfTwo(size)

	bitsLen := size / 64
	if size%64 != 0 {
		bitsLen += 1
	}

	return &Bitmap{
		Size: size,
		Bits: make([]uint64, bitsLen),
	}, nil
}

func (b *Bitmap) Set(pos uint64) error {
	if pos >= b.Size {
		return constant.ErrInvalidPos
	}

	idx := pos / 64
	offset := pos % 64
	b.Bits[idx] |= 1 << offset
	return nil
}

func (b *Bitmap) Clear(pos uint64) error {
	if pos >= b.Size {
		return constant.ErrInvalidPos
	}

	idx := pos / 64
	offset := pos % 64
	b.Bits[idx] &^= 1 << offset
	return nil
}

func (b *Bitmap) IsSet(pos uint64) bool {
	if pos >= b.Size {
		return false
	}

	idx := pos / 64
	offset := pos % 64
	return (b.Bits[idx] & (1 << offset)) != 0
}

func (b *Bitmap) FindFirstUnset() int {
	for i, chunk := range b.Bits {
		if chunk == ^uint64(0) {
			continue
		}
		posInChunk := bits.TrailingZeros64(^chunk)
		totalPos := i*64 + posInChunk
		if uint64(totalPos) >= b.Size {
			continue
		}
		return totalPos
	}
	return -1
}

type IPNetBitmap struct {
	SubnetBitmap *Bitmap            `json:"subnet_bitmap"`
	IpSize       uint64             `json:"ip_size"`
	SubnetIPMaps map[uint64]*Bitmap `json:"subnet_ip_maps"`
}

func NewIPNetBitmap(subnetSize, ipSize uint64) (*IPNetBitmap, error) {
	subnetBitmap, err := NewBitmap(subnetSize)
	if err != nil {
		return nil, fmt.Errorf("create subnet bitmap failed: %w", err)
	}

	return &IPNetBitmap{
		SubnetBitmap: subnetBitmap,
		IpSize:       util.NextPowerOfTwo(ipSize),
		SubnetIPMaps: make(map[uint64]*Bitmap),
	}, nil
}

func (i *IPNetBitmap) AllocateSubnet() (*Bitmap, uint64, error) {
	subnetPos := i.SubnetBitmap.FindFirstUnset()
	if subnetPos == -1 {
		return nil, 0, constant.ErrOutOfRange
	}
	subnetPosUint := uint64(subnetPos)

	if err := i.SubnetBitmap.Set(subnetPosUint); err != nil {
		return nil, 0, err
	}

	if _, exists := i.SubnetIPMaps[subnetPosUint]; !exists {
		ipBitmap, err := NewBitmap(i.IpSize)
		if err != nil {
			return nil, 0, fmt.Errorf("create ip bitmap failed: %w", err)
		}
		i.SubnetIPMaps[subnetPosUint] = ipBitmap
	}

	return i.SubnetIPMaps[subnetPosUint], subnetPosUint, nil
}

func (i *IPNetBitmap) ReleaseSubnet(subnetPos uint64) error {
	if !i.SubnetBitmap.IsSet(subnetPos) {
		return constant.ErrInvalidPos
	}

	if err := i.SubnetBitmap.Clear(subnetPos); err != nil {
		return err
	}

	if ipBitmap, exists := i.SubnetIPMaps[subnetPos]; exists {
		ipBitmap.Bits = make([]uint64, len(ipBitmap.Bits)) // 重置为全 0
	}

	return nil
}

func (i *IPNetBitmap) AllocateIPInSubnet(subnetPos uint64) (uint64, error) {
	if !i.SubnetBitmap.IsSet(subnetPos) {
		return 0, constant.ErrInvalidPos
	}

	ipBitmap, exists := i.SubnetIPMaps[subnetPos]
	if !exists {
		return 0, errors.New("subnet ip bitmap not initialized")
	}

	ipPos := ipBitmap.FindFirstUnset()
	if ipPos == -1 {
		return 0, constant.ErrOutOfRange
	}
	ipPosUint := uint64(ipPos)

	if err := ipBitmap.Set(ipPosUint); err != nil {
		return 0, err
	}

	return ipPosUint, nil
}
