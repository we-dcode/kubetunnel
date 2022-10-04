/*
Copyright 2018 Craig Johnston <cjimti@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"github.com/we-dcode/kube-tunnel/pkg"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/utils/logutil"
	"io/ioutil"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

func NewRootCmd() *cobra.Command {

	var kubeConfig, gcVersion, kubetunnelServerVersion, localIp, namespace, port string

	rootCommand := &cobra.Command{
		Use:   constants.KubetunnelSlug,
		Short: "Duplex interaction with K8s cluster.",
		Long:  "\"Deploy\" local service to running Kubernetes cluster and allow duplex interaction.",
		Example: fmt.Sprintf("  sudo -E %s svc --help\n", constants.KubetunnelSlug) +
			fmt.Sprintf("  sudo -E %s svc -n namespace -p 8080:80 svc_name", constants.KubetunnelSlug),

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
	rootCommand.PersistentFlags().StringVar(&gcVersion, "gc-version", Version, fmt.Sprintf("%s's Garbage Collector chart version", constants.KubetunnelSlug))
	rootCommand.PersistentFlags().StringVar(&kubetunnelServerVersion, "server-version", Version, fmt.Sprintf("%s's Server chart version", constants.KubetunnelSlug))
	rootCommand.PersistentFlags().StringVar(&localIp, "local-ip", "127.0.0.1", "local service binding ip, usually localhost")

	rootCommand.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Specify a namespace")

	// TODO: Change port to []string and allow multi -p ...
	rootCommand.Flags().StringVarP(&port, "port", "p", "", "Specify a namespace")
	rootCommand.MarkFlagRequired("port")

	return rootCommand
}

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05",
	})

	log.Print("")
	log.Print("https://github.com/we-dcode/kube-tunnel")
	log.Print("https://we.dcode.tech")
	log.Print("")

	cmd := NewRootCmd()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
