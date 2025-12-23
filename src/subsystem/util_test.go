package subsystem

import (
	"testing"
)

func TestSizeToBytes(t *testing.T) {
	testCases := []string{
		"100m",  // 100MB
		"100M",  // 100MB（大写）
		"1G",    // 1GB
		"1g",    // 1GB（小写）
		"512",   // 512字节
		"2.5KB", // 2.5KB
		"0.5G",  // 0.5GB
	}

	for _, tc := range testCases {
		_, err := SizeToBytes(tc)
		if err != nil {
			t.Error(err)
			continue
		}
	}

	_, err := SizeToBytes("abc")
	if err == nil {
		t.Error("parse abc should return error.")
	}
}
