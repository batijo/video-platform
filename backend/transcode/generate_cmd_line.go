package transcode

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Dzionys/video-platform/backend/models"
	"github.com/Dzionys/video-platform/backend/utils"
)

var prRes = map[string]string{
	"240p":  "352x240",
	"576p":  "720x576",
	"720p":  "1280x720",
	"360p":  "480x360",
	"1080p": "1920x1080",
}

func generatePresetCmdLine(prdata models.Pdata, vdata models.Vidinfo, sf string, sfname string, dfwe string) (string, []string, error) {
	var (
		cmd       = ""
		mapping   []string
		fcmaps    []string
		vcode     []string
		acode     []string
		scode     []string
		fc        = ""
		debugIntr = ""
		dfs       []string
		tempvc    = ""
		tempac    = ""
		tempsc    = ""
		tempmp    = ""
	)

	// Checks if debuging is set to true
	if utils.Conf.Debug {
		debugIntr += " -ss " + utils.Conf.DebugStart + " -t " + utils.Conf.DebugEnd
	}

	// Video part ---------------------------------------------

	if vdata.Videotrack[0].FrameRate < 25 {
		var (
			vo string
			ao string
		)
		fc = "-r 25 -filter_complex [0:v]setpts=%[1]v/25*PTS,split=%[2]v%[3]v;[0:a]atempo=25/%[1]v,asplit=%[2]v%[4]v"

		for i := range prdata.Streams {
			fcmaps = append(fcmaps, fmt.Sprintf(" -map [vo%[1]v] -map [ao%[1]v]", i))
			vo += fmt.Sprintf("[vo%v]", i)
			ao += fmt.Sprintf("[ao%v]", i)
		}

		fc = fmt.Sprintf(fc, vdata.Videotrack[0].FrameRate, len(prdata.Streams), vo, ao)
	}

	cmd = fmt.Sprintf("ffmpeg -i %v %v", sf, fc)

	for i, s := range prdata.Streams {

		vidpr, err := utils.GetPreset(s.VidPreset)
		if err != nil {
			return cmd, dfs, err
		}
		audpr, err := utils.GetPreset(s.AudPreset)
		if err != nil {
			return cmd, dfs, err
		}

		svtres := strconv.Itoa(vdata.Videotrack[0].Width) + "x" + strconv.Itoa(vdata.Videotrack[0].Height)
		if prRes[vidpr.Resolution] != svtres {
			tempvc += fmt.Sprintf(" -s %v", prRes[vidpr.Resolution])
		}

		if vdata.Videotrack[0].FrameRate < 25 && i != 0 {
			tempvc += " -r 25"
		}

		if vidpr.Codec != vdata.Videotrack[0].CodecName {
			switch vidpr.Codec {

			case "h264":
				tempvc += fmt.Sprintf(" -c:v:%[1]v libx264 -profile:v:%[1]v main -b:v:%[1]v %[2]vk -metadata:s:v:%[1]v name=\"%[3]v\"", s.VtId, vidpr.Bitrate, sfname)
				break

			case "hevc":
				templ := fmt.Sprintf(" -c:v:%v libx265 -x265-params \"preset=slower:me=hex:no-rect=1:no-amp=1:rd=4:aq-mode=2:", s.VtId)
				templ += "aq-strength=0.5:psy-rd=1.0:psy-rdoq=0.2:bframes=3:min-keyint=1\" "
				templ += fmt.Sprintf("-b:v:0 %vk -metadata:s:v:0 name=\"%v\"", vidpr.Bitrate, sfname)
				tempvc += templ

			case "default":
				if !(vdata.Videotrack[0].FrameRate < 25) {
					tempmp += fmt.Sprintf(" -map 0:%v", s.VtId)
				}
			}
		} else if !(vdata.Videotrack[0].FrameRate < 25) {
			tempmp += fmt.Sprintf(" -map 0:%v", s.VtId)
		}

		// Audio part ---------------------------------------------

		if s.AudioT[0].StreamID != -1 {
			for i, at := range s.AudioT {
				if !(vdata.Videotrack[0].FrameRate < 25) {
					tempmp += fmt.Sprintf(" -map 0:%v", at.StreamID)
				}
				tempac += fmt.Sprintf(" -c:a:%v libfdk_aac -ac 2 -b:a:%v %vk -metadata language=%v", i, i, audpr.Bitrate, at.Language)
			}
		} else {
			for i, at := range vdata.Audiotrack {
				if !(vdata.Videotrack[0].FrameRate < 25) {
					tempmp += fmt.Sprintf(" -map 0:%v", at.Index)
				}
				tempac += fmt.Sprintf(" -c:a:%v libfdk_aac -ac 2 -b:a:%v %vk -metadata language=%v", i, i, audpr.Bitrate, at.Language)
			}
		}

		// Subtitle part ------------------------------------------

		for i := range s.SubtitleT {
			if i > 0 {
				break
			}
			if s.SubtitleT[0].StreamID != -1 {
				for _, st := range s.SubtitleT {
					tempsc += fmt.Sprintf(" -c:s:%[1]v copy -metadata:s:s:%[1]v language=%[2]v", st.StreamID, st.Language)
				}
			} else {
				for _, st := range vdata.Subtitle {
					tempsc += fmt.Sprintf(" -c:s:%[1]v copy -metadata:s:s:%[1]v language=%[2]v", st.Index, st.Language)
				}
			}
		}

		// Creates output file names
		dfpat := "%v-%v-%v-%v%v"
		dfs = append(dfs, fmt.Sprintf(dfpat, dfwe, vidpr.Resolution, vidpr.Codec, audpr.Codec, ".mp4"))

		vcode = append(vcode, tempvc)
		acode = append(acode, tempac)
		scode = append(scode, tempsc)
		mapping = append(mapping, tempmp)

		if len(fcmaps) < i+1 {
			fcmaps = append(fcmaps, "")
		}

		// Create cmd line
		tcmd := "%v %v %v %v %v %v -async 1 -vsync 1 %v"
		cmd += fmt.Sprintf(tcmd, debugIntr, fcmaps[i], vcode[i], acode[i], scode[i], mapping[i], dfs[i])

		tempvc = ""
		tempac = ""
		tempsc = ""
		tempmp = ""

	}

	return cmd, dfs, nil
}

