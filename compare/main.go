package main

import (
	"flag"
	"fmt"
	"github.com/go-fingerprint/fingerprint"
	"github.com/go-fingerprint/gochroma"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

var (
	file1, file2 string
)

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		println("Usage: compare <file1> <file2>")
		os.Exit(0)
	}

	f1, err := os.Open(flag.Arg(0))

	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(flag.Arg(1))

	if err != nil {
		log.Fatal(err)
	}

	// Create new fingerprint calculator
	fpcalc := gochroma.New(gochroma.AlgorithmDefault)
	defer fpcalc.Close()

	// Get fingerprints as a slices of 32-bit integers
	fprint1, err := fpcalc.RawFingerprint(
		fingerprint.RawInfo{
			Src:        f1,
			Channels:   2,
			Rate:       44100,
			MaxSeconds: 120,
		})

	if err != nil {
		log.Fatal(err)
	}

	fprint2, err := fpcalc.RawFingerprint(
		fingerprint.RawInfo{
			Src:        f2,
			Channels:   2,
			Rate:       44100,
			MaxSeconds: 120,
		})

	if err != nil {
		log.Fatal(err)
	}

	// Compare fingerprints
	s, err := fingerprint.Compare(fprint1, fprint2)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Score: %v\n", s)

	if s > 0.95 {
		fmt.Println("Fingerprints do not differ a lot, maybe it's one track")
	} else {
		fmt.Println(`Fingerprints differ a lot, it's definitely different
			records`)
	}

	// Get graphical representation of distance between fingerprints
	i, err := fingerprint.ImageDistance(fprint1, fprint2)

	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(filepath.Join(filepath.Dir(flag.Arg(0)), "out.png"))

	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(out, i); err != nil {
		log.Fatal(nil)
	}
}
