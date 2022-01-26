package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/psilin/sluggin/scripts"
)

func main() {
	help := flag.Bool("h", false, "Prints this help message.")
	verbose := flag.Bool("v", false, "Verbose output of the script.")
	num := flag.Int("s", 50, "Number of slugs to download.")
	path := flag.String("p", "", "Path where HTML files should be written (no path - no files).")
	flag.Parse()

	fmt.Printf("help: %v\n", *help)

	if *help {
		fmt.Printf("Script downloads number of slugs from %s\n", scripts.URL)
		fmt.Println("Usage: downloader [-v] [-s] [-p]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	scripts.DownloadSlugs(*verbose, *num, *path)
}
