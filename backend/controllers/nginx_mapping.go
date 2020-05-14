package controllers

//"github.com/Dzionys/video-platform/backend/utils"

var tcvidpath = "/home/dzionys/go/src/github.com/Dzionys/video-platform/vod/transcoded/%v"

// NginxMappingHandler ...
// func NginxMappingHandler(w http.ResponseWriter, r *http.Request) {

// 	var sqncs models.Sequences

// 	vars := mux.Vars(r)

// 	if filepath.Ext(vars["name"]) == ".mp4" {
// 		//err := db.IsExist("Video", vars["name"])
// 		if err := utils.DB.Where("Name = ?", vars["name"]).First(&models.Video{}).Error; err != nil {
// 			log.Println(err)
// 			w.WriteHeader(404)
// 			return
// 		}

// 		temp := models.Clip{
// 			"source",
// 			fmt.Sprintf(tcvidpath, vars["name"]),
// 		}
// 		var tempclip models.Clips
// 		tempclip.Clips = append(tempclip.Clips, temp)
// 		sqncs.Sequences = append(sqncs.Sequences, tempclip)

// 	} else {
// 		names, err := db.GetAllStreamVideos(vars["name"])
// 		if err != nil {
// 			log.Println(err)
// 			w.WriteHeader(404)
// 			return
// 		}

// 		for _, n := range names {
// 			temp := vd.Clip{
// 				"source",
// 				fmt.Sprintf(tcvidpath, n),
// 			}
// 			var tempclip vd.Clips
// 			tempclip.Clips = append(tempclip.Clips, temp)
// 			sqncs.Sequences = append(sqncs.Sequences, tempclip)
// 		}
// 	}

// 	j, err := json.Marshal(sqncs)
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(500)
// 		return
// 	}

// 	w.Write(j)
// }

// func getAllStreamNames(sname string) ([]string, error) {
// 	models.Stream
// }
