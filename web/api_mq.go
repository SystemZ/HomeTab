package web

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/systemz/tasktab/config"
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
	authOk, device := ApiCheckAuth(w, r)
	if !authOk {
		w.Write([]byte{})
		return
	}
	deviceMqId := uuid.NewV4().String()
	credentials := MqCredentialApiRes{
		Id:   deviceMqId,
		Host: config.MQTT_EXTERNAL_SERVER_HOST,
		Port: uint(config.MQTT_EXTERNAL_SERVER_PORT),
		// tasktab:device-1
		Username: config.MQTT_VHOST + ":tasktab-device-" + strconv.Itoa(int(device.Id)),
		Password: device.Token,
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
