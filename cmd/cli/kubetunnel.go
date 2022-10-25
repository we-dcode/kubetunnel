package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/we-dcode/kube-tunnel/cmd/cli/cmds"
	"github.com/we-dcode/kube-tunnel/pkg"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/utils/logutil"
	"io/ioutil"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

var globalUsage = ``
var Version = "0.0.0"

func init() {
	// quiet version
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "version" && args[1] == "quiet" {
		fmt.Println(Version)
		os.Exit(0)
	}

	log.SetOutput(&logutil.LogOutputSplitter{})
	if len(args) > 0 && args[0] == "completion" {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05",
	})

	//log.Print("")
	//log.Print("https://github.com/we-dcode/kube-tunnel")
	//log.Print("https://dcode.tech")
	//log.Print("")

	cmd := NewRootCmd()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   constants.KubetunnelSlug,
		Short: "Duplex interaction with K8s cluster.",
		Long:  "\"Deploy\" local service to running Kubernetes cluster and allow duplex interaction.",
	}

	rootCmd.AddCommand(NewSvcCmd())
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(cmds.NewCmdCompletion(os.Stdout, ""))
	//rootCmd.SetHelpCommand()

	return rootCmd
}

func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the version of %s", constants.KubetunnelSlug),
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version: %s\nhttps://github.com/we-dcode/kube-tunnel\n", constants.KubetunnelSlug, Version)
		},
	}

	return versionCmd
}

func NewSvcCmd() *cobra.Command {

	var kubeConfig, gcVersion, kubetunnelServerVersion, localIp, namespace, port string

	svcCmd := &cobra.Command{
		Use:   "svc",
		Short: "Duplex interaction with K8s cluster.",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf("  sudo -E %s svc --help\n", constants.KubetunnelSlug) +
			fmt.Sprintf("  sudo -E %s svc -p '8080:80' svc_name\n", constants.KubetunnelSlug) +
			fmt.Sprintf("  sudo -E %s svc -c kubeconfig/path -p '8080:80' svc_name\n", constants.KubetunnelSlug) +
			fmt.Sprintf("  sudo -E %s svc -c kubeconfig/path -n namespace -p '8080:80' svc_name\n", constants.KubetunnelSlug),
		//ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		//
		//	log.SetLevel(log.PanicLevel)
		//
		//	k := kube.MustNew(kubeConfig, namespace)
		//
		//	serviceNames, _ := k.ListServiceNamesWithoutKubeTunnel()
		//
		//	return serviceNames, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		//},
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

			kubeTunnel.CreateTunnel(pkg.KubeTunnelConf{
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

	svcCmd.Flags().StringVarP(&kubeConfig, "kubeconfig", "c", "", "absolute path to a kubectl config file.")
	svcCmd.Flags().StringVar(&gcVersion, "gc-version", Version, fmt.Sprintf("%s's Garbage Collector chart version.", constants.KubetunnelSlug))
	svcCmd.Flags().StringVar(&kubetunnelServerVersion, "server-version", Version, fmt.Sprintf("%s's Server chart version.", constants.KubetunnelSlug))
	svcCmd.Flags().StringVar(&localIp, "local-ip", "127.0.0.1", "local service binding ip, usually localhost.")

	svcCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "specify namespace, default: taken from kubeconfig's context.")

	// TODO: Change port to []string and allow multi -p ...
	svcCmd.Flags().StringVarP(&port, "port", "p", "", "localPort:remotePort (example: 8080:80).")
	svcCmd.MarkFlagRequired("port")

	return svcCmd
}
