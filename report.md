# DATA COMPRESSION UTILIZING L-Z-W ALGORITHM (TEXT FILE COMPRESSION USING GOLANG)

## THEORY

Data compression is the process of reducing the size of digital data while preserving its essential information. It plays a crucial role in file storage, communication, and multimedia applications.

### Types of Data Compression

#### Lossless Compression
- No data is lost during the process.
- Achieves compression by identifying and eliminating redundancy within the data.
- Common algorithms: Run Length Encoding (RLE), Huffman coding, Lempel–Ziv–Welch (LZW).
- Ideal for text files, databases, and program executables.

#### Lossy Compression
- Sacrifices some data quality to achieve higher compression ratios.
- Commonly used for multimedia files (images, audio, video).
- Examples: JPEG (images), MP3 (audio).
- Trade-off: some details are permanently lost.

### The Lempel–Ziv–Welch (LZW) Algorithm
LZW is a **lossless dictionary-based compression algorithm**. It dynamically builds a dictionary of patterns found in the input and replaces repeated patterns with codes.

#### Initialization
- Start with a dictionary containing all single-character symbols (typically 256 ASCII characters).

#### Encoding Process
1. Read the input stream symbol by symbol.
2. Maintain a current string `w`.
3. For each next character `c`, form `w+c`.
4. If `w+c` is in the dictionary, set `w = w+c`.
5. Otherwise:
   - Output the code for `w`.
   - Add `w+c` to the dictionary.
   - Set `w = c`.
6. At the end, output the code for the remaining `w`.

#### Decoding Process
- Initialize the same starting dictionary.
- Read codes, map them back to strings, and rebuild the dictionary in the same order as encoding.
- This reconstructs the original text exactly.

---

## CODE AND OUTPUT

### go.mod
```text
module FilecompressionGolang

go 1.22.0
```

### main.go (Golang Program)
```go
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
```

---

### Input File (input.txt)
```text
NOTE: This is a test file for LZW compression. This is a test file for LZW compression.
NOTE: This is a test file for LZW compression. This is a test file for LZW compression.

The quick brown fox jumps over the lazy dog.
The quick brown fox jumps over the lazy dog.
The quick brown fox jumps over the lazy dog.
The quick brown fox jumps over the lazy dog.

Kathmandu Kathmandu Kathmandu Kathmandu Kathmandu Kathmandu Kathmandu Kathmandu Kathmandu Kathmandu
compression compression compression compression compression compression compression compression compression

PATTERN-ABC-123 PATTERN-ABC-123 PATTERN-ABC-123 PATTERN-ABC-123 PATTERN-ABC-123
PATTERN-ABC-123 PATTERN-ABC-123 PATTERN-ABC-123 PATTERN-ABC-123 PATTERN-ABC-123

This is a test file for LZW compression. This is a test file for LZW compression.
This is a test file for LZW compression. This is a test file for LZW compression.
```

---

### Commands Used

#### Compression
```powershell
go run . -mode=c -in ".\input.txt" -out ".\input.txt.lzw"
```

#### Decompression
```powershell
go run . -mode=d -in ".\input.txt.lzw" -out ".\restored.txt"
```

---

## RESULT AND DISCUSSION

- The program successfully compressed the text file `input.txt` into a compressed file `input.txt.lzw` using LZW in Golang.  
- After running decompression, the output file `restored.txt` was generated and the content matched the original `input.txt`, proving that the method is **lossless**.  
- Since the input file contains repeated substrings (e.g., repeated sentences and repeated patterns like `PATTERN-ABC-123`), the LZW dictionary grows with these patterns and later represents them using shorter codes, which improves compression efficiency.

---

## CONCLUSION

We performed **lossless data compression** on a text file utilizing the **L-Z-W algorithm** by running a Golang program.  
The file was successfully compressed and decompressed, and the decompressed file `restored.txt` matched the original `input.txt`.  
This confirms that the implemented LZW compression and decompression preserve the original data correctly.

---
