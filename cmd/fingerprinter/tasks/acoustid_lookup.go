package tasks

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/ocramh/fingerprinter/pkg/clients"
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
	Short: "fetches the AcoustID ID and the MusicBrainz recording ID matching the fingerprint",
	Run: func(cmd *cobra.Command, args []string) {
		chroma := fp.ChromaIO{}
		fingerprints, err := chroma.CalcFingerprint(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		acoustIDClient := clients.NewAcoustID(apikey)
		retryOnFail := true

		for _, fingerprint := range fingerprints {
			resp, err := acoustIDClient.LookupFingerprint(fingerprint, retryOnFail)
			if err != nil {
				log.Fatal(err)
			}

			for _, r := range resp.Results {
				log.Printf("[score] %f \n", r.Score)
				log.Printf("[acoustid] %s \n", r.ID)

				for _, recording := range r.Recordings {
					log.Printf("[mb recording ID] %s \n", recording.MBRecordingID)

					for _, releaseGroup := range recording.MBReleaseGroups {
						for _, release := range releaseGroup.Releases {
							log.Printf("[mb release ID] %s \n", release.ID)
						}
					}
				}
			}

			time.Sleep(clients.AcoustIDReqDelay)
		}
	},
}
