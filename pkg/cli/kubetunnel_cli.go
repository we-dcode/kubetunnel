package cli

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/we-dcode/kube-tunnel/pkg"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/utils/appcontext"
	"regexp"
)

func NewRootCmd() *cobra.Command {

	var kubeConfig, gcVersion, kubetunnelServerVersion, localIp, namespace, port string

	rootCommand := &cobra.Command{
		Use:   constants.KubetunnelSlug,
		Short: "Duplex interaction with K8s cluster.",
		Long:  "\"Deploy\" local service to running Kubernetes cluster and allow duplex interaction.",
		Example: fmt.Sprintf("  sudo -E %s svc --help\n", constants.KubetunnelSlug) +
			fmt.Sprintf("  sudo -E %s svc -n namespace -p '8080:80' svc_name", constants.KubetunnelSlug),

		Args: cobra.ExactArgs(1),

		// TODO: Consider change to RunE and modify all panic to return error
		Run: func(cmd *cobra.Command, args []string) {

			portForwardRegex := regexp.MustCompile(`^(?P<local>[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]):(?P<remote>[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`)

			if portForwardRegex.MatchString(port) == false {

				log.Panicf("port: '%s' is invalid. expected format example: '8080:80'", port)
			}

			matches := portForwardRegex.FindStringSubmatch(port)
			localIndex := portForwardRegex.SubexpIndex("local")
			remoteIndex := portForwardRegex.SubexpIndex("remote")

			kubeTunnel := pkg.MustNewKubeTunnel(kubeConfig, namespace)

			kubeTunnel.Run(pkg.KubeTunnelConf{
				GCVersion:         gcVersion,
				KubeTunnelVersion: kubetunnelServerVersion,
				ServiceName:       args[0],
				KubeTunnelPortMap: map[string]string{
					matches[localIndex]: matches[remoteIndex],
				},
				LocalIP: localIp,
			})

		},
	}

	// TODO: make a usage of kubeconfig explicit flag
	rootCommand.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "c", "", "absolute path to a kubectl config file")
	rootCommand.PersistentFlags().StringVar(&gcVersion, "gc-version", appcontext.Version, fmt.Sprintf("%s's Garbage Collector chart version", constants.KubetunnelSlug))
	rootCommand.PersistentFlags().StringVar(&kubetunnelServerVersion, "server-version", appcontext.Version, fmt.Sprintf("%s's Server chart version", constants.KubetunnelSlug))
	rootCommand.PersistentFlags().StringVar(&localIp, "local-ip", "127.0.0.1", "local service binding ip, usually localhost")

	rootCommand.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Specify a namespace")

	// TODO: Change port to []string and allow multi -p ...
	rootCommand.Flags().StringVarP(&port, "port", "p", "", "Specify a namespace")
	rootCommand.MarkFlagRequired("port")

	return rootCommand
}
