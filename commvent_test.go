package commvent_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/GiulianoDecesares/commvent/client"
	"github.com/GiulianoDecesares/commvent/primitives"
	"github.com/GiulianoDecesares/commvent/server"
	"github.com/sirupsen/logrus"
)

const (
	testingKey   = "TESTING_KEY"
	testingValue = "TESTING_STRING"
)

var (
	address = url.URL{
		Host: "localhost:5555",
	}

	sendEvents chan bool = make(chan bool)

	testingCommand *primitives.Message = primitives.NewMessage("TESTING_COMMAND", "")
	testingEvent   *primitives.Message = primitives.NewMessage("TESTING_EVENT", "")

	currentSever  *server.Server = server.NewServer(address)
	currentClient *server.Client

	eventCounter int = 0
)

func updateEventCounter(event *primitives.Message, context *testing.T) {
	if event != nil && event.Type == testingEvent.Type {
		eventCounter++
		context.Logf("Event counter: %d", eventCounter)
	}
}

func onNewServerClient(client *server.Client, context *testing.T) {
	context.Logf("[Server] Client connected. Regsitering")
	currentClient = client

	currentClient.HandleDisconnect(func(disconnectionState error) {
		context.Log("[Server] Client disconnected")
		currentClient = nil
	})

	currentClient.HandleEvents(func(event *primitives.Message) {
		context.Log("[Server] Event received")
		updateEventCounter(event, context)
	})

	client.SendCommand(testingCommand)
}

func clientCommandHandler(command *primitives.Message, context *testing.T) {
	context.Log("[Client] Command received")

	if command.Type != testingCommand.Type {
		context.Errorf("[Client] Command type should be %s but is %s", testingCommand.Type, command.Type)
	}

	if value, err := command.GetString(testingKey); err == nil {
		if value != testingValue {
			context.Errorf("[Client] Received command should have %s key with value %s", testingKey, testingValue)
		} else {
			sendEvents <- true
		}
	} else {
		context.Errorf("[Client] Received command should have %s key", testingKey)
	}

}

func Test(context *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)

	testingCommand.SetString(testingKey, testingValue)

	go currentSever.Listen(func(client *server.Client) {
		context.Log("[Server] Server listening")
		onNewServerClient(client, context)
	})

	client := client.NewClient(func(command *primitives.Message) {
		clientCommandHandler(command, context)
	})

	if err := client.Begin(address); err != nil {
		context.Fatalf("Error while starting client: %s", err.Error())
	}

	defer client.Stop()

	if <-sendEvents {
		for index := 0; index < 10; index++ {
			context.Log("[Client] Sending testing event")

			client.SendEvent(testingEvent)
			time.Sleep(time.Second)
		}

		if eventCounter != 10 {
			context.Errorf("Amount of received events after sending 10 is %d and should be %d", eventCounter, 10)
		}
	}
}
