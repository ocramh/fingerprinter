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
	verifyCmd.Flags().StringVarP(&appName, "appname", "n", "fingerprinter", "the name of the application")
	verifyCmd.Flags().StringVarP(&semVer, "semver", "s", "0.0.1", "the application semantic version")
	verifyCmd.Flags().StringVarP(&contactEmail, "email", "e", "", "contact email address")
	verifyCmd.MarkFlagRequired("apikey")
	verifyCmd.MarkFlagRequired("audiopath")
	verifyCmd.MarkFlagRequired("email")
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "verify inut audio metadata",
	Run: func(cmd *cobra.Command, args []string) {

		chromaMngr := &fp.ChromaIO{}
		acClient := meta.NewAcoustIDClient(apikey)
		mbClient := meta.NewMBClient(appName, semVer, contactEmail)

		verifier := vf.NewAudioVerifier(chromaMngr, acClient, mbClient)
		res, err := verifier.Analyze(audioPath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", res)
	},
}
