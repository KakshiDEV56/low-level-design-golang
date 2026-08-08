package cqrs

import (
	"errors"
	"reflect"
	"testing"
)

func TestCommandHandler_AddStockRejectsNonPositiveQty(t *testing.T) {
	c := NewCommandHandler()
	if err := c.AddStock("widget", 0); err == nil {
		t.Fatal("AddStock(widget, 0) = nil error, want an error")
	}
	if err := c.AddStock("widget", -5); err == nil {
		t.Fatal("AddStock(widget, -5) = nil error, want an error")
	}
}

func TestCommandHandler_RemoveStockRejectsInsufficientStock(t *testing.T) {
	c := NewCommandHandler()
	c.AddStock("widget", 5)
	err := c.RemoveStock("widget", 10)
	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("RemoveStock(widget, 10) with only 5 in stock = %v, want ErrInsufficientStock", err)
	}
}

func TestCommandHandler_RemoveStockSuccess(t *testing.T) {
	c := NewCommandHandler()
	c.AddStock("widget", 10)
	if err := c.RemoveStock("widget", 4); err != nil {
		t.Fatalf("RemoveStock(widget, 4) failed: %v", err)
	}
}

func TestCommandHandler_EmitsEventsToSubscribers(t *testing.T) {
	c := NewCommandHandler()
	var got []StockEvent
	c.Subscribe(func(e StockEvent) { got = append(got, e) })

	c.AddStock("widget", 10)
	c.RemoveStock("widget", 3)

	want := []StockEvent{
		{ProductID: "widget", Delta: 10},
		{ProductID: "widget", Delta: -3},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("emitted events = %v, want %v", got, want)
	}
}

func TestCommandHandler_FailedCommandsDoNotEmitEvents(t *testing.T) {
	c := NewCommandHandler()
	var got []StockEvent
	c.Subscribe(func(e StockEvent) { got = append(got, e) })

	c.AddStock("widget", -1)         // invalid, should not emit
	c.RemoveStock("widget", 100)     // insufficient stock, should not emit
	if len(got) != 0 {
		t.Fatalf("emitted events after only failed commands = %v, want none", got)
	}
}

func TestCommandHandler_MultipleSubscribersAllNotified(t *testing.T) {
	c := NewCommandHandler()
	var a, b []StockEvent
	c.Subscribe(func(e StockEvent) { a = append(a, e) })
	c.Subscribe(func(e StockEvent) { b = append(b, e) })

	c.AddStock("widget", 5)

	if len(a) != 1 || len(b) != 1 {
		t.Fatalf("subscriber event counts = (%d, %d), want (1, 1)", len(a), len(b))
	}
}

func TestQueryHandler_AppliesEventsAndTracksStock(t *testing.T) {
	q := NewQueryHandler()
	q.Apply(StockEvent{ProductID: "widget", Delta: 10})
	q.Apply(StockEvent{ProductID: "widget", Delta: -3})
	q.Apply(StockEvent{ProductID: "gadget", Delta: 5})

	if got := q.GetStock("widget"); got != 7 {
		t.Fatalf("GetStock(widget) = %d, want 7", got)
	}
	if got := q.GetStock("gadget"); got != 5 {
		t.Fatalf("GetStock(gadget) = %d, want 5", got)
	}
	if got := q.GetStock("never-seen"); got != 0 {
		t.Fatalf("GetStock(never-seen) = %d, want 0", got)
	}
}

func TestQueryHandler_TotalStock(t *testing.T) {
	q := NewQueryHandler()
	q.Apply(StockEvent{ProductID: "a", Delta: 3})
	q.Apply(StockEvent{ProductID: "b", Delta: 7})
	if got := q.TotalStock(); got != 10 {
		t.Fatalf("TotalStock() = %d, want 10", got)
	}
}

func TestQueryHandler_LowStockProductsSorted(t *testing.T) {
	q := NewQueryHandler()
	q.Apply(StockEvent{ProductID: "zebra", Delta: 1})
	q.Apply(StockEvent{ProductID: "apple", Delta: 2})
	q.Apply(StockEvent{ProductID: "mango", Delta: 100})

	got := q.LowStockProducts(5)
	want := []string{"apple", "zebra"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LowStockProducts(5) = %v, want %v", got, want)
	}
}

func TestQueryHandler_EmptyBeforeAnyEvents(t *testing.T) {
	q := NewQueryHandler()
	if got := q.TotalStock(); got != 0 {
		t.Fatalf("TotalStock() before any events = %d, want 0", got)
	}
	if got := q.LowStockProducts(100); len(got) != 0 {
		t.Fatalf("LowStockProducts() before any events = %v, want empty", got)
	}
}

func TestEndToEnd_WiringCommandToQuery(t *testing.T) {
	cmd := NewCommandHandler()
	query := NewQueryHandler()
	cmd.Subscribe(query.Apply)

	cmd.AddStock("widget", 10)
	cmd.AddStock("gadget", 2)
	cmd.RemoveStock("widget", 4)

	if got := query.GetStock("widget"); got != 6 {
		t.Fatalf("query.GetStock(widget) = %d, want 6", got)
	}
	if got := query.GetStock("gadget"); got != 2 {
		t.Fatalf("query.GetStock(gadget) = %d, want 2", got)
	}
}

func TestQueryHandler_UnwiredHandlerStaysAtZero(t *testing.T) {
	cmd := NewCommandHandler()
	query := NewQueryHandler() // deliberately NOT subscribed

	cmd.AddStock("widget", 10)

	if got := query.GetStock("widget"); got != 0 {
		t.Fatalf("query.GetStock(widget) for an unwired QueryHandler = %d, want 0 (no direct coupling to CommandHandler state)", got)
	}
}

// alternativeReadModel is a completely independent projector, defined only in the
// test file, with a different internal shape (a flat event log instead of an
// aggregated map). Subscribing it alongside QueryHandler to the SAME CommandHandler,
// with both working correctly and independently, proves the command side has no
// special coupling to any one read model.
type alternativeReadModel struct {
	log []StockEvent
}

func (r *alternativeReadModel) Apply(e StockEvent) {
	r.log = append(r.log, e)
}

func TestCQRS_MultipleIndependentReadModels(t *testing.T) {
	cmd := NewCommandHandler()
	query := NewQueryHandler()
	alt := &alternativeReadModel{}

	cmd.Subscribe(query.Apply)
	cmd.Subscribe(alt.Apply)

	cmd.AddStock("widget", 10)
	cmd.RemoveStock("widget", 2)

	if got := query.GetStock("widget"); got != 8 {
		t.Fatalf("query.GetStock(widget) = %d, want 8", got)
	}
	if len(alt.log) != 2 {
		t.Fatalf("alt.log has %d entries, want 2", len(alt.log))
	}
}
