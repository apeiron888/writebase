package config

import (
	"os"
	"testing"
)

func withEnv(kv map[string]string, fn func()) {
	// save old
	old := map[string]*string{}
	for k := range kv {
		if v, ok := os.LookupEnv(k); ok {
			vv := v
			old[k] = &vv
		} else {
			old[k] = nil
		}
		os.Setenv(k, kv[k])
	}
	defer func() {
		for k, v := range old {
			if v == nil {
				_ = os.Unsetenv(k)
			} else {
				os.Setenv(k, *v)
			}
		}
	}()
	fn()
}

func TestLoadEnv_Success(t *testing.T) {
	withEnv(map[string]string{
		"MONGODB_URI":    "mongodb://localhost:27017",
		"MONGODB_NAME":   "write_base",
		"JWT_SECRET":     "secret",
		"SERVER_PORT":    "8080",
		"GEMINI_API_KEY": "key",
	}, func() {
		cfg, err := LoadEnv()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.MongodbURI == "" || cfg.JwtSecret == "" {
			t.Fatalf("missing fields in cfg")
		}
	})
}

func TestLoadEnv_MissingVariables(t *testing.T) {
	withEnv(map[string]string{
		"MONGODB_URI":    "",
		"MONGODB_NAME":   "",
		"JWT_SECRET":     "",
		"SERVER_PORT":    "",
		"GEMINI_API_KEY": "",
	}, func() {
		_, err := LoadEnv()
		if err == nil {
			t.Fatalf("expected error for missing env, got nil")
		}
	})
}
