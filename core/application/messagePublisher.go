package application

type MessagePublisher interface {
	PublishMessages(messages []*Message, option map[string]interface{}) error
}