package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"runtime/pprof"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: findimg [options] <image> <subimage>\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	output      = flag.String("o", "text", "result output format (json, text)")
	random      = flag.Bool("random", false, "randomly pick subimage as test")
	verbose     = flag.Bool("v", false, "verbose output")
	cpuProfile  = flag.String("cpu-profile", "", "write cpu profile to file")
	imgMinWidth = flag.Int("img-min-width", 0, "minimum image width")
	imgMaxWidth = flag.Int("img-max-width", 0, "maximum image width")
	subMinArea  = flag.Int("sub-min-area", 0, "minimum subimage area")
	subMaxDiv   = flag.Int("sub-max-div", 0, "maximum subimage division")
	k           = flag.Int("k", 1, "number of top matches to keep")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("findimg: ")

	flag.Usage = usage
	flag.Parse()

	imgPath := flag.Arg(0)
	subimgPath := flag.Arg(1)

	if imgPath == "" || (subimgPath == "" && !*random) {
		usage()
	}

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	// Open the input images
	imgsrc, err := openImage(imgPath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	var subsrc image.Image
	if *random {
		subsrc = randomSubimage(imgsrc)
	} else {
		subsrc, err = openImage(subimgPath)
		if err != nil {
			log.Fatalf("failed to open image: %v", err)
		}
	}

	opts := Opts{}
	opts.html = *output == "html"
	opts.verbose = *verbose
	opts.imgMinWidth = *imgMinWidth
	opts.imgMaxWidth = *imgMaxWidth
	opts.subMinArea = *subMinArea
	opts.subMaxDiv = *subMaxDiv
	opts.k = *k

	if opts.html {
		opts.convolution = true
		opts.visualize = true
	}

	matches := findImage(imgsrc, subsrc, opts)

	switch *output {
	case "json":
		json.NewEncoder(os.Stdout).Encode(struct {
			Matches Matches `json:"matches"`
		}{
			Matches: matches,
		})
	case "html":
	default:
		for _, match := range matches {
			fmt.Printf(
				"%6f %4d %4d %4d %4d\n",
				match.Match,
				match.Bounds.Min.X,
				match.Bounds.Min.Y,
				match.Bounds.Dx(),
				match.Bounds.Dy(),
			)
		}
	}
}
