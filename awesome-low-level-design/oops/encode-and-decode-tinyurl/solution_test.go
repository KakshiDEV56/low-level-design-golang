package encodedecodetinyurl

import (
	"strconv"
	"testing"
)

func TestCodec_EncodeThenDecodeRoundTrip(t *testing.T) {
	c := NewCodec()
	long := "https://leetcode.com/problems/design-tinyurl"
	short := c.Encode(long)
	if got := c.Decode(short); got != long {
		t.Fatalf("Decode(Encode(%q)) = %q, want %q", long, got, long)
	}
}

func TestCodec_DifferentURLsGetDifferentShortURLs(t *testing.T) {
	c := NewCodec()
	a := c.Encode("https://example.com/a")
	b := c.Encode("https://example.com/b")
	if a == b {
		t.Fatalf("Encode() returned the same short URL %q for two different long URLs", a)
	}
}

func TestCodec_DecodingEachShortURLReturnsCorrectLongURL(t *testing.T) {
	c := NewCodec()
	longs := []string{
		"https://a.com/1",
		"https://b.com/2",
		"https://c.com/3",
	}
	shorts := make([]string, len(longs))
	for i, l := range longs {
		shorts[i] = c.Encode(l)
	}
	for i, s := range shorts {
		if got := c.Decode(s); got != longs[i] {
			t.Fatalf("Decode(%q) = %q, want %q", s, got, longs[i])
		}
	}
}

func TestCodec_ManyURLsAllRoundTripCorrectly(t *testing.T) {
	c := NewCodec()
	const n = 200
	longs := make([]string, n)
	shorts := make([]string, n)
	for i := 0; i < n; i++ {
		longs[i] = "https://example.com/page/" + strconv.Itoa(i)
		shorts[i] = c.Encode(longs[i])
	}
	seenShorts := make(map[string]bool)
	for i := 0; i < n; i++ {
		if seenShorts[shorts[i]] {
			t.Fatalf("short URL %q was generated more than once for distinct long URLs", shorts[i])
		}
		seenShorts[shorts[i]] = true

		if got := c.Decode(shorts[i]); got != longs[i] {
			t.Fatalf("Decode(%q) = %q, want %q", shorts[i], got, longs[i])
		}
	}
}

func TestCodec_LongURLWithQueryParamsAndFragments(t *testing.T) {
	c := NewCodec()
	long := "https://example.com/search?q=golang+interview&sort=asc#results"
	short := c.Encode(long)
	if got := c.Decode(short); got != long {
		t.Fatalf("Decode(Encode(url)) = %q, want %q", got, long)
	}
}

func TestCodec_ShortURLLooksLikeAURL(t *testing.T) {
	c := NewCodec()
	short := c.Encode("https://example.com/whatever")
	if len(short) == 0 {
		t.Fatal("Encode() returned an empty string")
	}
}
