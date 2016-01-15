package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("quiz: ")

	flagSet := new(flag.FlagSet)
	flagSet.Usage = usage
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if flagSet.NArg() != 1 {
		flagSet.Usage()
	}

	file, err := os.Open(flagSet.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err = do(os.Stdout, file); err != nil {
		log.Fatal(err)
	}
}

// do is the workhorse here, isolated from main to make testing easier.
func do(w io.Writer, r io.Reader) error {
	var lengths []int
	words := make(map[int]map[string]bool)

	// We suppose here that a line is no longer than 65536 characters (limitation of the scanner).
	// If that is an issue, then switch to more low-level procedure.
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		word := scanner.Text()
		length := len(word)
		if _, ok := words[length]; !ok {
			words[length] = make(map[string]bool)
			lengths = append(lengths, length)
		}
		words[length][word] = true
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	compound := findLongestCompoundWord(lengths, words)
	_, err := fmt.Fprintln(w, compound)
	return err
}

// usage prints program usage to stderr and exits the program.
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: quiz <file> a file name with a list of words, where each word is separated by newline\n\n")
	flag.PrintDefaults()
	os.Exit(2)
}

// findLongestCompoundWord iterates through words by its lengths and searches
// for the longest compound-word in the words, which is also a concatenation
// of other sub-words that exist in the words. The search process takes
// advantage by iterating through the map from longest to shortest words, which
// can speed-up the searches for some cases.
func findLongestCompoundWord(lengths []int, words map[int]map[string]bool) string {
	sort.Sort(sort.Reverse(sort.IntSlice(lengths)))
	for _, l := range lengths {
		for w := range words[l] {
			if isWordCompound(l, w, words) {
				return w
			}
		}
	}
	return ""
}

// isWordCompound recursively splits w into two parts according to the l and
// checks if both of the created sub-strings are presented in the words.
func isWordCompound(l int, w string, words map[int]map[string]bool) bool {
	if l < len(w) {
		if index, ok := words[len(w[l:])]; ok && index[w[l:]] {
			return true
		}
	}

	l--
	for i := l; i > 0; i-- {
		if index, ok := words[len(w[:i])]; ok && index[w[:i]] {
			if isWordCompound(i, w, words) {
				return true
			}
		}
	}

	return false
}
