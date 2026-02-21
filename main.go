package main

import (
	"compress/lzw"
	"flag"
	"fmt"
	"io"
	"os"
)

func compressFile(inPath, outPath string) error {
	in, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer in.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	// LSB + litWidth=8 is the common setup for byte-oriented data.
	w := lzw.NewWriter(out, lzw.LSB, 8)
	defer w.Close() // must close to flush final codes

	if _, err := io.Copy(w, in); err != nil {
		return fmt.Errorf("compress copy: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close compressor: %w", err)
	}
	return nil
}

func decompressFile(inPath, outPath string) error {
	in, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer in.Close()

	r := lzw.NewReader(in, lzw.LSB, 8)
	defer r.Close()

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return fmt.Errorf("decompress copy: %w", err)
	}
	return nil
}

func main() {
	mode := flag.String("mode", "c", "c=compress, d=decompress")
	in := flag.String("in", "", "input file path")
	out := flag.String("out", "", "output file path")
	flag.Parse()

	if *in == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "Usage:\n  go run . -mode=c -in=input.txt -out=output.lzw\n  go run . -mode=d -in=output.lzw -out=restored.txt")
		os.Exit(2)
	}

	var err error
	switch *mode {
	case "c":
		err = compressFile(*in, *out)
	case "d":
		err = decompressFile(*in, *out)
	default:
		err = fmt.Errorf("unknown mode: %q (use c or d)", *mode)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
