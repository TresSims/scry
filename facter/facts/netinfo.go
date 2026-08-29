package facts

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"strconv"
	"strings"
	"time"
)

type NetInfo struct {
	Rx int
	Tx int
}

func Net(ctx context.Context, publish func(val any)) error {
	timer := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-timer.C:
			netbytes, err := os.ReadFile("/proc/net/dev")
			if err != nil {
				return err
			}

			netInfo, err := processNetBytes(netbytes)
			if err != nil {
				return err
			}

			publish(netInfo)
		case <-ctx.Done():
			return nil
		}
	}
}

const (
	Interface int = iota
	RXBytes
	RXPackets
	RXErrors
	RXDrops
	RXFIFO
	RXFrame
	RXCompressed
	RXMulticast
	TXBytes
	TXPackets
	TXErrors
	TXDrops
	TXFIFO
	TXColls
	TXCarrier
	TXCompressed
)

// processNetBytes takes byte reads from /proc/net/dev data
// and return a netinfo object with Rx and Tx calculated by summing
// all received and transmitted bytes
func processNetBytes(netbytes []byte) (*NetInfo, error) {
	netInfo := &NetInfo{}
	netInfoReader := bufio.NewScanner(bytes.NewReader(netbytes))
	netInfoReader.Split(bufio.ScanLines)

	// skip header sections
	netInfoReader.Scan()
	netInfoReader.Scan()

	for netInfoReader.Scan() {
		fields := strings.Fields(netInfoReader.Text())
		rxInt, err := strconv.Atoi(fields[RXBytes])
		if err != nil {
			return nil, err
		}

		txInt, err := strconv.Atoi(fields[TXBytes])
		if err != nil {
			return nil, err
		}

		netInfo.Rx += rxInt
		netInfo.Tx += txInt
	}

	return netInfo, nil
}
