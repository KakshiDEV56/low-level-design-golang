package notificationsystem

import (
	"errors"
	"sync"
	"testing"
)

// recordingChannel is a test double that records every (address, message) pair it was
// asked to send, and always "succeeds."
type recordingChannel struct {
	name string
	mu   sync.Mutex
	sent []sentMessage
}

type sentMessage struct {
	address string
	message string
}

func newRecordingChannel(name string) *recordingChannel {
	return &recordingChannel{name: name}
}

func (c *recordingChannel) Name() string { return c.name }

func (c *recordingChannel) Send(address, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sent = append(c.sent, sentMessage{address, message})
	return nil
}

func (c *recordingChannel) sentCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.sent)
}

// failingChannel always returns an error, to test partial-failure handling.
type failingChannel struct {
	name string
	err  error
}

func (c failingChannel) Name() string                        { return c.name }
func (c failingChannel) Send(address, message string) error { return c.err }

func TestSubscribe_UnregisteredChannelFails(t *testing.T) {
	svc := NewNotificationService()
	if err := svc.Subscribe("user1", "email", "user1@example.com"); err == nil {
		t.Fatal("Subscribe to an unregistered channel = nil error, want an error")
	}
}

func TestNotify_SendsToAllSubscribedChannels(t *testing.T) {
	svc := NewNotificationService()
	email := newRecordingChannel("email")
	sms := newRecordingChannel("sms")
	svc.RegisterChannel(email)
	svc.RegisterChannel(sms)

	if err := svc.Subscribe("user1", "email", "user1@example.com"); err != nil {
		t.Fatalf("Subscribe(email) failed: %v", err)
	}
	if err := svc.Subscribe("user1", "sms", "+1-555-0100"); err != nil {
		t.Fatalf("Subscribe(sms) failed: %v", err)
	}

	errs := svc.Notify("user1", "hello")
	if len(errs) != 0 {
		t.Fatalf("Notify() errors = %v, want none", errs)
	}
	if email.sentCount() != 1 {
		t.Fatalf("email channel received %d sends, want 1", email.sentCount())
	}
	if sms.sentCount() != 1 {
		t.Fatalf("sms channel received %d sends, want 1", sms.sentCount())
	}
	if email.sent[0].address != "user1@example.com" || email.sent[0].message != "hello" {
		t.Fatalf("email.sent[0] = %+v, want address=user1@example.com message=hello", email.sent[0])
	}
}

func TestNotify_UserWithNoSubscriptionsReturnsEmpty(t *testing.T) {
	svc := NewNotificationService()
	svc.RegisterChannel(newRecordingChannel("email"))

	errs := svc.Notify("nobody", "hello")
	if len(errs) != 0 {
		t.Fatalf("Notify() for a user with no subscriptions = %v, want empty", errs)
	}
}

func TestNotify_PartialFailureDoesNotStopOtherChannels(t *testing.T) {
	svc := NewNotificationService()
	email := newRecordingChannel("email")
	failWith := errors.New("sms provider down")
	sms := failingChannel{name: "sms", err: failWith}
	push := newRecordingChannel("push")

	svc.RegisterChannel(email)
	svc.RegisterChannel(sms)
	svc.RegisterChannel(push)
	svc.Subscribe("user1", "email", "user1@example.com")
	svc.Subscribe("user1", "sms", "+1-555-0100")
	svc.Subscribe("user1", "push", "device-token-abc")

	errs := svc.Notify("user1", "hello")
	if len(errs) != 1 {
		t.Fatalf("Notify() errors = %v, want exactly 1 (sms)", errs)
	}
	if !errors.Is(errs["sms"], failWith) {
		t.Fatalf("errs[\"sms\"] = %v, want %v", errs["sms"], failWith)
	}
	// Both OTHER channels must still have been attempted despite sms failing.
	if email.sentCount() != 1 {
		t.Fatalf("email channel received %d sends, want 1 (must not be blocked by sms failure)", email.sentCount())
	}
	if push.sentCount() != 1 {
		t.Fatalf("push channel received %d sends, want 1 (must not be blocked by sms failure)", push.sentCount())
	}
}

