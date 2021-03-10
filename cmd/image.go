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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
)

type favicon struct {
	height, width int
	name, tag string
}

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Convert an image to its favicon formats.",
	Long: `Convert an image that you want to use as your favicon then use this 
	tool to convert an image to its favicon formats.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Requires `source` and `target`")
		}
		if len(args) >= 3 {
			return errors.New("Too many arguments provided")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		target := args[1]

		files := map[string]favicon{
			"apple-touch-icon-57x57.png":   favicon{height: 50, width: 50, name: "apple-touch-icon-precomposed", tag: "link"},
			"apple-touch-icon-60x60.png":   favicon{height: 60, width: 60, name: "apple-touch-icon-precomposed", tag: "link"},
			"apple-touch-icon-72x72.png":   favicon{height: 72, width: 72, name: "apple-touch-icon-precomposed", tag: "link"},
			"apple-touch-icon-76x76.png":   favicon{height: 76, width: 76, name: "apple-touch-icon-precomposed", tag: "link"},
			"apple-touch-icon-114x114.png": favicon{height: 114, width: 114, name: "apple-touch-icon-precomposed", tag: "link"},
			"apple-touch-icon-120x120.png": favicon{height: 120, width: 120, name: "apple-touch-icon-precomposed", tag: "link"},
			"apple-touch-icon-144x144.png": favicon{height: 144, width: 144, name: "apple-touch-icon-precomposed", tag: "link"},
			"apple-touch-icon-152x152.png": favicon{height: 152, width: 152, name: "apple-touch-icon-precomposed", tag: "link"},
			"favicon-16x16.png":            favicon{height: 16, width: 16, name: "icon", tag: "link"},
			"favicon-32x32.png":            favicon{height: 32, width: 32, name: "icon", tag: "link"},
			"favicon-96x96.png":            favicon{height: 96, width: 96, name: "icon", tag: "link"},
			"favicon-128.png":              favicon{height: 128, width: 128, name: "icon", tag: "link"},
			"favicon-196x196.png":          favicon{height: 196, width: 196, name: "icon", tag: "link"},
			"mstile-70x70.png":             favicon{height: 70, width: 70, name: "msapplication-square70x70logo", tag: "meta"},
			"ms-title-144x144.png":         favicon{height: 144, width: 144, name: "msapplication-TileImage", tag: "meta"},
			"mstile-150x150.png":           favicon{height: 150, width: 150, name: "msapplication-wide150x150logo", tag: "meta"},
			"mstile-310x310.png":           favicon{height: 310, width: 310, name: "msapplication-wide310x310logo", tag: "meta"},
			"favicon.ico":                  favicon{height: 64, width: 64, name: "shortcut icon", tag: "meta"},
			"mstile-310x150.png":           favicon{height: 310, width: 150, name: "msapplication-wide310x150logo", tag: "meta"},
		}

		for filename, favicon := range files {
			filepath := filepath.Join(target, filename)

			if filename == "favicon.ico" {
				file, err := os.Open(source)
				if err != nil {
					fmt.Println(err)
				}
				defer file.Close()

				pngFile, err := png.Decode(file)
				if err != nil {
					fmt.Println(err)
				}

				writerICO, _ := os.Create(filepath)
				defer writerICO.Close()
				bounds := pngFile.Bounds()
				rgba := image.NewRGBA(bounds)
				draw.Draw(rgba, bounds, pngFile, bounds.Min, draw.Src)
				buffer := new(bytes.Buffer)
				pngWriter := bufio.NewWriter(buffer)
				png.Encode(pngWriter, rgba)
				pngWriter.Flush()
				writerICO.Write(buffer.Bytes())
			} else {
				imagefile, err := imaging.Open(source)
				if err != nil {
					log.Fatalf("Failed to open image: %v", err)
				}

				imagefile = imaging.Resize(
					imagefile,
					favicon.height,
					favicon.width,
					imaging.Lanczos,
				)
				err = imaging.Save(imagefile, filepath)
				fmt.Println(favicon.name, favicon.tag)
				if err != nil {
					log.Fatalf("Failed to save image: %v", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
}
