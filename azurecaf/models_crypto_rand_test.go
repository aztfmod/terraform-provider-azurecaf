package azurecaf

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"
)

// failingReader always returns the provided error to simulate a broken OS
// entropy source (unreadable /dev/urandom on Linux, failing BCryptGenRandom on
// Windows). It is used to exercise the error-return branch of randomLetter and
// randSeq without depending on environment-specific behaviour.
func failingReader() func() {
	originalReader := cryptoRandReader
	cryptoRandReader = iotest.ErrReader(errors.New("simulated crypto/rand failure"))
	return func() { cryptoRandReader = originalReader }
}

// TestRandomLetterReturnsErrorWhenCryptoRandFails verifies the error-return
// branch added to randomLetter when crypto/rand.Reader fails. Without this
// test the branch is unreachable in normal CI runs and the function reports
// 50% coverage; with the injection it reports 100%.
func TestRandomLetterReturnsErrorWhenCryptoRandFails(t *testing.T) {
	restore := failingReader()
	defer restore()

	_, err := randomLetter()
	if err == nil {
		t.Fatal("expected randomLetter to return an error when crypto/rand reader fails")
	}
	if !strings.Contains(err.Error(), "crypto/rand.Reader failed") {
		t.Fatalf("expected wrapped crypto/rand failure message, got: %v", err)
	}
}

// TestRandSeqPropagatesRandomLetterError verifies that randSeq surfaces the
// underlying randomLetter error rather than swallowing it or panicking. This
// is the non-deterministic branch (seed == nil) where each rune is drawn from
// crypto/rand.
func TestRandSeqPropagatesRandomLetterError(t *testing.T) {
	restore := failingReader()
	defer restore()

	out, err := randSeq(16, nil)
	if err == nil {
		t.Fatal("expected randSeq to return an error when crypto/rand reader fails")
	}
	if out != "" {
		t.Fatalf("expected empty string on error, got %q", out)
	}
	if !strings.Contains(err.Error(), "crypto/rand.Reader failed") {
		t.Fatalf("expected wrapped crypto/rand failure message, got: %v", err)
	}
}

// TestRandSeqDeterministicPathIgnoresCryptoReader verifies the seeded branch
// continues to operate even when cryptoRandReader is broken, because it draws
// from math/rand instead. This protects the issue #336 plan-time visibility
// contract from regressions when the crypto path is unavailable.
func TestRandSeqDeterministicPathIgnoresCryptoReader(t *testing.T) {
	restore := failingReader()
	defer restore()

	seed := int64(42)
	out, err := randSeq(12, &seed)
	if err != nil {
		t.Fatalf("seeded randSeq must not consult crypto/rand: %v", err)
	}
	if len(out) != 12 {
		t.Fatalf("expected length 12, got %d (%q)", len(out), out)
	}
}
