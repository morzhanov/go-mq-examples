package mq

// MQ defines MQ interface for testing
// each MQ should produce some messages and listen to them
type MQ interface {
	Produce()
	Listen()
}
