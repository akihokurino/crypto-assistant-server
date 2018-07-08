package utils

import (
	"fmt"
	"time"
	"encoding/binary"
	"crypto/rand"
	"strconv"
)

type IDUtil interface {
	MakeRandomKey() string
}

type iDUtil struct {
}

func NewIDUtil() IDUtil {
	return &iDUtil{}
}

func (u *iDUtil) MakeRandomKey() string {
	return fmt.Sprintf("%s-%d", u.random(), 10000000000-time.Now().Unix())
}

func (u *iDUtil) random() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}