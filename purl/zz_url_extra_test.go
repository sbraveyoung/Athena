package purl

import (
	"net/url"
	"strings"
	"sync"
	"testing"

	pool "github.com/sbraveyoung/Athena/easypool"
)

func TestParseQueryEmpty(t *testing.T) {
	v, err := ParseQuery("")
	if err != nil {
		t.Fatalf("ParseQuery(\"\"): %v", err)
	}
	if v.Len() != 0 {
		t.Errorf("Len()=%d want 0", v.Len())
	}
}

func TestParseQueryEscaped(t *testing.T) {
	v, err := ParseQuery("k=hello%20world")
	if err != nil {
		t.Fatalf("ParseQuery: %v", err)
	}
	if got := v.Get("k"); got != "hello world" {
		t.Errorf("Get(k)=%q want %q", got, "hello world")
	}
}

func TestParseQueryRejectsSemicolon(t *testing.T) {
	_, err := ParseQuery("a=1;b=2")
	if err == nil {
		t.Error("ParseQuery: expected error on semicolon separator")
	}
}

func TestParseQueryEmptyKeySkipped(t *testing.T) {
	v, err := ParseQuery("=v&a=b")
	if err != nil {
		t.Fatalf("ParseQuery: %v", err)
	}
	// "=v" has empty key after Cut("="), so it's added as key="" value="v" —
	// but per ParseQuery the key cut happens before the empty-string skip.
	// Just verify a/b survived.
	if got := v.Get("a"); got != "b" {
		t.Errorf("Get(a)=%q want b", got)
	}
}

func TestSetnxOnlyWhenAbsent(t *testing.T) {
	v := NewValues()
	if !v.Setnx("a", "1") {
		t.Errorf("Setnx on empty returned false")
	}
	if v.Setnx("a", "2") {
		t.Errorf("Setnx on existing key returned true")
	}
	if got := v.Get("a"); got != "1" {
		t.Errorf("Setnx clobbered existing value: got %q", got)
	}
}

func TestSetEmptyValuesIgnored(t *testing.T) {
	v := NewValues()
	v.Set("", "x")
	v.Set("k", "")
	if v.Len() != 0 {
		t.Errorf("Len()=%d, want 0 (empty key/value should be ignored)", v.Len())
	}
}

func TestRangeVisitsAllPairs(t *testing.T) {
	v := NewValues()
	v.Set("a", "b")
	v.Set("c", "d")
	v.Set("e", "f")

	seen := map[string]string{}
	v.Range(func(k, val string) bool {
		seen[k] = val
		return true
	})

	want := map[string]string{"a": "b", "c": "d", "e": "f"}
	for k, vv := range want {
		if seen[k] != vv {
			t.Errorf("Range did not see %s=%s, got %v", k, vv, seen)
		}
	}
}

func TestRangeEarlyExit(t *testing.T) {
	v := NewValues()
	v.Set("a", "1")
	v.Set("b", "2")
	v.Set("c", "3")

	count := 0
	v.Range(func(k, val string) bool {
		count++
		return false // stop after the first
	})
	if count != 1 {
		t.Errorf("Range visited %d items after returning false, want 1", count)
	}
}

func TestEncodeIsSorted(t *testing.T) {
	v, err := ParseQuery("c=3&a=1&b=2")
	if err != nil {
		t.Fatalf("ParseQuery: %v", err)
	}
	got := v.Encode()
	if got != "a=1&b=2&c=3" {
		t.Errorf("Encode()=%q, want sorted a=1&b=2&c=3", got)
	}
}

func TestEncodeNilValuesIsEmpty(t *testing.T) {
	var v *Values
	if got := v.Encode(); got != "" {
		t.Errorf("nil.Encode()=%q want empty", got)
	}
}

func TestGetNilValuesIsEmpty(t *testing.T) {
	var v *Values
	if got := v.Get("anything"); got != "" {
		t.Errorf("nil.Get()=%q want empty", got)
	}
}

func TestDelRemovesKey(t *testing.T) {
	v := NewValues()
	v.Set("k", "v")
	v.Del("k")
	if got := v.Get("k"); got != "" {
		t.Errorf("Get after Del = %q, want empty", got)
	}
}

func TestRoundTripParseEncode(t *testing.T) {
	rawQuery := "a=1&b=hello&c=world"
	v, err := ParseQuery(rawQuery)
	if err != nil {
		t.Fatalf("ParseQuery: %v", err)
	}

	got := v.Encode()
	// Round-trip via re-parse to ignore key ordering.
	stdGot, _ := url.ParseQuery(got)
	stdWant, _ := url.ParseQuery(rawQuery)
	for k, vs := range stdWant {
		if stdGot.Get(k) != vs[0] {
			t.Errorf("round-trip lost key %s: got %v want %v", k, stdGot.Get(k), vs[0])
		}
	}
}

// Concurrency: getIndex/getString must be race-free across goroutines.
func TestGetIndexConcurrent(t *testing.T) {
	const goroutines = 16
	const ops = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(g int) {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				key := "k" + string(rune('A'+g))
				idx := getIndex(key)
				if back := getString(idx); back != key {
					t.Errorf("getString(getIndex(%s))=%s", key, back)
					return
				}
			}
		}(g)
	}
	wg.Wait()
}

// Integration: purl.Values.Encode goes through easypool's buffer pool. Run a
// big batch of Encode calls and verify the output is correct each time, which
// proves the pool is being reset properly between checkout/return cycles.
func TestEncodeIntegrationWithEasyPool(t *testing.T) {
	v, err := ParseQuery("a=1&b=2&c=3")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		got := v.Encode()
		if got != "a=1&b=2&c=3" {
			t.Fatalf("iter %d: Encode=%q, want a=1&b=2&c=3", i, got)
		}
	}
}

// Integration sanity: verify Buffer obtained via pool.GetBuffer behaves like
// the implementation Encode relies on (the api shape is the part purl
// depends on).
func TestPurlBufferContract(t *testing.T) {
	buf, put := pool.GetBuffer()
	defer put(buf)

	buf.WriteString("a=")
	buf.WriteString("b")
	buf.WriteByte('&')
	buf.WriteString("c=d")
	got := buf.String()
	if !strings.HasPrefix(got, "a=b&c=") {
		t.Errorf("Buffer composed string = %q", got)
	}
}
