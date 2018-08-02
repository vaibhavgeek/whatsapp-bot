package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"os"
	"strings"
	"time"
	. "github.com/mlabouardy/dialogflow-go-client"
	apiai "github.com/mlabouardy/dialogflow-go-client/models"
	"log"
)

var wac, err = whatsapp.NewConn(5 * time.Second)

type waHandler struct{}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (*waHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "error occoured: %v", err)
}

//Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (*waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.FromMe == false	{
		var result = GetResponse(message.Text)
		fmt.Printf("Info : %v \n" , message.Info)
		fmt.Printf("Message: %v \n", message.Text)
		if result.Fulfillment.Speech != "" {
			msg := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: message.Info.RemoteJid,
			},
				Text: result.Fulfillment.Speech,
			}
			wac.Send(msg)

		}
		fmt.Printf("Response : %v \n\n" , result.Fulfillment.Speech)
	}

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
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	//Add handler
	wac.AddHandler(&waHandler{})

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	<-time.After(25 * time.Hour)
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open("temp.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create("temp.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
func GetResponse(input string) apiai.Result {
	err, client := NewDialogFlowClient(apiai.Options{
		AccessToken: "6a988106a0cb4d19afd53e4c00769bda",
	})
	if err != nil {
		log.Fatal(err)
	}

	query := apiai.Query{
		Query: input,
	}
	resp, err := client.QueryFindRequest(query)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Result
}
