package tasks

import (
	"fmt"

	"github.com/spf13/cobra"

	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
	meta "github.com/ocramh/fingerprinter/pkg/meta"
	vf "github.com/ocramh/fingerprinter/pkg/verifier"
)

var (
	audioPath string
)

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVarP(&apikey, "apikey", "k", "", "acoustid key")
	verifyCmd.Flags().StringVarP(&audioPath, "audiopath", "a", "", "audio file(s) path")
	verifyCmd.MarkFlagRequired("apikey")
	verifyCmd.MarkFlagRequired("audiopath")
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify inut audio metadata",
	Run: func(cmd *cobra.Command, args []string) {

		chromaMngr := &fp.ChromaIO{}
		acClient := meta.NewAcoustIDClient(apikey)
		mbClient := meta.NewMBClient("fingerprinter", "0.0.1", "marco@sygma.io")

		verifier := vf.NewAudioVerifier(chromaMngr, acClient, mbClient)
		res, err := verifier.Analyze(audioPath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", res)
	},
}
