package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "net/http"
    "bufio"
    "os/exec"
    b64 "encoding/base64"
)

var conn2 *websocket.Conn = nil

func handle2(){
    //cmd := exec.Command("ffmpeg","-f", "dshow", "-i", "video=HP HD Camera", "-preset", "ultrafast", "-f", "mjpeg", "pipe:1")
    cmd := exec.Command("ffmpeg","-f", "v4l2", "-framerate", "10", "-video_size", "380x200", "-i", "/dev/video0", "-preset", "ultrafast", "-s", "800x600", "-f", "mjpeg", "pipe:1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	
    cmd.Start()
	chunk := make([]byte, 40*1024)
	for {
        nr, _ := stdout.Read(chunk)
        if nr > 0 {
            validData := chunk[:nr]
            sEnc := b64.StdEncoding.EncodeToString(validData)
            conn2.WriteMessage(websocket.TextMessage, []byte(sEnc))
        }
	}
	cmd.Wait()
}

func handle(){
    //cmd := exec.Command("ffmpeg","-f", "dshow", "-i", "video=\"HP", "HD", "Camera\"", "-preset", "ultrafast", "-f", "mjpeg", "pipe:1")
    cmd := exec.Command("ffmpeg","-f", "dshow", "-i", "video=HP HD Camera", "-preset", "ultrafast", "-f", "mjpeg", "pipe:1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
    cmd.Start()
	scanner := bufio.NewScanner(stdout)
	
	for scanner.Scan() {
        m := scanner.Bytes()
        sEnc := b64.StdEncoding.EncodeToString(m)
        conn2.WriteMessage(websocket.TextMessage, []byte(sEnc))
        
        //time.Sleep(500)
	}
	cmd.Wait()
}


func main() {
    r := gin.Default()
    r.LoadHTMLFiles("assets/index.html")

    r.GET("/", func(c *gin.Context) {
        c.HTML(200, "index.html", nil)
    })

    r.GET("/ws", func(c *gin.Context) {
        wshandler(c.Writer, c.Request)
    })

    r.Run("0.0.0.0:8089")
}

var wsupgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
    conn, err := wsupgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Failed to set websocket upgrade: %+v", err)
        return
    }
    conn2 = conn

    go handle2()
}