func TestUnsubscribe_StopsDelivery(t *testing.T) {
	svc := NewNotificationService()
	email := newRecordingChannel("email")
	svc.RegisterChannel(email)
	svc.Subscribe("user1", "email", "user1@example.com")

	svc.Unsubscribe("user1", "email")
	svc.Notify("user1", "hello")

	if email.sentCount() != 0 {
		t.Fatalf("email channel received %d sends after Unsubscribe, want 0", email.sentCount())
	}
}

func TestUnsubscribe_UnknownSubscriptionIsNoOp(t *testing.T) {
	svc := NewNotificationService()
	svc.RegisterChannel(newRecordingChannel("email"))
	svc.Unsubscribe("nobody", "email") // must not panic
}

func TestResubscribe_UpdatesAddressWithoutDuplicateSend(t *testing.T) {
	svc := NewNotificationService()
	email := newRecordingChannel("email")
	svc.RegisterChannel(email)

	svc.Subscribe("user1", "email", "old@example.com")
	svc.Subscribe("user1", "email", "new@example.com") // re-subscribe with a new address

	svc.Notify("user1", "hello")
	if email.sentCount() != 1 {
		t.Fatalf("email channel received %d sends, want exactly 1 (re-subscribing shouldn't duplicate)", email.sentCount())
	}
	if email.sent[0].address != "new@example.com" {
		t.Fatalf("email.sent[0].address = %q, want %q (the updated address)", email.sent[0].address, "new@example.com")
	}
}

func TestRegisterChannel_ReplacesExistingByName(t *testing.T) {
	svc := NewNotificationService()
	oldChannel := newRecordingChannel("email")
	newChannel := newRecordingChannel("email")
	svc.RegisterChannel(oldChannel)
	svc.RegisterChannel(newChannel) // same Name(), should replace

	svc.Subscribe("user1", "email", "user1@example.com")
	svc.Notify("user1", "hello")

	if oldChannel.sentCount() != 0 {
		t.Fatalf("old channel received %d sends, want 0 (should have been replaced)", oldChannel.sentCount())
	}
	if newChannel.sentCount() != 1 {
		t.Fatalf("new channel received %d sends, want 1", newChannel.sentCount())
	}
}

func TestNotify_MultipleUsersIndependent(t *testing.T) {
	svc := NewNotificationService()
	email := newRecordingChannel("email")
	svc.RegisterChannel(email)
	svc.Subscribe("user1", "email", "user1@example.com")
	// user2 never subscribes.

	svc.Notify("user1", "hi user1")
	svc.Notify("user2", "hi user2")

	if email.sentCount() != 1 {
		t.Fatalf("email channel received %d sends, want 1 (only user1 is subscribed)", email.sentCount())
	}
}

// carrierPigeonChannel is a BRAND NEW channel type, defined only here in the test
// file. If it works with RegisterChannel/Subscribe/Notify without any change to
// solution.go, that proves NotificationService is genuinely extensible.
type carrierPigeonChannel struct {
	delivered []string
}

func (c *carrierPigeonChannel) Name() string { return "carrier-pigeon" }
func (c *carrierPigeonChannel) Send(address, message string) error {
	c.delivered = append(c.delivered, address+": "+message)
	return nil
}

func TestNotify_ExtensibilityWithBrandNewChannelType(t *testing.T) {
	svc := NewNotificationService()
	pigeon := &carrierPigeonChannel{}
	svc.RegisterChannel(pigeon)

	if err := svc.Subscribe("user1", "carrier-pigeon", "Coop #4"); err != nil {
		t.Fatalf("Subscribe to the new channel type failed: %v", err)
	}
	svc.Notify("user1", "hello via pigeon")

	if len(pigeon.delivered) != 1 || pigeon.delivered[0] != "Coop #4: hello via pigeon" {
		t.Fatalf("pigeon.delivered = %v, want [\"Coop #4: hello via pigeon\"]", pigeon.delivered)
	}
}