func generateClientCmdLine(crdata models.Video, vdata models.Vidinfo, sf string, sfname string, df string) string {
	var (
		cmd       = ""
		mapping   = ""
		vcode     = ""
		acode     = ""
		scode     = ""
		debugIntr = ""
	)

	// Checks if debuging is set to true
	if utils.Conf.Debug {
		debugIntr = "-ss " + utils.Conf.DebugStart + " -t " + utils.Conf.DebugEnd
	}

	// Video part ---------------------------------------------

	// Change frames per second and/or resolution
	svtres := strconv.Itoa(vdata.Videotrack[0].Width) + ":" + strconv.Itoa(vdata.Videotrack[0].Height)
	crres := strconv.Itoa(crdata.Width) + ":" + strconv.Itoa(crdata.Height)
	if crdata.FrameRate != vdata.Videotrack[0].FrameRate || crres != svtres {

		res := strings.Split(crres, ":")
		var (
			filterComplex = ""
			vpipe         = "0:v"
			fps           = ""
			maps          = ""
		)

		// Change resolution
		if crres != svtres {
			vpipe = "scaled"
			maps += "-map [scaled]"
			filterComplex += fmt.Sprintf("[0:v]scale=%v:%v[%v]", res[0], res[1], vpipe)
		}

		// Change frame rate
		if crdata.FrameRate != vdata.Videotrack[0].FrameRate {
			maps = "-map [v] -map [a]"
			fps = fmt.Sprintf(" -r %v", crdata.FrameRate)
			var bline string
			if vpipe == "scaled" {
				bline = ";[%[3]v]setpts=%[2]v/%[1]v*PTS[v];[0:a]atempo=%[1]v/%[2]v[a]"
			} else {
				bline = "[%[3]v]setpts=%[2]v/%[1]v*PTS[v];[0:a]atempo=%[1]v/%[2]v[a]"
			}
			filterComplex += fmt.Sprintf(bline, crdata.FrameRate, vdata.Videotrack[0].FrameRate, vpipe)

		} else {

			// Map all audio tracks if not mapped while changing fps
			for _, at := range crdata.AudioT {
				mapping += " -map 0:" + strconv.Itoa(at.StreamID)
			}
		}

		// Combine all "-filter_copmplex" filters
		vcode += fmt.Sprintf("%v -filter_complex %v %v ", fps, filterComplex, maps)

	} else {
		mapping += fmt.Sprintf(" -map 0:%v", crdata.StrID)
		for _, at := range crdata.AudioT {
			mapping += " -map 0:" + strconv.Itoa(at.StreamID)
		}
	}

	// Changes video codec
	if crdata.VideoCodec != "nochange" {
		switch crdata.VideoCodec {

		case "h264":
			vcode += fmt.Sprintf(" -c:v:%[1]v libx264 -profile:v:%[1]v main -b:v:%[1]v %[2]vk -metadata:s:v:%[1]v name=\"%[3]v\"", crdata.StrID, utils.Conf.VBW, sfname)
			break

		case "h265":
			templ := fmt.Sprintf(" -c:v:%v libx265 -x265-params \"preset=slower:me=hex:no-rect=1:no-amp=1:rd=4:aq-mode=2:", crdata.StrID)
			templ += "aq-strength=0.5:psy-rd=1.0:psy-rdoq=0.2:bframes=3:min-keyint=1\" "
			templ += fmt.Sprintf("-b:v:0 %vk -metadata:s:v:0 name=\"%v\"", utils.Conf.VBW, sfname)
			vcode += templ
		}
	} else {
		vcode += fmt.Sprintf(" -c:v:%[1]v copy -metadata:s:v:%[1]v name=\"%[2]v\"", crdata.StrID, sfname)
	}

	// Audio part ---------------------------------------------

	for _, cAt := range crdata.AudioT {
		for _, sAt := range vdata.Audiotrack {
			if cAt.StreamID == sAt.Index {

				channels := ""
				bline := " -c:a:%[1]v libfdk_aac%[4]v -b:a:%[1]v %[2]vk -metadata language=%[3]v"

				// If frame rates changed do not map
				if !(crdata.FrameRate != vdata.Videotrack[0].FrameRate) {
					mapping += fmt.Sprintf(" -map 0:%v", cAt.StreamID)
				}

				// Change layout to stereo or mono
				if cAt.Channels != sAt.Channels {
					switch cAt.Channels {

					case 2:
						channels = " -ac 2"

					case 1:
						channels = " -ac 1"
					}
				}

				// Change audio codec to aac
				if cAt.AtCodec != sAt.CodecName {
					acode += fmt.Sprintf(bline, cAt.StreamID, cAt.Channels*64, cAt.Language, channels)

				} else {
					acode += channels
				}
			}
		}
	}

	// Subtitle part ------------------------------------------

	for _, st := range crdata.SubtitleT {
		scode += fmt.Sprintf(" -c:s:%[1]v copy -metadata:s:s:%[1]v language=%[2]v", st.StreamID, st.Language)
	}

	cmd = fmt.Sprintf("ffmpeg -i %v %v %v %v %v %v -async 1 -vsync 1 %v", sf, debugIntr, acode, vcode, scode, mapping, df)

	return cmd
}
