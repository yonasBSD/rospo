package cmd

import (
	"path/filepath"

	"github.com/ferama/rospo/pkg/autocomplete"
	"github.com/ferama/rospo/pkg/sshc"
	"github.com/ferama/rospo/pkg/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(grabpubkeyCmd)

	usr := utils.CurrentUser()
	knownHostFile := filepath.Join(usr.HomeDir, ".ssh", "known_hosts")
	grabpubkeyCmd.PersistentFlags().StringP("known-hosts", "k", knownHostFile, "the known_hosts file absolute path")
}

var grabpubkeyCmd = &cobra.Command{
	Use:   "grabpubkey host:port",
	Short: "Grab the host pubkey and put it into the known_hosts file",
	Long:  `Grab the host pubkey and put it into the known_hosts file`,
	Example: `
 # grabs the pubkey from the server at host:port and put it into ./known file
 $ rospo grabpubkey -k ./known host:port
	`,
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: autocomplete.Host(),
	Run: func(cmd *cobra.Command, args []string) {
		knownHosts, _ := cmd.Flags().GetString("known-hosts")
		sshcConf := &sshc.SshClientConf{
			KnownHosts: knownHosts,
			ServerURI:  args[0],
		}
		client := sshc.NewSshConnection(sshcConf)
		client.GrabPubKey()
	},
}
