package cigarettesmokers

/*
Implement the cigarette smokers problem using the "pusher" / dedicated-coordinator
pattern. See README.md for the full explanation — this is The Little Book of
Semaphores, §4.5.

TODO:
 1. Add the fields you need on Table below. Hint:
      - three booleans: tobaccoOnTable, paperOnTable, matchOnTable
      - a sync.Mutex protecting them
      - three "trigger" channels (buffered size 1 is enough), one per smoker: e.g.
        matchesSmokerCh, tobaccoSmokerCh, paperSmokerCh — each is how a Provide* call
        tells the corresponding Smoke* call "your pair is ready."
 2. Implement NewTable: initialize the flags to false and the channels.
 3. Implement ProvideTobacco(), ProvidePaper(), ProvideMatch() using the "pusher" logic
    described in the README: lock the mutex; check whether the OTHER ingredient needed
    to complete a pair with what just arrived is already on the table. If so, clear
    that flag and send a value on the corresponding smoker's trigger channel (this is
    what wakes them). If not, set this ingredient's own flag to true (so the NEXT
    Provide* call can find it). Unlock the mutex.
      - ProvideTobacco: pairs with paper (wakes SmokeWithMatches) or with match (wakes
        SmokeWithPaper).
      - ProvidePaper: pairs with tobacco (wakes SmokeWithMatches) or with match (wakes
        SmokeWithTobacco).
      - ProvideMatch: pairs with tobacco (wakes SmokeWithPaper) or with paper (wakes
        SmokeWithTobacco).
 4. Implement SmokeWithMatches(smoke func()), SmokeWithTobacco(smoke func()),
    SmokeWithPaper(smoke func()): each just receives from its own trigger channel
    (blocking until a Provide* call signals it), then calls smoke() once, then returns.
*/

type Table struct {
	// TODO: add your fields here
}

func NewTable() *Table {
	// TODO: implement
	return &Table{}
}

func (t *Table) ProvideTobacco() {
	// TODO: implement
}

func (t *Table) ProvidePaper() {
	// TODO: implement
}

func (t *Table) ProvideMatch() {
	// TODO: implement
}

func (t *Table) SmokeWithMatches(smoke func()) {
	// TODO: implement
}

func (t *Table) SmokeWithTobacco(smoke func()) {
	// TODO: implement
}

func (t *Table) SmokeWithPaper(smoke func()) {
	// TODO: implement
}
