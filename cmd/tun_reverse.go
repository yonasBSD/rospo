package cmd

import (
	"github.com/ferama/rospo/conf"
	"github.com/ferama/rospo/sshc"
	"github.com/ferama/rospo/sshd"
	"github.com/ferama/rospo/tun"

	"github.com/spf13/cobra"
)

func init() {
	tunCmd.AddCommand(tunReverseCmd)

	tunReverseCmd.Flags().BoolP("start-sshd", "S", false, "optional start the embedded sshd")
	tunReverseCmd.Flags().StringP("sshd-authorized-keys", "K", "./authorized_keys", "ssh server authorized keys path")
	tunReverseCmd.Flags().StringP("sshd-port", "P", "2222", "the ssh server tcp port")
	tunReverseCmd.Flags().StringP("sshd-key", "I", "./server_key", "the ssh server key path")
}

var tunReverseCmd = &cobra.Command{
	Use:   "reverse [user@][server]:port",
	Short: "Creates a reverse ssh tunnel",
	Long:  `Creates a reverse ssh tunnel`,
	Example: `
  # Starts an embedded sshd and reverse proxy it to the remote server
  $ rospo tun reverse -S -r :8888 user@server:port

  # Start a reverse tunnelt from the local port 5000 to the remote 8888
  # proxing through a jump host server
  $ rospo tun reverse -l :5000 -r :8888 -j jump_host_server user@server
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startSshD, _ := cmd.Flags().GetBool("start-sshd")
		local, _ := cmd.Flags().GetString("local")
		remote, _ := cmd.Flags().GetString("remote")
		jumpHost, _ := cmd.Flags().GetString("jump-host")
		identity, _ := cmd.Flags().GetString("user-identity")
		insecure, _ := cmd.Flags().GetBool("insecure")

		config := &conf.Config{
			SshClient: &conf.SshClientConf{
				Identity:  identity,
				ServerURI: args[0],
				JumpHosts: make([]*conf.JumpHostConf, 0),
				Insecure:  insecure,
			},
			Tunnel: []*conf.TunnnelConf{
				{
					Remote:  remote,
					Local:   local,
					Forward: false,
				},
			},
		}

		if jumpHost != "" {
			config.SshClient.JumpHosts = append(config.SshClient.JumpHosts, &conf.JumpHostConf{
				URI:      jumpHost,
				Identity: identity,
			})
		}

		if startSshD {
			sshdKey, _ := cmd.Flags().GetString("sshd-key")
			sshdAuthorizedKeys, _ := cmd.Flags().GetString("sshd-authorized-keys")
			sshdPort, _ := cmd.Flags().GetString("sshd-port")

			config.SshD = &conf.SshDConf{
				Key:                sshdKey,
				AuthorizedKeysFile: sshdAuthorizedKeys,
				Port:               sshdPort,
			}
			s := sshd.NewSshServer(config.SshD)
			go s.Start()
		}

		client := sshc.NewSshConnection(config.SshClient)
		go client.Start()
		// I can easily run multiple tunnels in their respective
		// go routine here using the same client
		tun.NewTunnel(client, config.Tunnel[0]).Start()
	},
}