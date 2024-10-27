package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"path/filepath"

	"github.com/JPFrancoia/thumbnailer/thumbnailer"
	"github.com/spf13/cobra"
)

// Define a variable for the -s flag
var sFlag int

func init() {
	rootCmd.Flags().IntVar(&sFlag, "size", 300, "Size of the thumbnail")
}

var rootCmd = &cobra.Command{
	Use:   "thumbnailer [image]",
	Short: "Creates thumbnails for images",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer f.Close()

		// Get the full path
		fullPath := f.Name() // f.Name() returns the path you used in os.Open()

		// Get the file extension
		fileExt := filepath.Ext(fullPath)

		var destination string
		if fileExt != "" {
			destination = strings.Replace(fullPath, fileExt, fmt.Sprintf("_%d%s", sFlag, fileExt), 1)
		} else {
			destination = fmt.Sprintf("%s_%d", args[0], sFlag)
		}

		fmt.Println("Input file: ", fullPath)
		fmt.Println("Exporting to: ", destination)

		outFile, err := os.Create(destination)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		defer outFile.Close()

		err = thumbnailer.Thumbnail(f, outFile, sFlag)

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
