package public

import (
	"time"
	"sync"
)

var g_scene_id int64
var g_scene_id_mutex sync.Mutex

func GetSceneID()(int64){
	g_scene_id_mutex.Lock()
	nowNumber:=time.Now().Unix()
	if nowNumber>g_scene_id {
		g_scene_id=nowNumber
	} else {
		g_scene_id+=1
	}
	g_scene_id_mutex.Unlock()
	return g_scene_id
}