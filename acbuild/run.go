// Copyright 2015 The appc Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"strings"

	"github.com/appc/acbuild/engine"
	"github.com/appc/acbuild/engine/chroot"
	"github.com/appc/acbuild/engine/systemdnspawn"

	"github.com/spf13/cobra"
)

var (
	insecure      = false
	workingdir    = ""
	engineName    = ""
	authConfigDir = ""

	cmdRun = &cobra.Command{
		Use:     "run -- CMD [ARGS]",
		Short:   "Run a command in an ACI",
		Long:    "Run a given command in an ACI, and save the resulting container as a new ACI",
		Example: "acbuild run -- yum install nginx",
		Run:     runWrapper(runRun),
	}

	engines = map[string]engine.Engine{
		"systemd-nspawn": systemdnspawn.Engine{},
		"chroot":         chroot.Engine{},
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdRun)

	var engineNames []string
	for engine, _ := range engines {
		engineNames = append(engineNames, engine)
	}
	engineList := fmt.Sprintf("[%s]", strings.Join(engineNames, ","))

	cmdRun.Flags().BoolVar(&insecure, "insecure", false, "Allows fetching dependencies over http")
	cmdRun.Flags().StringVar(&workingdir, "working-dir", "", "The working directory inside the container for this command")
	cmdRun.Flags().StringVar(&engineName, "engine", "systemd-nspawn", "The engine used to run the command. Supported engines: "+engineList)
	cmdRun.Flags().StringVar(&authConfigDir, "auth-config-dir", "auth.d", "Directory with authentication config file(s)")
}

func runRun(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}

	if debug {
		stderr("Running: %v", args)
	}

	engine, ok := engines[engineName]
	if !ok {
		stderr("run: no such engine %q", engineName)
		return 1
	}

	err := newACBuild().Run(args, workingdir, insecure, engine, authConfigDir)

	if err != nil {
		stderr("run: %v", err)
		return getErrorCode(err)
	}

	return 0
}
