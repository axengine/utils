package natx

import "context"

type Handle func(subj string, req interface{}) interface{}

type NatXAPP interface {
	Request(context.Context, string, interface{}) (interface{}, error)
	Subscribe(string, Handle) error
	QueueSubscribe(string, Handle) error
	Publish(string, interface{}) error
}
