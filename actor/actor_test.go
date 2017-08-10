package actor

import (
	"github.com/davyxu/actornet/proto"
	"sync"
	"testing"
)

func TestHelloWorld(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(1)

	pid := Spawn("hello", func(c Context) {

		switch msg := c.Msg().(type) {
		case string:
			t.Log(msg)
			wg.Done()
		}

	})

	pid.Exec("hello", nil)

	wg.Wait()
}

func Test2ActorCommunicate(t *testing.T) {

	echoFunc := func(c Context) {

		switch msg := c.Msg().(type) {
		case string:
			t.Log("server recv", msg)

			if c.Source() != nil {
				c.Source().Exec(msg, c.Self())
			}

		}

	}

	server := Spawn("server", echoFunc)

	var wg sync.WaitGroup

	wg.Add(1)

	Spawn("client", func(c Context) {

		switch data := c.Msg().(type) {
		case *proto.Start:
			t.Log("client start")

			server.Exec("hello", c.Self())
		case string:
			t.Log("client recv", data)

			wg.Done()
		}

	})

	wg.Wait()
}
