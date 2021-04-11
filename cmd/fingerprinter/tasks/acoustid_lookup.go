package tasks

import (
	"log"

	"github.com/spf13/cobra"

	fp "github.com/ocramh/fingerprinter/pkg/fingerprint"
	"github.com/ocramh/fingerprinter/pkg/meta"
)

var (
	fingerprint string
	duration    float32
	apikey      string
)

func init() {
	rootCmd.AddCommand(acoustidCmd)
	fpCmd.Flags().StringVarP(&fingerprint, "fingerp", "f", "", "audio fingerprint")
	fpCmd.Flags().Float32VarP(&duration, "duration", "d", 0, "audio duration")
	fpCmd.Flags().StringVarP(&apikey, "apikey", "k", "", "acustid key")
	fpCmd.MarkFlagRequired("fingerp")
	fpCmd.MarkFlagRequired("duration")
	fpCmd.MarkFlagRequired("apikey")
}

var acoustidCmd = &cobra.Command{
	Use:   "acoustid",
	Short: "fetches the AcoustID ID and the MusicBrainz recording ID matching the fingerprint",
	Run: func(cmd *cobra.Command, args []string) {
		acoustIDClient := meta.NewAcoustIDClient(apikey)
		resp, err := acoustIDClient.LookupFingerprint(&fp.Fingerprint{
			Duration: duration,
			Value:    fingerprint,
		})

		if err != nil {
			panic(err)
		}

		log.Printf("[status] %s \n", resp.Status)
		for _, r := range resp.Results {
			log.Printf("[score] %f \n", r.Score)
			log.Printf("[acustid] %s \n", r.ID)

			for _, recording := range r.Recordings {
				log.Printf("[mb recording ID] %f \n", recording.MBRecordingsID)

				for _, release := range recording.MBReleaseGroupsID {
					log.Printf("[mb release ID] %f \n", release.ID)
				}
			}
		}
	},
}
