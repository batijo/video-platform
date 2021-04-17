package transcode

import (
	"sync"
	"testing"
)

type Info struct {
	source string
	expout string
	experr string
}

var infos = []Info{
	{"", `{}`, "no video file provided"},
	{"sdd", `{}`, "json data is empty"},
}

func TestGetMediaInfoJSON(t *testing.T) {
	var wg sync.WaitGroup
	for _, inf := range infos {
		wg.Add(1)
		out, err := getMediaInfoJSON(inf.source, &wg)
		wg.Wait()
		if err != nil {
			if err.Error() != inf.experr {
				t.Errorf("DATA: %v EXPECTED: %v, GOT: %v", inf.source, inf.experr, err.Error())
			}
		} else if string(out) != inf.expout {
			t.Errorf("DATA: %v EXPECTED: %v, GOT: %v", inf.source, inf.expout, string(out))
		}
	}

}

type VidInfo struct {
	path     string
	filename string
	clientID string
	vidID    int
	experr   string
}

func TestGetVidInfo(t *testing.T) {

}
