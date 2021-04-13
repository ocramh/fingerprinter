package tasks

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ocramh/fingerprinter/pkg/meta"
)

var (
	appName      string
	semVer       string
	contactEmail string
	releaseID    string
)

func init() {
	rootCmd.AddCommand(mbCmd)
	mbCmd.Flags().StringVarP(&appName, "appname", "n", "fingerprinter", "the name of the application")
	mbCmd.Flags().StringVarP(&semVer, "semver", "s", "0.0.1", "the application semantic version")
	mbCmd.Flags().StringVarP(&contactEmail, "email", "e", "", "contact email address")
	mbCmd.Flags().StringVarP(&releaseID, "release", "r", "", "the release ID to lookup")
	mbCmd.MarkFlagRequired("email")
}

var mbCmd = &cobra.Command{
	Use:   "mblookup",
	Short: "queries the MusicBrainz API and returns recordings and releases metadata",
	Run: func(cmd *cobra.Command, args []string) {
		mbClient := meta.NewMBClient(appName, appName, contactEmail)
		recInfo, err := mbClient.GetReleaseInfo(releaseID)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[recording title]: %s \n", recInfo.Title)

		for _, md := range recInfo.Media {
			for _, tr := range md.Tracks {
				log.Printf("[recording id]: %s \n", tr.Recording.ID)
				log.Printf("[track title]: %s \n", tr.Title)
				log.Printf("[track id]: %s \n", tr.ID)
				log.Printf("[track position]: %d \n", tr.Position)
				log.Printf("[track isrc]: %v \n", tr.Recording.ISRCs)
			}
		}
	},
}
