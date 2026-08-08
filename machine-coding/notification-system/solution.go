package notificationsystem

/*
See README.md for the full problem statement, requirements, and rubric.

Channel is a fixed interface — do not change its shape. Deliberately, no concrete
Channel implementations (email/SMS/push) are provided here: the whole point of this
exercise is that NotificationService works with ANY Channel implementation without
needing to know about it, which is proven by solution_test.go defining its own
Channel implementations entirely outside this file.

TODO: implement NotificationService.
  - NewNotificationService() *NotificationService
  - RegisterChannel(ch Channel)
  - Subscribe(userID, channelName, address string) error
  - Unsubscribe(userID, channelName string)
  - Notify(userID, message string) map[string]error

See README.md hints for a suggested internal data-structure layout.
*/

// Channel is one way of delivering a notification (email, SMS, push, or anything
// else). NotificationService only ever depends on this interface.
type Channel interface {
	Name() string
	Send(address string, message string) error
}

type NotificationService struct {
	// TODO: add your fields here
}

func NewNotificationService() *NotificationService {
	// TODO: implement
	return &NotificationService{}
}

func (n *NotificationService) RegisterChannel(ch Channel) {
	// TODO: implement
}

func (n *NotificationService) Subscribe(userID, channelName, address string) error {
	// TODO: implement
	return nil
}

func (n *NotificationService) Unsubscribe(userID, channelName string) {
	// TODO: implement
}

func (n *NotificationService) Notify(userID, message string) map[string]error {
	// TODO: implement
	return nil
}
