package events

import (
	"encoding/binary"
	"fmt"
	"github.com/FactomProject/live-feed-api/EventRouter/config"
	"github.com/FactomProject/live-feed-api/EventRouter/eventmessages/generated/eventmessages"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

var eventsQueue chan *eventmessages.FactomEvent
var address string

func init() {
	configuration := &config.ReceiverConfig{
		Protocol:    "tcp",
		BindAddress: "",
		Port:        0,
	}

	// Start the new server at random port
	server := NewReceiver(configuration)
	server.Start()
	time.Sleep(10 * time.Millisecond) // sleep to allow the server to start before making a connection
	address = server.GetAddress()
	fmt.Printf("start server at: '%s'\n", address)
	eventsQueue = server.GetEventQueue()
}

func TestWriteEvents(t *testing.T) {
	n := 10
	data := mockData(t)
	dataSize := int32(len(data))

	conn := connect(t)
	defer conn.Close()

	for i := 1; i <= n; i++ {
		err := binary.Write(conn, binary.LittleEndian, supportedProtocolVersion)
		if err != nil {
			t.Fatal(err)
		}
		err = binary.Write(conn, binary.LittleEndian, dataSize)
		if err != nil {
			t.Fatal(err)
		}

		status, err := conn.Write(data)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("bytes offered: %d, bytes written: %v", dataSize, status)
	}

	correctSendEvents := listenForEvents("WRITE", n, 20*time.Second)
	t.Logf("number of events sent: %d, number of events received: %d", n, correctSendEvents)
	assert.EqualValues(t, n, correctSendEvents, "failed to receive the correct number of events %d != %d", n, correctSendEvents)
}

func TestEOFConnection(t *testing.T) {
	n := 10
	data := mockData(t)
	dataSize := int32(len(data))

	// test in parallel
	for i := 0; i < n; i++ {
		go func() {
			// prevent every thread making connection at the same time
			r := rand.Intn(10)
			time.Sleep(time.Duration(r) * time.Millisecond)

			conn := connect(t)
			defer conn.Close()

			// send one event correctly
			binary.Write(conn, binary.LittleEndian, supportedProtocolVersion)
			err := binary.Write(conn, binary.LittleEndian, dataSize)
			if err != nil {
				t.Fatal(err)
			}
			_, err = conn.Write(data)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}

	correctSendEvents := listenForEvents("EOF", n, 20*time.Second)
	assert.EqualValues(t, n, correctSendEvents, "failed to receive the correct number of events %d != %d", n, correctSendEvents)
}

func connect(t *testing.T) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func listenForEvents(testID string, n int, timeLimit time.Duration) int {
	var correctSendEvents int32 = 0
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-eventsQueue:
				atomic.AddInt32(&correctSendEvents, 1)
				fmt.Printf("[%s] > received event in queue\n", testID)
			case <-quit:
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	deadline := time.Now().Add(timeLimit)
	for int(atomic.LoadInt32(&correctSendEvents)) != n && time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
	}
	quit <- true
	return int(correctSendEvents)
}

func mockData(t *testing.T) []byte {
	event := eventmessages.NewPopulatedFactomEvent(randomizer, true)
	data, err := proto.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
