package main


import (
	"github.com/vivekvasvani/Slack-Bot/server"
)



func main() {
	wait := make(chan struct{})
	server.NewServer()
	<-wait
}






