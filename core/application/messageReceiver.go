package application

type MessageReceiver interface {
	ReceiveMessage(originalMessage interface{}, option map[string]interface{}) (*Message, bool, error)
	ConfirmReceive(message *Message) error
}
