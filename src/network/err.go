package network

import (
	"errors"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

func IsResourceNotFound(err error) bool {
	return err != nil && errors.Is(err, constant.ErrResourceNotFound)
}
