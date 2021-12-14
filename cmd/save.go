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
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/chwetion/ezsave-for-docker-image/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "easy for saving docker image",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		f := cmd.Flag("file").Value.String()
		b, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Printf("failed to read file: %s: %s", f, err)
			os.Exit(1)
		}
		cfg := model.ImportConfiguration{}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			fmt.Printf("failed to parse import cfg: %s: %s", f, err)
			os.Exit(1)
		}
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Printf("failed to create docker cli: %s", err)
		}
		ctx := context.Background()
		for _, p := range cfg.Packages {
			imageNames := make([]string, 0, 5)
			for _, content := range p.Content {
				// docker pull
				rc, err := cli.ImagePull(ctx, content.From, types.ImagePullOptions{})
				if err != nil {
					fmt.Printf("skip to pull image %s: %s", content.From, err)
					_ = rc.Close()
					continue
				}
				_, _ = io.Copy(os.Stdout, rc)
				// docker tag
				err = cli.ImageTag(ctx, content.From, content.Name)
				if err != nil {
					fmt.Printf("skip to tag image %s: %s", content.From, err)
					_ = rc.Close()
					continue
				}
				// prepare save
				_ = rc.Close()
				imageNames = append(imageNames, content.Name)
			}
			// docker save
			rc, err := cli.ImageSave(ctx, imageNames)
			if err != nil {
				fmt.Printf("save image error: %s", err)
				continue
			}
			of, err := os.Create(p.File)
			if err != nil {
				fmt.Printf("create file error: %s", err)
				_ = rc.Close()
				continue
			}
			_, err = io.Copy(of, rc)
			if err != nil {
				fmt.Printf("write file error: %s", err)
			}
			_ = rc.Close()
			_ = of.Close()
		}
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// saveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// saveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	saveCmd.Flags().StringP("file", "f", "", "Save image configuration file")
}
