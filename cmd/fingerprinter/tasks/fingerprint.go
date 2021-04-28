package tasks

import (
	"log"

	"github.com/spf13/cobra"

	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
)

var (
	inputFile string
)

func init() {
	rootCmd.AddCommand(fpCmd)
	fpCmd.Flags().StringVarP(&inputFile, "audiofile", "a", "", "path to input audio file or directory")
	fpCmd.MarkFlagRequired("audiofile")
}

var fpCmd = &cobra.Command{
	Use:   "fpcalc",
	Short: "calculates the fingerprint of the input audio file",
	Run: func(cmd *cobra.Command, args []string) {
		chroma := fp.ChromaIO{}
		fingerprints, err := chroma.CalcFingerprint(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		for _, fingerprint := range fingerprints {
			log.Printf("[duration] %f\n", fingerprint.Duration)
			log.Printf("[fingerprint] %s\n", fingerprint.Value)
		}
	},
}
