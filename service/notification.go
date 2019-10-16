package service

import (
	"gitlab.com/systemz/tasktab/model"
	"gitlab.com/systemz/tasktab/queue"
	"log"
	"strconv"
)

func SendCounterNotification(start bool, sourceUser model.User, counterId uint) {
	// get DB info
	var counter model.Counter
	model.DB.Where(model.Counter{Id: counterId}).First(&counter)
	var devices []model.Device
	model.DB.Find(&devices)

	// send message to each device
	for _, device := range devices {
		msgKey := "device" + strconv.Itoa(int(device.Id))
		msgBody := counter.Name
		// for devices of other users we craft other message
		if device.UserId != sourceUser.Id {
			msgBody = sourceUser.Username + " @ " + counter.Name
		}
		// add or remove notification from device
		msgType := "startNotification"
		if !start {
			msgType = "stopNotification"
		}

		// finally craft queue message
		msg := queue.Notification{
			Id:   counterId,
			Type: msgType,
			Msg:  msgBody,
		}
		log.Printf("Sending push msg to %v", device.Name)
		queue.SendNotification(msg, msgKey)
	}
}
