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
	"github.com/we-dcode/kube-tunnel/pkg/cli"
	"github.com/we-dcode/kube-tunnel/pkg/utils/appcontext"
	"github.com/we-dcode/kube-tunnel/pkg/utils/logutil"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

var globalUsage = ``

func init() {
	// quiet version
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "version" && args[1] == "quiet" {
		fmt.Println(appcontext.Version)
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

	log.Print("")
	log.Print("https://github.com/we-dcode/kube-tunnel")
	log.Print("https://we.dcode.tech")
	log.Print("")

	cmd := cli.NewRootCmd()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
