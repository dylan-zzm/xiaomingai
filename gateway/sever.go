package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func StartServer(addr string) {
	http.HandleFunc("/api/extract", func(w http.ResponseWriter, r *http.Request) {
		home, err := os.UserHomeDir()
		if err != nil {
			http.Error(w, "could not get user home directory", 500)
			return
		}
		configDir := filepath.Join(home, ".xiaomingai")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			http.Error(w, "could not create config directory", 500)
			return
		}
		keyPath := filepath.Join(configDir, "key.dat")
		if err := RunHelper("extract", map[string]string{"out": keyPath}); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"keyPath": keyPath})
	})

	http.HandleFunc("/api/dump", func(w http.ResponseWriter, r *http.Request) {
		src := r.URL.Query().Get("src") // EnMicroMsg.db
		key := r.URL.Query().Get("key") // 32B hex
		if src == "" || key == "" {
			http.Error(w, "need src & key param", 400)
			return
		}
		dst := "/tmp/plain.db"
		if err := RunHelper("dump", map[string]string{
			"src": src,
			"dst": dst,
			"key": key,
		}); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"db": dst})
	})

	// --- 新增 /api/status：返回剩余积分、密钥是否存在、解密进度等 ---
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		home, _ := os.UserHomeDir()
		keyPath := filepath.Join(home, ".xiaomingai", "key.dat")
		_, keyErr := os.Stat(keyPath)
		// TODO: 积分 & 进度 后续接入数据库
		_ = json.NewEncoder(w).Encode(map[string]any{
			"keyExists": keyErr == nil,
			"points":    0, // placeholder
			"progress":  "idle",
		})
	})
	log.Printf("gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
