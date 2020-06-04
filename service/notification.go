package service

import (
	"bytes"
	"encoding/json"
	"gitlab.com/systemz/tasktab/config"
	"gitlab.com/systemz/tasktab/model"
	"gitlab.com/systemz/tasktab/queue"
	"log"
	"net/http"
	"strconv"
)

func SendCounterNotification(start bool, sourceUser model.User, counterId uint, sessionId uint) {
	// get DB info
	var counter model.Counter
	model.DB.Where(model.Counter{Id: counterId}).First(&counter)
	var devices []model.Device
	model.DB.Find(&devices)

	// send message to each device
	for _, device := range devices {
		msgKey := "device" + strconv.Itoa(int(device.Id))
		msgTitle := sourceUser.Username + " @ " + counter.Name
		msgBody := "Counting..."
		// add or remove notification from device
		msgType := "startNotification"
		if !start {
			msgType = "stopNotification"
			msgTitle = ""
			msgBody = ""
		}

		// finally craft queue message
		msg := queue.Notification{
			Id:        counterId,
			SessionId: sessionId,
			Type:      msgType,
			Title:     msgTitle,
			Msg:       msgBody,
		}
		log.Printf("Sending AMQP/MQTT push msg to %v", device.Name)
		queue.SendNotification(msg, msgKey)

		pushReqRaw := PushyMeReq{
			To:           device.TokenPush,
			Notification: msg,
		}
		pushReq, err := json.Marshal(&pushReqRaw)
		if err != nil {
			log.Printf("failed preparing msg for push notification: %v", err)
			return
		}
		log.Printf("Sending pushy.me msg to %v", device.Name)
		SendPushyMe(pushReq)
	}
}

type PushyMeReq struct {
	To           string             `json:"to"`
	Notification queue.Notification `json:"data"`
}

func SendPushyMe(body []byte) error {
	c := &http.Client{}
	reqUrl := "https://api.pushy.me/push?api_key=" + config.PUSHY_ME_SECRET
	r, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("fail when creating request for api.pushy.me: %v", err)
	}
	r.Header.Set("Content-Type", "application/json")
	res, err := c.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Printf("HTTP %v @ api.pushy.me", res.StatusCode)
	return nil
}
