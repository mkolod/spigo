// eureka tests
package eureka

import (
	"fmt"
	"github.com/adrianco/spigo/archaius"
	"github.com/adrianco/spigo/edda"
	"github.com/adrianco/spigo/gotocol"
	"testing"
	"time"
)

// Test the discovery process by writing and reading back service information
func TestDiscovery(t *testing.T) {
	fmt.Println("eureka_test start")
	listener := make(chan gotocol.Message)
	edda.Logchan = make(chan gotocol.Message, 10) // buffered channel
	go edda.Start()
	archaius.Conf.Msglog = true
	archaius.Conf.GraphjsonFile = "test"
	archaius.Conf.GraphmlFile = "test"
	eureka := make(chan gotocol.Message, 10)
	go Start(eureka)
	// stack up a series of requests in the buffered channel
	eureka <- gotocol.Message{gotocol.Hello, listener, time.Now(), "test0" + " " + "test"}
	eureka <- gotocol.Message{gotocol.Hello, listener, time.Now(), "test1" + " " + "test"}
	eureka <- gotocol.Message{gotocol.Hello, listener, time.Now(), "thing0" + " " + "thing"}
	eureka <- gotocol.Message{gotocol.GetRequest, listener, time.Now(), "test0"}
	eureka <- gotocol.Message{gotocol.Goodbye, listener, time.Now(), ""}
	// pick up responses until we see the Googbye response
	for {
		msg := <-listener
		if archaius.Conf.Msglog {
			fmt.Printf("test_eureka: %v\n", msg)
		}
		if msg.Imposition == gotocol.Goodbye {
			break
		}
		switch msg.Imposition {
		case gotocol.GetResponse:
			if msg.Intention != "test" {
				t.Fail()
			}
		}
	}
	if edda.Logchan != nil {
		for {
			//fmt.Printf("Logger has %v messages left to flush\n", len(edda.Logchan))
			if len(edda.Logchan) == 0 {
				break
			}
			time.Sleep(time.Second)
		}
	}
	//wait until edda and eureak finish flushing and close files
	Wg.Wait()
	edda.Wg.Wait()
	fmt.Println("eureka_test end")
}
