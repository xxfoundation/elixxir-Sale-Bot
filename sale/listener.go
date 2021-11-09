package sale

import (
	"fmt"
	"git.xx.network/elixxir/sale-bot/storage"
	"github.com/golang/protobuf/proto"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/interfaces/message"
	"gitlab.com/elixxir/client/interfaces/params"
	"net/mail"
	"strings"
)

type listener struct {
	s *storage.Storage
	c *api.Client
}

// Hear messages from users to the coupon bot & respond appropriately
func (l *listener) Hear(item message.Receive) {
	// Confirm that authenticated channels
	if !l.c.HasAuthenticatedChannel(item.Sender) {
		jww.ERROR.Printf("No authenticated channel exists to %+v", item.Sender)
	}

	responseStr := ""

	// Parse the trigger
	in := &CMIXText{}
	err := proto.Unmarshal(item.Payload, in)
	if err != nil {
		jww.ERROR.Printf("Could not unmartial message from messenger: %+v", err)
	}

	addr, err := mail.ParseAddress(in.Text)
	if err != nil {
		responseStr = fmt.Sprintf("Expected a valid email, but could not parse address: %+v", err)
	} else {
		// Do sale stuff
		err = l.s.UpsertMember(addr.String(), item.Sender.String())
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				responseStr = fmt.Sprintf("Error: this account has already registered for the sale")
			}
			responseStr = fmt.Sprintf("Failed to insert sale member to db: %+v", err)
		} else {
			responseStr = "Successfully inserted email to sale database"
		}
	}

	contact, err := l.c.GetAuthenticatedChannelRequest(item.Sender)
	if err != nil {
		jww.ERROR.Printf("Could not get authenticated channel request info: %+v", err)
	}

	out := &CMIXText{Text: responseStr}
	payload, err := proto.Marshal(out)
	if err != nil {
		jww.ERROR.Printf("Failed to marshal proto response: %+v", err)
	}

	// Create response message
	resp := message.Send{
		Recipient:   contact.ID,
		Payload:     payload,
		MessageType: message.Text,
	}

	// Send response message to sender over cmix
	rids, mid, t, err := l.c.SendE2E(resp, params.GetDefaultE2E())
	if err != nil {
		jww.ERROR.Printf("Failed to send message: %+v", err)
	}
	jww.INFO.Printf("Sent ... %s [%+v] to %+v on rounds %+v [%+v]", responseStr, mid, item.Sender.String(), rids, t)
}

// Name returns a name, used for debugging
func (l *listener) Name() string {
	return "sale-listener"
}
