package task

import (
	"testing"
)

func TestXuiScan(t *testing.T) {
	links := []string{"http://178.253.52.188:9090", "http://47.238.103.25", "http://117.18.124.83:2333"}
	res := ScanLinks(links)
	t.Logf("Total: %d", len(res))
	for _, r := range res {
		t.Log(r)
	}
}
