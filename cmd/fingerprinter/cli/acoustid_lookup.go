package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	ac "github.com/ocramh/fingerprinter/pkg/acoustid"
	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
)

var (
	apikey string
)

func init() {
	rootCmd.AddCommand(acoustidCmd)
	acoustidCmd.Flags().StringVarP(&apikey, "apikey", "k", "", "acoustid key")
	acoustidCmd.Flags().StringVarP(&inputFile, "audiofile", "a", "", "audio file path")
	acoustidCmd.MarkFlagRequired("apikey")
	acoustidCmd.MarkFlagRequired("audiofile")
}

var acoustidCmd = &cobra.Command{
	Use:   "acoustid",
	Short: "Generate an audio fingerprint and queries the AcoustID API to find matching recording ID(s)",
	Run: func(cmd *cobra.Command, args []string) {
		chroma := fp.NewChromaPrint(exec.Command, afero.NewOsFs())
		fingerprints, err := chroma.CalcFingerprint(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		acoustIDClient := ac.NewAcoustID(apikey)
		retryOnFail := true

		var lookupRes []ac.ACLookupResult
		for _, fingerprint := range fingerprints {
			resp, err := acoustIDClient.LookupFingerprint(fingerprint, retryOnFail)
			if err != nil {
				log.Fatal(err)
			}

			lookupRes = append(lookupRes, resp.Results...)

			time.Sleep(ac.AcoustIDReqDelay)
		}

		b, err := json.Marshal(lookupRes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(os.Stdout, string(b))
	},
}
