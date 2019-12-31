package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/config"
	"log"
	"net/http"
	"strconv"
)

type MqCredentialApiRes struct {
	Id       string `json:"id"`
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func ApiMqCredential(w http.ResponseWriter, r *http.Request) {
	authOk, device := DeviceApiCheckAuth(w, r)
	if !authOk {
		w.Write([]byte{})
		return
	}
	credentials := MqCredentialApiRes{
		Id:   "device" + strconv.Itoa(int(device.Id)),
		Host: config.MQTT_EXTERNAL_SERVER_HOST,
		Port: uint(config.MQTT_EXTERNAL_SERVER_PORT),
		// tasktab:device-1
		Username: config.MQTT_VHOST + ":tasktab-device-" + strconv.Itoa(int(device.Id)),
		Password: device.Token,
	}
	if config.DEV_MODE {
		credentials.Username = "guest"
		credentials.Password = "guest"
		log.Printf("Device getting credentials: %v", credentials)
	}

	// prepare JSON
	res, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
