package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	mb "github.com/ocramh/fingerprinter/pkg/musicbrainz"
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
	Short: "Queries the MusicBrainz API and returns recordings and releases metadata associated with a recording ID",
	Run: func(cmd *cobra.Command, args []string) {
		mbClient := mb.NewMusicBrainz(appName, semVer, contactEmail)
		recInfo, err := mbClient.GetReleaseInfo(releaseID)
		if err != nil {
			log.Fatal(err)
		}

		b, err := json.Marshal(recInfo)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(os.Stdout, string(b))
	},
}
