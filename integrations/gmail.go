package integrations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"gitlab.systemz.pl/systemz/tasktab/types"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func GmailGetInboxMessages(credentials types.Credentials) *gmail.ListMessagesResponse {
	srv := GmailAuth(credentials.Token)
	r, err := srv.Users.Messages.List("me").LabelIds("INBOX").Do()
	if err != nil {
		log.Printf("%v", err)
	}
	return r
}

func GmailGetMessage(credentials types.Credentials, msgId string) *gmail.Message {
	srv := GmailAuth(credentials.Token)
	r, err := srv.Users.Messages.Get("me", msgId).Do()
	if err != nil {
		log.Printf("%v", err)
	}
	return r
}

func GmailMarkMessageAsDone(credentials types.Credentials, msgId string) error {
	srv := GmailAuth(credentials.Token)
	labels := []string{"UNREAD", "INBOX"}
	_, err := srv.Users.Messages.Modify("me", msgId, &gmail.ModifyMessageRequest{RemoveLabelIds: labels}).Do()
	log.Printf("%v", err)
	return err
}

func GmailMarkMessageAsToDo(credentials types.Credentials, msgId string) error {
	srv := GmailAuth(credentials.Token)
	labels := []string{"INBOX"}
	_, err := srv.Users.Messages.Modify("me", msgId, &gmail.ModifyMessageRequest{AddLabelIds: labels}).Do()
	log.Printf("%v", err)
	return err
}

func GmailAuth(token string) *gmail.Service {
	ctx := context.Background()

	// get server side app secret
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// prepare client side app token
	t := &oauth2.Token{}
	err = json.Unmarshal([]byte(token), t)
	client := config.Client(ctx, t)

	// create instanc of Gmail API Client
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	return srv
}

func GmailGetNewTokenStep1() string {
	// get server side app secret
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	return authURL
}

func GmailGetNewTokenStep2(code string) *oauth2.Token {
	// get server side app secret
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}
