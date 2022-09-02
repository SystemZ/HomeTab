package service

import (
	"time"

	"github.com/gregdel/pushover"
	"github.com/systemz/hometab/internal/config"
	"github.com/systemz/hometab/internal/model"
)

type Notification struct {
	Id        uint   `json:"id"`
	SessionId uint   `json:"sessionId"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Msg       string `json:"msg"`
}

func SendGenericNotificationToAllDevices(title string, body string, ignoreUsers []int) {
	// get DB info
	var devices []model.Device
	model.DB.Find(&devices)

	// send message to each device
	for _, device := range devices {
		send := true
		// ignore specified users
		for _, userId := range ignoreUsers {
			if uint(userId) == device.UserId {
				send = false
				break
			}
		}
		if send {
			SendGenericNotification(title, body, device)
		}
	}
}

func SendGenericNotification(title string, body string, device model.Device) {
	// summary of session
	msg := Notification{
		Id:        0,
		SessionId: 0,
		Type:      "showNotification",
		Title:     title,
		Msg:       body,
	}
	//SendPushyMe(msg, device)
	PushoverNotification(msg, device)
}

func SendCounterNotification(start bool, sourceUser model.User, counterId uint, sessionId uint, sessionTaken string) {
	// get DB info
	var counter model.Counter
	model.DB.Where(model.Counter{Id: counterId}).First(&counter)
	var devices []model.Device
	model.DB.Find(&devices)

	// send message to each device
	for _, device := range devices {
		currentTime := time.Now()
		msgTitle := sourceUser.Username + " @ " + counter.Name
		msgBody := "Counting from " + currentTime.Format("15:04:05 02-01-2006")
		// add or remove notification from device
		msgType := "startNotification"
		if !start {
			// sum up this session as separate notification
			SendGenericNotification(msgTitle, "Session: "+sessionTaken+"\nEnded "+currentTime.Format("15:04:05 02-01-2006"), device)

			// send empty notification to remove ongoing notification
			msgType = "stopNotification"
			msgTitle = ""
			msgBody = ""
		}

		// finally craft queue message
		msg := Notification{
			Id:        counterId,
			SessionId: sessionId,
			Type:      msgType,
			Title:     msgTitle,
			Msg:       msgBody,
		}
		//SendPushyMe(msg, device)
		PushoverNotification(msg, device)
	}
}

// TODO insert switch case for PushyMe notifications
//type PushyMeReq struct {
//	To           string       `json:"to"`
//	Notification Notification `json:"data"`
//}

/*
func EncryptAES(key []byte, plaintext string) string {
	// create cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("encrypting msg err: %v", err)
		return ""
	}

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, []byte(plaintext))
	// return hex string
	return hex.EncodeToString(out)
}
*/

// send push notification to device via pushy.me service
//func SendPushyMe(msg Notification, device model.Device) (err error) {
//	if len(device.TokenPush) < 1 {
//		return errors.New("wrong push token for pushy.me")
//	}
//	// https://stackoverflow.com/questions/40123319/easy-way-to-encrypt-decrypt-string-in-android
//	// use first 32 characters from device token as basic encryption measure in transit
//	//encryptKey := []byte(device.Token[0:32])
//	//msg.Title = EncryptAES(encryptKey, msg.Title)
//
//	log.Printf("Sending pushy.me msg to %v", device.Name)
//	pushReqRaw := PushyMeReq{
//		To:           device.TokenPush,
//		Notification: msg,
//	}
//	pushReq, err := json.Marshal(&pushReqRaw)
//	if err != nil {
//		log.Printf("failed preparing msg for push notification: %v", err)
//		return err
//	}
//
//	c := &http.Client{}
//	reqUrl := "https://api.pushy.me/push?api_key=" + config.PUSHY_ME_SECRET
//	r, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(pushReq))
//	if err != nil {
//		log.Printf("fail when creating request for api.pushy.me: %v", err)
//		return err
//	}
//	r.Header.Set("Content-Type", "application/json")
//	res, err := c.Do(r)
//	if err != nil {
//		log.Printf("fail when sending request to api.pushy.me: %v", err)
//		return err
//	}
//	defer res.Body.Close()
//	log.Printf("HTTP %v @ api.pushy.me", res.StatusCode)
//	return nil
//}

func PushoverNotification(msg Notification, device model.Device) error {
	app := pushover.New(config.PUSHOVER_API_TOKEN)
	recipient := pushover.NewRecipient(config.PUSHOVER_USER_KEY)
	message := &pushover.Message{
		Message:    msg.Msg,
		Title:      msg.Title,
		Priority:   pushover.PriorityNormal,
		DeviceName: device.Name,
		Sound:      pushover.SoundGamelan,
	}

	// Send the message to the recipient
	_, err := app.SendMessage(message, recipient)
	return err
}
