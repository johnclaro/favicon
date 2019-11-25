/*
Copyright © 2019 John Claro <jkrclaro@gmail.com>

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
	"log"

	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "A brief description of your command",
	Long:  "TODO",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]

		src, err := imaging.Open(source)
		if err != nil {
			log.Fatalf("Failed to open image: %v", err)
		}

		src = imaging.Resize(src, 16, 16, imaging.Lanczos)

		err = imaging.Save(src, "output.png")
		if err != nil {
			log.Fatalf("Failed to save image: %v", err)
		}

	},
}

var source string
var target string

func init() {
	rootCmd.AddCommand(imageCmd)
}
