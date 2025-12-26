package main

import (
	"fmt"
	"os"
	"path/filepath"

	"vb6enc/internal/converter"
	"vb6enc/internal/detector"
	"vb6enc/internal/walker"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "scan":
		runScan(args)
	case "to-utf8":
		runConvert(args, true)
	case "to-gb":
		runConvert(args, false)
	case "verify":
		runVerify(args)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: vb6enc <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  scan [path]      Scan directory and report encoding of files")
	fmt.Println("  to-utf8 [path]   Convert GBK files to UTF-8")
	fmt.Println("  to-gb [path]     Convert UTF-8 files to GBK")
	fmt.Println("  verify [path]    Report files with unknown encoding")
}

func parsePath(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	cwd, _ := os.Getwd()
	return cwd
}

func runScan(args []string) {
	root := parsePath(args)
	fmt.Printf("Scanning directory: %s\n", root)

	w := walker.New(walker.DefaultConfig())
	files, err := w.Walk(root)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		enc, err := detector.DetectFile(file)
		if err != nil {
			fmt.Printf("[Error] %s: %v\n", file, err)
			continue
		}
		fmt.Printf("[%s] %s\n", enc, file)
	}
}

func runConvert(args []string, toUTF8 bool) {
	root := parsePath(args)
	targetEnc := "UTF-8"
	if !toUTF8 {
		targetEnc = "GBK"
	}
	fmt.Printf("Converting files to %s in: %s\n", targetEnc, root)

	w := walker.New(walker.DefaultConfig())
	files, err := w.Walk(root)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	// var wg sync.WaitGroup
	// Limit concurrency? Maybe, but disk I/O is the bottleneck.
	// Let's simple semaphore it or just sequential for safety and log clarity.
	// Sequential is safer for "atomic" feeling and logging.

	successCount := 0
	skipCount := 0
	failCount := 0

	for _, file := range files {
		enc, err := detector.DetectFile(file)
		if err != nil {
			fmt.Printf("[Error] Detect %s: %v\n", file, err)
			failCount++
			continue
		}

		shouldConvert := false
		if toUTF8 {
			if enc == detector.GBK {
				shouldConvert = true
			}
		} else {
			if enc == detector.UTF8 {
				shouldConvert = true
			}
		}

		if !shouldConvert {
			skipCount++
			continue
		}

		fmt.Printf("Converting %s (%s -> %s)... ", filepath.Base(file), enc, targetEnc)
		err = converter.ConvertFile(file, toUTF8)
		if err != nil {
			fmt.Printf("Failed: %v\n", err)
			failCount++
		} else {
			fmt.Printf("Done\n")
			successCount++
		}
	}

	fmt.Printf("\nSummary: %d converted, %d skipped, %d failed\n", successCount, skipCount, failCount)
}

func runVerify(args []string) {
	root := parsePath(args)
	fmt.Printf("Verifying files in: %s\n", root)

	w := walker.New(walker.DefaultConfig())
	files, err := w.Walk(root)
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	issues := 0
	for _, file := range files {
		enc, err := detector.DetectFile(file)
		if err != nil {
			fmt.Printf("[Error] %s: %v\n", file, err)
			issues++
			continue
		}
		if enc == detector.Unknown {
			fmt.Printf("[UNKNOWN] %s\n", file)
			issues++
		}
	}

	if issues == 0 {
		fmt.Println("All files have valid encodings (UTF-8 or GBK).")
	} else {
		fmt.Printf("Found %d files with issues.\n", issues)
	}
}
