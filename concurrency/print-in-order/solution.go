package printinorder

import "sync"

/*
Suppose we have a class:

	public class Foo {
	  public void first() { print("first"); }
	  public void second() { print("second"); }
	  public void third() { print("third"); }
	}

The same instance of Foo will be passed to three different threads. Thread A will call first(), thread B will call second(), and thread C will call third(). Design a mechanism and modify the program to ensure that second() is executed after first(), and third() is executed after second().

Note:

We do not know how the threads will be scheduled in the operating system, even though the numbers in the input seem to imply the ordering. The input format you see is mainly to ensure our tests' comprehensiveness.

Example 1:

Input: nums = [1,2,3]
Output: "firstsecondthird"
Explanation: There are three threads being fired asynchronously. The input [1,2,3] means thread A calls first(), thread B calls second(), and thread C calls third(). "firstsecondthird" is the correct output.
Example 2:

Input: nums = [1,3,2]
Output: "firstsecondthird"
Explanation: The input [1,3,2] means thread A calls first(), thread B calls third(), and thread C calls second(). "firstsecondthird" is the correct output.
*/
type Foo struct {
	done1 chan struct{}
	done2 chan struct{}
	done3 chan struct{}
	wg    sync.WaitGroup
}

func NewFoo() *Foo {
	return &Foo{
		done1: make(chan struct{}),
		done2: make(chan struct{}),
		done3: make(chan struct{}),
	}
}

func (f *Foo) First(printFirst func()) {
	f.wg.Add(3)
	// Do not change this line
	defer f.wg.Done()
	printFirst()
	close(f.done1)

}

func (f *Foo) Second(printSecond func()) {
	/// Do not change this line
	<-f.done1
	defer f.wg.Done()
	printSecond()
	close(f.done2)
}
func (f *Foo) Third(printThird func()) {
	// Do not change this line
	<-f.done2
	defer f.wg.Done()
	printThird()
	close(f.done3)
}
