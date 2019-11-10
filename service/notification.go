package service

import (
	"gitlab.com/systemz/tasktab/model"
	"gitlab.com/systemz/tasktab/queue"
	"log"
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
		log.Printf("Sending push msg to %v", device.Name)
		queue.SendNotification(msg, msgKey)
	}
}
