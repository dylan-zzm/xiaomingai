package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/dylan-zzm/xiaomingai/cmd/helper/extract"
	"github.com/dylan-zzm/xiaomingai/internal/decrypt"
)

// KeyResult 用于 extract 模式下的 JSON 输出结构
type KeyResult struct {
	Status string   `json:"status"`
	Keys   []string `json:"keys,omitempty"`
	Error  string   `json:"error,omitempty"`
}

func main() {
	// 自定义 Usage 方便 --help 查看
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	mode := flag.String("mode", "", "Execution mode: 'extract' (get key) or 'dump' (decrypt db)")
	key := flag.String("key", "", "Decryption key (for 'dump' mode)")
	src := flag.String("src", "", "Source DB path (for 'dump' mode, e.g., EnMicroMsg.db)")
	dst := flag.String("dst", "decrypted.db", "Destination for decrypted DB (for 'dump' mode)")
	out := flag.String("out", "", "Alias of --dst / optional output path for 'extract' mode")
	flag.Parse()

	if *mode == "" {
		fmt.Fprintln(os.Stderr, "Error: --mode is required.")
		flag.Usage()
		os.Exit(1)
	}

	switch *mode {
	case "extract":
		handleExtract(*out)
	case "dump":
		handleDump(*key, *src, *dst)
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid mode '%s'. Use 'extract' or 'dump'.\n", *mode)
		flag.Usage()
		os.Exit(1)
	}
}

// handleExtract 调用密钥提取逻辑并以 JSON 形式输出结果
func handleExtract(outPath string) {
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

	if outPath != "" && output.Status == "success" {
		_ = os.WriteFile(outPath, []byte(keys[0]), 0600)
	}

	if output.Status != "success" {
		os.Exit(1)
	}
	os.Exit(0)
}

// handleDump 解密数据库
func handleDump(key, src, dst string) {
	if key == "" || src == "" {
		fmt.Fprintln(os.Stderr, "Error: --key and --src are required for 'dump' mode.")
		flag.Usage()
		os.Exit(1)
	}

	if err := decrypt.DumpDB(src, dst, []byte(key)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Decryption failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Decrypted database saved to '%s'\n", dst)
	os.Exit(0)
}
