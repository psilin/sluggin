package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/psilin/sluggin/scripts"
)

func main() {
	help := flag.Bool("h", false, "Prints this help message.")
	verbose := flag.Bool("v", false, "Verbose output of the script.")
	num := flag.Int("s", 50, "Number of slugs to download.")
	dsn := flag.String("dsn", "user=postgres password=postgres host=localhost dbname=postgres sslmode=disable", "DB connection sting.")
	flag.Parse()

	if *help {
		fmt.Printf("Script downloads number of slugs from %s\n", scripts.URL)
		fmt.Println("Usage: downloader [-v] [-s] [-p] [-dsn]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	start := time.Now()
	scripts.DownloadSlugs(*verbose, *num, *dsn)
	elapsed := time.Since(start)
	fmt.Printf("Slugs processing took %v\n", elapsed)
}
