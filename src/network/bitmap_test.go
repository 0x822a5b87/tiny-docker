package network

import (
	"fmt"
	"testing"

	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/stretchr/testify/assert"
)

func TestNewBitmap(t *testing.T) {
	bitmap, err := NewBitmap(0)
	assert.ErrorIs(t, err, ErrInvalidSize)
	assert.Nil(t, bitmap)

	bitmap, err = NewBitmap(14)
	assert.NoError(t, err)
	assert.Equal(t, uint64(16), bitmap.Size)
	assert.Len(t, bitmap.Bits, 1)

	bitmap, err = NewBitmap(64)
	assert.NoError(t, err)
	assert.Equal(t, uint64(64), bitmap.Size)
	assert.Len(t, bitmap.Bits, 1)

	bitmap, err = NewBitmap(65)
	assert.NoError(t, err)
	assert.Equal(t, uint64(128), bitmap.Size)
	assert.Len(t, bitmap.Bits, 2)
}

func TestBitmap_Set_Clear_IsSet(t *testing.T) {
	bitmap, err := NewBitmap(16)
	assert.NoError(t, err)

	err = bitmap.Set(5)
	assert.NoError(t, err)
	assert.True(t, bitmap.IsSet(5))

	err = bitmap.Set(20)
	assert.ErrorIs(t, err, ErrInvalidPos)
	assert.False(t, bitmap.IsSet(20))

	err = bitmap.Clear(5)
	assert.NoError(t, err)
	assert.False(t, bitmap.IsSet(5))

	err = bitmap.Clear(20)
	assert.ErrorIs(t, err, ErrInvalidPos)
}

func TestBitmap_FindFirstUnset(t *testing.T) {
	bitmap, err := NewBitmap(16)
	assert.NoError(t, err)
	pos := bitmap.FindFirstUnset()
	assert.Equal(t, 0, pos)

	err = bitmap.Set(0)
	assert.NoError(t, err)
	pos = bitmap.FindFirstUnset()
	assert.Equal(t, 1, pos)

	for i := 0; i < 16; i++ {
		err = bitmap.Set(uint64(i))
		assert.NoError(t, err)
	}
	pos = bitmap.FindFirstUnset()
	assert.Equal(t, -1, pos)
}

func TestNewIPNetBitmap(t *testing.T) {
	ipNetBitmap, err := NewIPNetBitmap(16, 256)
	assert.NoError(t, err)
	assert.Equal(t, uint64(16), ipNetBitmap.SubnetBitmap.Size)
	assert.Equal(t, uint64(256), ipNetBitmap.IpSize)
	assert.Empty(t, ipNetBitmap.SubnetIPMaps)
}

func TestIPNetBitmap_AllocateSubnet(t *testing.T) {
	ipNetBitmap, err := NewIPNetBitmap(16, 256)
	assert.NoError(t, err)

	ipBitmap, subnetPos, err := ipNetBitmap.AllocateSubnet()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), subnetPos)
	assert.Equal(t, uint64(256), ipBitmap.Size)
	assert.True(t, ipNetBitmap.SubnetBitmap.IsSet(0))
	assert.Contains(t, ipNetBitmap.SubnetIPMaps, uint64(0))

	ipBitmap2, subnetPos2, err := ipNetBitmap.AllocateSubnet()
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), subnetPos2)
	assert.Equal(t, uint64(256), ipBitmap2.Size)

	for i := 2; i < 16; i++ {
		_, _, err := ipNetBitmap.AllocateSubnet()
		assert.NoError(t, err)
	}
	_, _, err = ipNetBitmap.AllocateSubnet()
	assert.ErrorIs(t, err, ErrOutOfRange)
}

func TestIPNetBitmap_ReleaseSubnet(t *testing.T) {
	ipNetBitmap, err := NewIPNetBitmap(16, 256)
	assert.NoError(t, err)

	_, subnetPos, err := ipNetBitmap.AllocateSubnet()
	assert.NoError(t, err)
	assert.True(t, ipNetBitmap.SubnetBitmap.IsSet(subnetPos))

	err = ipNetBitmap.ReleaseSubnet(subnetPos)
	assert.NoError(t, err)
	assert.False(t, ipNetBitmap.SubnetBitmap.IsSet(subnetPos))

	ipBitmap := ipNetBitmap.SubnetIPMaps[subnetPos]
	pos := ipBitmap.FindFirstUnset()
	assert.Equal(t, 0, pos)

	err = ipNetBitmap.ReleaseSubnet(20)
	assert.ErrorIs(t, err, ErrInvalidPos)
}

func TestIPNetBitmap_AllocateIPInSubnet(t *testing.T) {
	ipNetBitmap, err := NewIPNetBitmap(16, 256)
	assert.NoError(t, err)

	_, subnetPos, err := ipNetBitmap.AllocateSubnet()
	assert.NoError(t, err)

	ipPos, err := ipNetBitmap.AllocateIPInSubnet(subnetPos)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), ipPos)
	assert.True(t, ipNetBitmap.SubnetIPMaps[subnetPos].IsSet(ipPos))

	ipPos2, err := ipNetBitmap.AllocateIPInSubnet(subnetPos)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), ipPos2)

	err = ipNetBitmap.ReleaseSubnet(subnetPos)
	assert.NoError(t, err)
	_, err = ipNetBitmap.AllocateIPInSubnet(subnetPos)
	assert.ErrorIs(t, err, ErrInvalidPos)

	_ = ipNetBitmap.SubnetIPMaps[subnetPos]
	_, _, err = ipNetBitmap.AllocateSubnet()
	assert.NoError(t, err)
	for i := 0; i < 256; i++ {
		_, err := ipNetBitmap.AllocateIPInSubnet(subnetPos)
		assert.NoError(t, err)
	}
	_, err = ipNetBitmap.AllocateIPInSubnet(subnetPos)
	assert.ErrorIs(t, err, ErrOutOfRange)
}

func TestIPNetBitmap_Integration(t *testing.T) {
	ipNetBitmap, err := NewIPNetBitmap(16, 65536)
	assert.NoError(t, err)

	ipBitmap, subnetPos, err := ipNetBitmap.AllocateSubnet()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), subnetPos)
	assert.Equal(t, uint64(65536), ipBitmap.Size)

	ipPos, err := ipNetBitmap.AllocateIPInSubnet(subnetPos)
	for i := 1; i < 100; i++ {
		ipPos, err = ipNetBitmap.AllocateIPInSubnet(subnetPos)
		assert.NoError(t, err)
	}
	assert.Equal(t, uint64(99), ipPos)
	assert.True(t, ipBitmap.IsSet(ipPos))

	err = ipNetBitmap.ReleaseSubnet(subnetPos)
	assert.NoError(t, err)

	ipBitmap2, _, err := ipNetBitmap.AllocateSubnet()
	assert.NoError(t, err)
	assert.Equal(t, 0, ipBitmap2.FindFirstUnset())
}

func TestNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		input  uint64
		output uint64
	}{
		{0, 1},
		{1, 1},
		{14, 16},
		{16, 16},
		{65, 128},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input_%d", tt.input), func(t *testing.T) {
			res := util.NextPowerOfTwo(tt.input)
			assert.Equal(t, tt.output, res)
		})
	}
}
