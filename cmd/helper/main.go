

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/sjzar/chatlog/cmd/helper/extract"
	"github.com/sjzar/chatlog/internal/decrypt"
)

type KeyResult struct {
	Status string   `json:"status"`
	Keys   []string `json:"keys,omitempty"`
	Error  string   `json:"error,omitempty"`
}

func main() {
	mode := flag.String("mode", "", "Execution mode: 'extract' or 'dump'")
	key := flag.String("key", "", "Decryption key for 'dump' mode")
	src := flag.String("src", "", "Source DB path for 'dump' mode (e.g., EnMicroMsg.db)")
	dst := flag.String("dst", "decrypted.db", "Destination for decrypted DB in 'dump' mode")
	flag.Parse()

	switch *mode {
	case "extract":
		handleExtract()
	case "dump":
		handleDump(*key, *src, *dst)
	default:
		fmt.Println("Error: Invalid or missing --mode flag. Use 'extract' or 'dump'.")
		flag.Usage()
		os.Exit(1)
	}
}

func handleExtract() {
	keys, err := extract.ExtractPlatformKeys()

	output := KeyResult{}
	if err != nil {
		output.Status = "error"
		output.Error = err.Error()
	} else if len(keys) == 0 {
		output.Status = "error"
		output.Error = "no keys found"
	} else {
		output.Status = "success"
		output.Keys = keys
	}

	jsonOutput, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(jsonOutput))
}

func handleDump(key, src, dst string) {
	if key == "" || src == "" {
		fmt.Println("Error: --key and --src are required for 'dump' mode.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Attempting to decrypt '%s' with provided key...\n", src)
	err := decrypt.DecryptDB(key, src, dst)
	if err != nil {
		fmt.Printf("Error: Decryption failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Decrypted database saved to '%s'\n", dst)
}

