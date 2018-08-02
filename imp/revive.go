
package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"os"
	"strings"
	"time"
)

type waHandler struct{}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (*waHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "error occoured: %v", err)
}

//Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (*waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Printf("time:%v from:%v\n\t message:%v\n\n Lot More Info: %v\n\n\n ", message.Info.Timestamp, message.Info.RemoteJid, message.Text, message.Info)
}

//Example for media handling. Video, Audio, Document are also possible in the same way
func (*waHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	data, err := message.Download()
	if err != nil {
		return
	}
	filename := fmt.Sprintf("%v/%v.%v", os.TempDir(), message.Info.Id, strings.Split(message.Type, "/")[1])
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	fmt.Printf("%v %v\n\timage reveived, saved at:%v\n", message.Info.Timestamp, message.Info.RemoteJid, filename)
}

func main() {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}
