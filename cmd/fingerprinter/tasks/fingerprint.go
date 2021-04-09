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
	fpCmd.Flags().StringVarP(&inputFile, "audiofile", "f", "", "audio file path")
	fpCmd.MarkFlagRequired("audiofile")
}

var fpCmd = &cobra.Command{
	Use:   "fpcalc",
	Short: "calculates the fingerprint of the input audio file",
	Run: func(cmd *cobra.Command, args []string) {
		chroma := fp.ChromaIO{}
		calc, err := chroma.CalcFingerprint(inputFile)
		if err != nil {
			panic(err)
		}

		log.Println(calc)
	},
}
