package client

import (
	"net/url"

	"github.com/GiulianoDecesares/commvent/primitives"
)

type IClient interface {
	Begin(url url.URL) error
	Stop() error

	SendEvent(event *primitives.Message)
}
