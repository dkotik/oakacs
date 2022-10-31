//+build ignore
//go:build ignore

package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"text/scanner"
	"time"

	"github.com/spf13/pflag"
)

var (
	source      = pflag.String("source", "", "file to load")
	destination = pflag.String("destination", "", "file to save to")
	variable    = pflag.String("variable", "", "variable name")
)

func createDictionary(destination, source, variable string) (err error) {
	words := make([]string, 256)

	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()
	s := &scanner.Scanner{}
	s.Init(f)
	s.Error = func(s *scanner.Scanner, msg string) {
		err = errors.New(msg)
	}

	for i := 0; i < 256; i++ {
		s.Scan()
		if err != nil {
			return err
		}
		words[i] = s.TokenText()
	}

	// b, err := os.ReadFile(source)
	// if err != nil {
	// 	return err
	// }
	// words := bytes.Fields(b)[:256]
	rand.Shuffle(256, func(i int, j int) {
		words[i], words[j] = words[j], words[i]
	})

	out, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.WriteString(out, "package oakwords\n\n// Autogenerated file from "+source+"\n\nvar "+variable+" = Dictionary{\n")
	if err != nil {
		return err
	}
	for _, word := range words {
		if _, err = fmt.Fprintf(out, "\t%q,\n", word); err != nil {
			return err
		}
	}
	_, err = io.WriteString(out, "}\n")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	pflag.Parse()
	_, err := os.Stat(*destination)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("[INFO] file %s already exists, skipping\n", *destination)
		return
	}
	rand.Seed(time.Now().UnixNano())
	if err := createDictionary(*destination, *source, *variable); err != nil {
		panic(err)
	}
}