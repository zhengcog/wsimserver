package wsimserver

import (
	"log"
	"net/http"
	"runtime"
	"wsimserver/chat"
	"wsimserver/utils"
)

func main() {
	runtime.GOMAXPROCS(4)
	log.SetFlags(log.Lshortfile)

	utils.InitRedis()

	//websocket server
	server := chat.NewServer("/server")
	go server.Listen()

	log.Fatal(http.ListenAndServe(":9090", nil))
}
