package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter     uint32
	counterLock sync.Mutex
	lastNano    int64
)

func GenerateUniqeID(prefix *string) string {
	counterLock.Lock()
	defer counterLock.Unlock()

	nowNano := time.Now().UnixNano()

	if nowNano == lastNano {
		counter++
	} else {
		counter = 0
		lastNano = nowNano
	}

	var orderID string
	if prefix != nil && *prefix != "" {
		orderID = fmt.Sprintf("%s%d%02d", *prefix, nowNano, counter)
	} else {
		orderID = fmt.Sprintf("TRX%d%02d", nowNano, counter)
	}

	return orderID
}
