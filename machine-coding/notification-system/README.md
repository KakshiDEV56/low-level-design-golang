# Multi-Channel Notification System

**Suggested time box: 45 minutes.**

## Problem statement

Design a notification service that can send a message to a user across whichever
channels (email, SMS, push, or anything added later) that user has subscribed to — and
critically, **the core service must never need to change when a new channel type is
added.**

## Functional requirements

1. `Channel` is a fixed interface (given below — don't change its shape):
   `Name() string` and `Send(address, message string) error`.
2. `NewNotificationService()` creates an empty service (no channels registered, no
   subscriptions).
3. `RegisterChannel(ch Channel)` adds support for a channel type, keyed by
   `ch.Name()`. Registering a channel with a name that's already registered replaces
   the old one.
4. `Subscribe(userID, channelName, address string) error` records that `userID` wants
   notifications via `channelName` (e.g. `"email"`) sent to `address` (e.g. an email
   address, phone number, or device token — just an opaque string as far as this
   service is concerned). Returns an error if `channelName` was never registered via
   `RegisterChannel`. Subscribing again for the same `(userID, channelName)` pair
   updates the address.
5. `Unsubscribe(userID, channelName string)` removes that one subscription, if it
   exists (no-op otherwise).
6. `Notify(userID, message string) map[string]error` sends `message` via **every**
   channel `userID` is currently subscribed to. It must attempt ALL of them even if one
   fails partway — a failure on one channel must never prevent delivery attempts on the
   others. Returns a map of `channelName -> error` for only the channels that failed
   (channels that succeeded, and a user with no subscriptions at all, both result in an
   empty map — never `nil` vs. empty being meaningfully different, so returning an
   empty non-nil map is fine either way as long as `len() == 0`).

## Rubric

- [ ] `NotificationService` contains **zero** channel-specific logic (no `if
      channelName == "email"` anywhere) — it only ever calls the `Channel` interface.
- [ ] A brand-new `Channel` implementation, written entirely outside `solution.go`,
      works with `RegisterChannel`/`Subscribe`/`Notify` with no changes to
      `solution.go` — this is the concrete proof of the extensibility requirement.
- [ ] One channel failing during `Notify` does not stop delivery attempts on the
      remaining channels for that user.
- [ ] `Subscribe` on an unregistered channel name fails cleanly (no panic, clear
      error).
- [ ] Re-subscribing updates the address rather than creating a duplicate
      subscription.

## Hints

- Store registered channels as `map[string]Channel` (name -> implementation).
- Store subscriptions as `map[string]map[string]string` — `userID -> channelName ->
  address`. This makes `Unsubscribe` and re-`Subscribe` both simple map operations.
- `Notify` just iterates the user's subscription map, looks up each channel by name in
  the registered-channels map, calls `Send`, and collects errors — no need for
  goroutines/concurrency here unless you want the extra challenge of sending
  concurrently and aggregating results safely.
