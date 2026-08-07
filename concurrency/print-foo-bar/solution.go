package printfoobar

/*
Suppose you are given the following code:

class FooBar {
  public void foo() {
    for (int i = 0; i < n; i++) {
      print("foo");
    }
  }

  public void bar() {
    for (int i = 0; i < n; i++) {
      print("bar");
    }
  }
}
The same instance of FooBar will be passed to two different threads:

thread A will call foo(), while
thread B will call bar().
Modify the given program to output "foobar" n times.

*/

import "sync"

type FooBar struct {
	n     int
	turnA chan struct{}
	turnB chan struct{}
	wg    sync.WaitGroup
}

func NewFooBar(n int) *FooBar {

	return &FooBar{n: n,
		turnA: make(chan struct{}, 1),
		turnB: make(chan struct{}, 1),
	}
}

func (fb *FooBar) Foo(printFoo func()) {
	fb.wg.Add(2)
	defer fb.wg.Done()
	for i := 0; i < fb.n; i++ {
		fb.turnA <- struct{}{}
		// printFoo() outputs "foo". Do not change or remove this line.
		printFoo()
		fb.turnB <- struct{}{}

	}
}

func (fb *FooBar) Bar(printBar func()) {
	defer fb.wg.Done()
	for i := 0; i < fb.n; i++ {
		_ = <-fb.turnB
		// printBar() outputs "bar". Do not change or remove this line.
		printBar()
		_ = <-fb.turnA
	}

}
