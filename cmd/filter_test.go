package cmd

import (
	"os"
	"testing"
)

func TestFilterCheck(t *testing.T) {
	if os.Getenv("ADCTL_HOST") == "" || os.Getenv("ADCTL_USERNAME") == "" || os.Getenv("ADCTL_PASSWORD") == "" {
		t.Skip("integration test requires ADCTL_HOST, ADCTL_USERNAME, and ADCTL_PASSWORD")
	}

	cfa := CheckFilterArgs{name: "www.doubleclick.net"}
	_, err := GetFilter(nil, cfa)
	if err != nil {
		t.Errorf("error in GetFilter: %v", err)
	}
}
