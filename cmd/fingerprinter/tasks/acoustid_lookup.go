package tasks

import (
	"log"

	"github.com/spf13/cobra"

	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
	"github.com/ocramh/fingerprinter/pkg/meta"
)

var (
	apikey string
)

func init() {
	rootCmd.AddCommand(acoustidCmd)
	acoustidCmd.Flags().StringVarP(&apikey, "apikey", "k", "", "acustid key")
	acoustidCmd.Flags().StringVarP(&inputFile, "audiofile", "a", "", "audio file path")
	acoustidCmd.MarkFlagRequired("apikey")
	acoustidCmd.MarkFlagRequired("audiofile")
}

var acoustidCmd = &cobra.Command{
	Use:   "acoustid",
	Short: "fetches the AcoustID ID and the MusicBrainz recording ID matching the fingerprint",
	Run: func(cmd *cobra.Command, args []string) {
		chroma := fp.ChromaIO{}
		fingerp, err := chroma.CalcFingerprint(inputFile)
		if err != nil {
			panic(err)
		}

		acoustIDClient := meta.NewAcoustIDClient(apikey)
		resp, err := acoustIDClient.LookupFingerprint(fingerp)

		if err != nil {
			panic(err)
		}

		log.Printf("[status] %s \n", resp.Status)
		for _, r := range resp.Results {
			log.Printf("[score] %f \n", r.Score)
			log.Printf("[acustid] %s \n", r.ID)

			for _, recording := range r.Recordings {
				log.Printf("[mb recording ID] %s \n", recording.MBRecordingsID)

				for _, release := range recording.MBReleaseGroupsID {
					log.Printf("[mb release ID] %s \n", release.ID)
				}
			}
		}
	},
}
