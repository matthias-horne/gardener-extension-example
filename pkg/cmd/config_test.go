package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompleteControllerConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name    string
		fields  string
		qps     float32
		burst   int32
		wantErr bool
	}{
		{name: "default client connection", qps: 100, burst: 130},
		{name: "custom client connection", fields: "clientConnection:\n  qps: 25\n  burst: 40\n", qps: 25, burst: 40},
		{name: "unknown configuration", fields: "unknownField: true\n", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.yaml")
			data := "apiVersion: example.extensions.config.gardener.cloud/v1alpha1\nkind: ControllerConfiguration\n" + tc.fields
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			opts := &ExampleOptions{ConfigLocation: path}
			if err := opts.Complete(); (err != nil) != tc.wantErr {
				t.Fatalf("Complete() error = %v, want error: %t", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			conn := opts.Completed().Config().ClientConnection
			if conn == nil || conn.QPS != tc.qps || conn.Burst != tc.burst {
				t.Fatalf("unexpected client connection configuration: %+v", conn)
			}
		})
	}
}
