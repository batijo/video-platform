package transcode

import (
	"os"

	"github.com/Dzionys/video-platform/backend/utils"
)

func removeFile(path string, filename string, clid string) {
	if _, err := os.Stat(path + filename); os.Remove(path+filename) != nil && !os.IsNotExist(err) {
		utils.WLog("Error: failed removing file", clid)
	}
	//db.RemoveRowByName(filename, "Video")
	return
}
