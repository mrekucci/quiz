package main

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
	"io/ioutil"
)

// randStr returns a byte slice of characters of
// size n built randomly from characters in t.
func randStr(n int, t string, s rand.Source) []byte {
	l := len(t) - 1
	buf := make([]byte, n)
	if l > 0 {
		for i, p := range rand.New(s).Perm(n) {
			buf[i] = t[p%l]
		}
	}
	return buf
}

func TestQuiz(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
	}{
		{"", "\n"},
		{"\n", "\n"},
		{"\n\n", "\n"},
		{"A", "\n"},
		{"A\n", "\n"},
		{"A\n" + "BC", "\n"},
		{"ABC\n" + "AB\n" + "C", "ABC\n"},
		{"ABC\n" + "AB\n" + "C\n" + "XYZW", "ABC\n"},
		{"☺世☺\n" + "☺世\n" + "☺", "☺世☺\n"},
		{"CAT\n" +
			"CATS\n" +
			"CATSDOGCATS\n" +
			"CATXDOGCATSRAT\n" +
			"DOG\n" +
			"DOGCATSDOG\n" +
			"HIPPOPOTAMUSES\n" +
			"RAT\n" +
			"RATCAT\n" +
			"RATCATDOG\n" +
			"RATCATDOGCAT", "RATCATDOGCAT\n"},
	} {
		w := new(bytes.Buffer)
		err := do(w, strings.NewReader(test.in))
		if err != nil {
			t.Fatalf("do: %q: got error %s\n", test.in, err)
		}
		if w.String() != test.want {
			t.Errorf("Quiz input: %q\ngot:  %q\nwant: %q", test.in, w, test.want)
		}
	}
}

func benchQuiz(b *testing.B, size int) {
	b.StopTimer()
	size = size * 10 // Every word will be 10 characters long.
	str := randStr(size, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", rand.NewSource(int64(size)))
	var in []byte
	for i := 10; i <= len(str); i += 10 { // Cut the generated string by length 10.
		in = append(in, str[i-10:i]...)
		in = append(in, '\n')
	}
	w, r := ioutil.Discard, bytes.NewReader(in)
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		do(w, r)
		b.StopTimer()
		r.Seek(0, 0)
	}
}

func BenchmarkQuiz1e4Words(b *testing.B) { benchQuiz(b, 1e4) }
func BenchmarkQuiz1e5Words(b *testing.B) { benchQuiz(b, 1e5) }
func BenchmarkQuiz1e6Words(b *testing.B) { benchQuiz(b, 1e6) }
