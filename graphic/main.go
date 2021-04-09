package graphic

// import (
// 	"flag"
// 	"github.com/go-fingerprint/fingerprint"
// 	"github.com/go-fingerprint/gochroma"
// 	"image/png"
// 	"log"
// 	"os"
// )

// func main() {
// 	flag.Parse()

// 	if flag.NArg() < 1 {
// 		println("Usage: graphic <file>")
// 	}

// 	f, err := os.Open(flag.Arg(0))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create new fingerprint calculator
// 	fpcalc := gochroma.New(gochroma.AlgorithmDefault)
// 	defer fpcalc.Close()

// 	fprint, err := fpcalc.RawFingerprint(
// 		fingerprint.RawInfo{
// 			Src:        f,
// 			Channels:   2,
// 			Rate:       44100,
// 			MaxSeconds: 120,
// 		})

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Get image representation of our fingerprint
// 	i := fingerprint.ToImage(fprint)

// 	out, err := os.Create(flag.Arg(0) + ".png")

// 	if err := png.Encode(out, i); err != nil {
// 		log.Fatal(err)
// 	}
// }
