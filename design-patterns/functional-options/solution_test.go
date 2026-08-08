package functionaloptions

import (
	"testing"
	"time"
)

func TestNewServer_Defaults(t *testing.T) {
	s := NewServer("localhost")
	if s.Host != "localhost" {
		t.Errorf("Host = %q, want %q", s.Host, "localhost")
	}
	if s.Port != 8080 {
		t.Errorf("Port = %d, want 8080", s.Port)
	}
	if s.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", s.Timeout)
	}
	if s.Retries != 3 {
		t.Errorf("Retries = %d, want 3", s.Retries)
	}
	if s.TLS != false {
		t.Errorf("TLS = %v, want false", s.TLS)
	}
}

func TestWithPort(t *testing.T) {
	s := NewServer("localhost", WithPort(9090))
	if s.Port != 9090 {
		t.Fatalf("Port = %d, want 9090", s.Port)
	}
}

func TestWithTimeout(t *testing.T) {
	s := NewServer("localhost", WithTimeout(5*time.Second))
	if s.Timeout != 5*time.Second {
		t.Fatalf("Timeout = %v, want 5s", s.Timeout)
	}
}

func TestWithRetries(t *testing.T) {
	s := NewServer("localhost", WithRetries(10))
	if s.Retries != 10 {
		t.Fatalf("Retries = %d, want 10", s.Retries)
	}
}

func TestWithTLS(t *testing.T) {
	s := NewServer("localhost", WithTLS())
	if !s.TLS {
		t.Fatal("TLS = false, want true")
	}
}

func TestOptions_UnspecifiedFieldsKeepDefaults(t *testing.T) {
	s := NewServer("localhost", WithPort(443))
	if s.Timeout != 30*time.Second || s.Retries != 3 || s.TLS != false {
		t.Fatalf("unrelated fields changed: Timeout=%v Retries=%d TLS=%v, want defaults", s.Timeout, s.Retries, s.TLS)
	}
}

func TestOptions_LaterOptionForSameFieldWins(t *testing.T) {
	s := NewServer("localhost", WithPort(80), WithPort(443))
	if s.Port != 443 {
		t.Fatalf("Port = %d, want 443 (later option should win)", s.Port)
	}
}

func TestOptions_MultipleCombined(t *testing.T) {
	s := NewServer("example.com", WithPort(443), WithTLS(), WithRetries(5), WithTimeout(time.Minute))
	if s.Host != "example.com" || s.Port != 443 || !s.TLS || s.Retries != 5 || s.Timeout != time.Minute {
		t.Fatalf("Server = %+v, want Host=example.com Port=443 TLS=true Retries=5 Timeout=1m", s)
	}
}

// withProductionDefaults is a NEW option, defined only in the test file, that sets
// multiple fields at once by composing existing options. If this works with
// NewServer unmodified, it proves Option is a genuinely open type NewServer doesn't
// special-case internally.
func withProductionDefaults() Option {
	return func(s *Server) {
		WithPort(443)(s)
		WithTLS()(s)
		WithRetries(5)(s)
	}
}

func TestCustomOption_DefinedOutsideSolutionGo(t *testing.T) {
	s := NewServer("prod.example.com", withProductionDefaults())
	if s.Port != 443 || !s.TLS || s.Retries != 5 {
		t.Fatalf("Server = %+v, want Port=443 TLS=true Retries=5", s)
	}
	// Default (untouched by the custom option) should remain.
	if s.Timeout != 30*time.Second {
		t.Fatalf("Timeout = %v, want unchanged default 30s", s.Timeout)
	}
}

func TestCustomOption_CanBeOverriddenByLaterBuiltinOption(t *testing.T) {
	s := NewServer("prod.example.com", withProductionDefaults(), WithPort(8443))
	if s.Port != 8443 {
		t.Fatalf("Port = %d, want 8443 (explicit later option should override the custom option's port)", s.Port)
	}
}

func TestNewServer_NoOptionsIsEquivalentToAllDefaults(t *testing.T) {
	a := NewServer("host-a")
	b := NewServer("host-a")
	if *a != *b {
		t.Fatalf("two NewServer() calls with identical (no) options produced different results: %+v vs %+v", a, b)
	}
}
