/*
Copyright © 2021 author chwetion chwetion@foxmail.com

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
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/chwetion/ezsave-for-docker-image/pkg/docker"
	"github.com/chwetion/ezsave-for-docker-image/pkg/util"
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "easy for loading image",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		filePath := cmd.Flag("file").Value.String()
		cfg, err := util.LoadConfiguration(filePath)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		client, err := docker.NewClient(cfg, path.Dir(filePath))
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		client.Load()
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loadCmd.Flags().StringP("file", "f", "", "Load image configuration file")
	// saveCmd.Flags().BoolP("local", "", false, "load compressed package to local image not push")
}
