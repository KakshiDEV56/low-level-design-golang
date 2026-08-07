package movingaveragefromdatastream

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-4
}

func TestMovingAverage_Example1(t *testing.T) {
	m := NewMovingAverage(3)
	if got := m.Next(1); !almostEqual(got, 1.0) {
		t.Fatalf("Next(1) = %v, want 1.0", got)
	}
	if got := m.Next(10); !almostEqual(got, 5.5) {
		t.Fatalf("Next(10) = %v, want 5.5", got)
	}
	if got := m.Next(3); !almostEqual(got, 4.66667) {
		t.Fatalf("Next(3) = %v, want 4.66667", got)
	}
	if got := m.Next(5); !almostEqual(got, 6.0) {
		t.Fatalf("Next(5) = %v, want 6.0", got)
	}
}

func TestMovingAverage_WindowSizeOne(t *testing.T) {
	m := NewMovingAverage(1)
	if got := m.Next(5); !almostEqual(got, 5.0) {
		t.Fatalf("Next(5) = %v, want 5.0", got)
	}
	if got := m.Next(10); !almostEqual(got, 10.0) {
		t.Fatalf("Next(10) = %v, want 10.0", got)
	}
	if got := m.Next(-3); !almostEqual(got, -3.0) {
		t.Fatalf("Next(-3) = %v, want -3.0", got)
	}
}

func TestMovingAverage_NegativeValues(t *testing.T) {
	m := NewMovingAverage(2)
	if got := m.Next(-4); !almostEqual(got, -4.0) {
		t.Fatalf("Next(-4) = %v, want -4.0", got)
	}
	if got := m.Next(-6); !almostEqual(got, -5.0) {
		t.Fatalf("Next(-6) = %v, want -5.0", got)
	}
	if got := m.Next(10); !almostEqual(got, 2.0) {
		t.Fatalf("Next(10) = %v, want 2.0", got)
	}
}

func TestMovingAverage_ExactWindowFillThenSlide(t *testing.T) {
	m := NewMovingAverage(4)
	vals := []int{1, 2, 3, 4}
	want := []float64{1.0, 1.5, 2.0, 2.5}
	for i, v := range vals {
		if got := m.Next(v); !almostEqual(got, want[i]) {
			t.Fatalf("Next(%d) = %v, want %v", v, got, want[i])
		}
	}
	// Window is now full [1,2,3,4]; next value slides out the 1.
	if got := m.Next(8); !almostEqual(got, 4.25) {
		t.Fatalf("Next(8) = %v, want 4.25", got)
	}
}

func TestMovingAverage_AllSameValue(t *testing.T) {
	m := NewMovingAverage(5)
	for i := 0; i < 10; i++ {
		if got := m.Next(7); !almostEqual(got, 7.0) {
			t.Fatalf("Next(7) call %d = %v, want 7.0", i, got)
		}
	}
}

func TestMovingAverage_LargeWindowFewValues(t *testing.T) {
	m := NewMovingAverage(1000)
	if got := m.Next(2); !almostEqual(got, 2.0) {
		t.Fatalf("Next(2) = %v, want 2.0", got)
	}
	if got := m.Next(4); !almostEqual(got, 3.0) {
		t.Fatalf("Next(4) = %v, want 3.0", got)
	}
}

func TestMovingAverage_MixedPositiveNegativeAverageToZero(t *testing.T) {
	m := NewMovingAverage(2)
	m.Next(5)
	if got := m.Next(-5); !almostEqual(got, 0.0) {
		t.Fatalf("Next(-5) = %v, want 0.0", got)
	}
}

func TestMovingAverage_StressSlidingWindowCorrectness(t *testing.T) {
	const size = 5
	m := NewMovingAverage(size)
	var window []int

	for v := 1; v <= 200; v++ {
		got := m.Next(v)

		window = append(window, v)
		if len(window) > size {
			window = window[1:]
		}
		sum := 0
		for _, w := range window {
			sum += w
		}
		want := float64(sum) / float64(len(window))

		if !almostEqual(got, want) {
			t.Fatalf("Next(%d) = %v, want %v", v, got, want)
		}
	}
}
