package designparkingsystem

import "testing"

func TestParkingSystem_Example1(t *testing.T) {
	p := NewParkingSystem(1, 1, 0)
	if got := p.AddCar(1); !got {
		t.Fatal("AddCar(1) = false, want true")
	}
	if got := p.AddCar(2); !got {
		t.Fatal("AddCar(2) = false, want true")
	}
	if got := p.AddCar(3); got {
		t.Fatal("AddCar(3) = true, want false (no small slots)")
	}
	if got := p.AddCar(1); got {
		t.Fatal("AddCar(1) second time = true, want false (big slot taken)")
	}
}

func TestParkingSystem_ZeroSlotsForAllSizes(t *testing.T) {
	p := NewParkingSystem(0, 0, 0)
	if got := p.AddCar(1); got {
		t.Fatal("AddCar(1) with 0 big slots = true, want false")
	}
	if got := p.AddCar(2); got {
		t.Fatal("AddCar(2) with 0 medium slots = true, want false")
	}
	if got := p.AddCar(3); got {
		t.Fatal("AddCar(3) with 0 small slots = true, want false")
	}
}

func TestParkingSystem_MultipleSlotsSameSize(t *testing.T) {
	p := NewParkingSystem(0, 0, 3)
	for i := 0; i < 3; i++ {
		if got := p.AddCar(3); !got {
			t.Fatalf("AddCar(3) call %d = false, want true", i)
		}
	}
	if got := p.AddCar(3); got {
		t.Fatal("AddCar(3) 4th call = true, want false (only 3 slots)")
	}
}

func TestParkingSystem_IndependentSizeCounters(t *testing.T) {
	p := NewParkingSystem(2, 1, 1)
	p.AddCar(1)
	if got := p.AddCar(1); !got {
		t.Fatal("AddCar(1) second big car = false, want true (2 big slots)")
	}
	if got := p.AddCar(1); got {
		t.Fatal("AddCar(1) third big car = true, want false")
	}
	// Medium and small should be unaffected by big slots being exhausted.
	if got := p.AddCar(2); !got {
		t.Fatal("AddCar(2) = false, want true")
	}
	if got := p.AddCar(3); !got {
		t.Fatal("AddCar(3) = false, want true")
	}
}

func TestParkingSystem_LargeCapacityStress(t *testing.T) {
	p := NewParkingSystem(1000, 1000, 1000)
	for i := 0; i < 1000; i++ {
		if got := p.AddCar(1); !got {
			t.Fatalf("AddCar(1) call %d = false, want true", i)
		}
		if got := p.AddCar(2); !got {
			t.Fatalf("AddCar(2) call %d = false, want true", i)
		}
		if got := p.AddCar(3); !got {
			t.Fatalf("AddCar(3) call %d = false, want true", i)
		}
	}
	if got := p.AddCar(1); got {
		t.Fatal("AddCar(1) after exhausting 1000 slots = true, want false")
	}
	if got := p.AddCar(2); got {
		t.Fatal("AddCar(2) after exhausting 1000 slots = true, want false")
	}
	if got := p.AddCar(3); got {
		t.Fatal("AddCar(3) after exhausting 1000 slots = true, want false")
	}
}

func TestParkingSystem_OnlyMediumAvailable(t *testing.T) {
	p := NewParkingSystem(0, 5, 0)
	if got := p.AddCar(1); got {
		t.Fatal("AddCar(1) = true, want false")
	}
	if got := p.AddCar(3); got {
		t.Fatal("AddCar(3) = true, want false")
	}
	for i := 0; i < 5; i++ {
		if got := p.AddCar(2); !got {
			t.Fatalf("AddCar(2) call %d = false, want true", i)
		}
	}
	if got := p.AddCar(2); got {
		t.Fatal("AddCar(2) 6th call = true, want false")
	}
}
