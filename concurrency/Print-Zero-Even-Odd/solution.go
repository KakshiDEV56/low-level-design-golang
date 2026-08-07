package printzeroevenodd

/*
You have a function printNumber that can be called with an integer parameter and prints it to the console.

For example, calling printNumber(7) prints 7 to the console.
You are given an instance of the class ZeroEvenOdd that has three functions: zero, even, and odd. The same instance of ZeroEvenOdd will be passed to three different threads:

Thread A: calls zero() that should only output 0's.
Thread B: calls even() that should only output even numbers.
Thread C: calls odd() that should only output odd numbers.
Modify the given class to output the series "010203040506..." where the length of the series must be 2n.

Implement the ZeroEvenOdd class:

ZeroEvenOdd(int n) Initializes the object with the number n that represents the numbers that should be printed.
void zero(printNumber) Calls printNumber to output one zero.
void even(printNumber) Calls printNumber to output one even number.
void odd(printNumber) Calls printNumber to output one odd number.


Example 1:

Input: n = 2
Output: "0102"
Explanation: There are three threads being fired asynchronously.
One of them calls zero(), the other calls even(), and the last one calls odd().
"0102" is the correct output.
Example 2:

Input: n = 5
Output: "0102030405"


Constraints:

1 <= n <= 1000
*/

import "sync"

type ZeroEvenOdd struct {
	n   int
	wg  sync.WaitGroup
	ch1 chan struct{}
	ch2 chan struct{}
	ch3 chan struct{}
}

/*
Channel Closing Pattern — Close only after all senders are finished

A channel must be closed only when we are certain that NO goroutine
will ever send to it again.

Here, ch1/ch2/ch3 have multiple goroutines involved in their
synchronization, so an individual sender must NOT close them directly.

Instead, we use the WaitGroup to track the completion of all goroutines:

	Odd  ──────┐
	           │
	Even ──────┼──→ wg.Wait() ──→ close(channel)
	           │
	Zero ──────┘

The important pattern is:

 1. wg.Add(N) before starting/tracking N goroutines.
 2. Each goroutine calls wg.Done() when it finishes.
 3. A coordinator calls wg.Wait().
 4. Only after Wait() returns do we close the channel.

Why is this useful?

Closing a channel means:

	"No more values/signals will ever be sent on this channel."

Therefore, closing before all senders finish can cause:

	panic: send on closed channel

Waiting for all senders first guarantees that closing the channel is safe.

This pattern is especially useful for:
  - worker pools
  - pipelines
  - producer/consumer systems
  - graceful goroutine shutdown
  - broadcasting "work is completely finished"

IMPORTANT:
A WaitGroup and close() solve different problems:

	WaitGroup → tracks WHEN goroutines finish.
	close()    → tells receivers that NO MORE values will be sent.

Also remember:
A buffered channel's capacity does NOT mean it is initialized
with a value. close() also does NOT mean "destroy the channel";
it means "no more sends are allowed."

For synchronization/token channels like ch1/ch2/ch3 in this
ZeroEvenOdd problem, closing is generally unnecessary because
the loops already know exactly when to stop. A separate
"done" channel is usually a cleaner choice when you need to
broadcast that the entire operation has finished.
*/
func NewZeroEvenOdd(n int) *ZeroEvenOdd {
	zeo := &ZeroEvenOdd{
		n:   n,
		ch1: make(chan struct{}, 1),
		ch2: make(chan struct{}, 1),
		ch3: make(chan struct{}, 1),
	}
	zeo.wg.Add(3)
	return zeo
}

func (z *ZeroEvenOdd) Zero(printNumber func(int)) {
	defer z.wg.Done()
	for i := 1; i <= z.n; i++ {
		if i > 1 {
			_ = <-z.ch1
		}
		printNumber(0)

		if i%2 == 0 {
			z.ch2 <- struct{}{}
		} else if i%2 != 0 {
			z.ch3 <- struct{}{}
		}
	}
	go func() {
		z.wg.Wait()

		close(z.ch1)
		close(z.ch2)
		close(z.ch3)
	}()
}

func (z *ZeroEvenOdd) Even(printNumber func(int)) {
	defer z.wg.Done()
	for i := 2; i <= z.n; i += 2 {
		_ = <-z.ch2
		printNumber(i)
		z.ch1 <- struct{}{}
	}

}

func (z *ZeroEvenOdd) Odd(printNumber func(int)) {
	defer z.wg.Done()
	for i := 1; i <= z.n; i += 2 {
		_ = <-z.ch3
		printNumber(i)
		z.ch1 <- struct{}{}
	}

}
