# Golang - FFMPEG stream to web page via websocket

This small software is an example that can be used in order to stream ffmpeg inputs to webpage via websocket in go programming language.
ffmpeg frames are sended to the client as base64 strings.
The client displays frames in canvas. 
It uses [gin-gonic](https://github.com/gin-gonic/gin) and [gorilla websocket](https://github.com/gorilla/websocket).