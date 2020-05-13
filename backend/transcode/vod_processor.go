package transcode

import (
	"os"

	"github.com/Dzionys/video-platform/backend/utils"
)

// func ProcessVodFile(source string, data vd.Vidinfo, cldata vd.Video, prdata vd.PData, conf cf.Config, clid string) {
// 	lp.WLog("Starting VOD Processor..", clid)
// 	var (
// 		err error
// 		cmd string
// 	)

// 	CONF = conf

// 	// Path to source file
// 	sfpath := CONF.SD + source

// 	// Checks if source file exists
// 	if source != "" {
// 		if _, err := os.Stat(sfpath); err == nil {
// 			lp.WLog("File found", clid)
// 		} else if os.IsNotExist(err) {
// 			lp.WLog("Error: file does not exist", clid)
// 			return
// 		} else {
// 			log.Println(err)
// 			lp.WLog("Error: file may or may not exist", clid)
// 			removeFile("videos/", source, clid)
// 			return
// 		}
// 	} else {
// 		removeFile("videos/", source, clid)
// 		return
// 	}

// 	// Full source file name
// 	fullsfname, err := filepath.EvalSymlinks(sfpath)
// 	if err != nil {
// 		log.Println(err)
// 		lp.WLog("Error: failed to get full file name", clid)
// 		removeFile("videos/", source, clid)
// 		return
// 	}

// 	// Source file name without extension
// 	sfnamewe := strings.Split(source, filepath.Ext(fullsfname))[0]

// 	// If transcoding directory does not exist creat it
// 	if _, err = os.Stat(CONF.TD); os.IsNotExist(err) {
// 		os.Mkdir(CONF.TD, 0777)
// 	}

// 	// File name after transcoding
// 	tempfile := fmt.Sprintf("%v%v.mp4", CONF.TD, sfnamewe)

// 	// f
// 	destinationfile := fmt.Sprintf("%v%v.mp4", CONF.DD, sfnamewe)

// 	// Checks if transcoded file with the same name already exists
// 	if _, err := os.Stat(tempfile); err == nil {
// 		lp.WLog(fmt.Sprintf("Error: file \"%v\" already transcoding", sfnamewe+".mp4"), clid)
// 		removeFile("videos/", source, clid)
// 		return
// 	} else if _, err := os.Stat(destinationfile); err == nil {
// 		lp.WLog(fmt.Sprintf("Error: file \"%v\" already exist in transcoded folder", sfnamewe+".mp4"), clid)
// 		removeFile("videos/", source, clid)
// 		return
// 	}

// 	lp.WLog(fmt.Sprintf("Starting to process %s", source), clid)

// 	// If data is empty get video info
// 	if data.IsEmpty() {
// 		data, err = GetVidInfo(CONF.SD, source, CONF.TempJson, CONF.DataGen, CONF.TempTxt, clid)
// 		if err != nil {
// 			log.Println(err)
// 			removeFile("videos/", source, clid)
// 			return
// 		}
// 	}

// 	// Generate thumbnails
// 	lp.WLog("Generating thumbnail", clid)
// 	wg.Add(1)
// 	err = generateThumbnail(&wg, fullsfname, sfnamewe, data)
// 	if err != nil {
// 		log.Printf("Generate thumbnail exited with error: %v", err)
// 	}
// 	wg.Wait()

// 	msg := "%v video track(s), %v audio track(s) and %v subtitle(s) found"
// 	frmt := fmt.Sprintf(msg, data.Videotracks, data.Audiotracks, data.Subtitles)
// 	lp.WLog(frmt, clid)

// 	// Generate command line
// 	var save bool
// 	var tempdfs []string
// 	if CONF.Advanced {
// 		if CONF.Presets {
// 			save = prdata.Save
// 			cmd, tempdfs, err = generatePresetCmdLine(prdata, data, sfpath, fullsfname, fmt.Sprintf("%v%v", CONF.TD, sfnamewe))
// 			tempfile = tempdfs[0]
// 			if err != nil {
// 				lp.WLog("Error: failed to generate cmd line", clid)
// 				log.Println(err)
// 				removeFile("videos/", source, clid)
// 				return
// 			}
// 		} else {
// 			save = cldata.Save
// 			cmd = generateClientCmdLine(cldata, data, sfpath, fullsfname, tempfile)
// 		}
// 	} else {
// 		cmd = generateBaseCmdLine(data, sfpath, tempfile, fullsfname)
// 	}

// 	// check if client wants to save cmd line
// 	if save {
// 		err := db.AddCmdLine(source, cmd, tempdfs)
// 		if err != nil {
// 			lp.WLog("Error: failed to insert command line in database", clid)
// 			log.Println(err)
// 			removeFile(CONF.SD, source, clid)
// 		} else {
// 			lp.WLog("Transcoding parameters saved", clid)
// 		}
// 	} else {
// 		var dfsl string
// 		for i, d := range tempdfs {
// 			if i != len(tempdfs)-1 {
// 				dfsl += d + " "
// 			} else {
// 				dfsl += d
// 			}
// 		}
// 		go StartTranscode(source, CONF, cmd, dfsl, clid)
// 	}
// }

func removeFile(path string, filename string, clid string) {
	if _, err := os.Stat(path + filename); os.Remove(path+filename) != nil && !os.IsNotExist(err) {
		utils.WLog("Error: failed removing file", clid)
	}
	//db.RemoveRowByName(filename, "Video")
	return
}
