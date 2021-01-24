package server

import (
	"bytes"
	"encoding/json"
	// gitlabWh "github.com/go-playground/webhooks/gitlab"
	"github.com/systemz/hometab/internal/model"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func WebhookGitlab(w http.ResponseWriter, r *http.Request) {
	//enforce POST only
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("X-Gitlab-Event") == "Issue Hook" {
		GitlabWhIssueToDiscord(w, r)
	} else if r.Header.Get("X-Gitlab-Event") == "Note Hook" {
		// new comment on commit, merge_request, issue, snippet
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

/*
func WebhookGitlab(w http.ResponseWriter, r *http.Request) {
	//enforce POST only
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hook, _ := gitlabWh.New()
	payload, err := hook.Parse(r, gitlabWh.PushEvents)
	if err != nil {
		if err == gitlabWh.ErrEventNotFound {
			log.Printf("%v", "Event from gitlab not parsed")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	switch payload.(type) {
	case gitlabWh.PushEventPayload:
		data := payload.(gitlabWh.PushEventPayload)
		log.Printf("Push! ProjectID: %v", data.ProjectID)
	}


	w.WriteHeader(http.StatusOK)
	return
}
*/

func GitlabWhIssueToDiscord(w http.ResponseWriter, r *http.Request) {
	//get raw JSON
	EventRaw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//parse from JSON
	var webhook GitlabIssueWebhook
	err = json.Unmarshal(EventRaw, &webhook)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO check auth

	// check what to do with legit issue
	var whReceiver model.WebhookReceiver
	model.DB.Where(model.WebhookReceiver{ProjectId: uint(webhook.Project.ID)}).First(&whReceiver)
	actionId := whReceiver.Id
	if actionId < 1 {
		log.Print("Webhook - project with ID not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get action to take
	var whAction model.WebhookAction
	model.DB.Where(model.WebhookAction{Id: actionId}).First(&whAction)
	if whAction.Id < 0 {
		log.Print("Webhook - action with ID not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	discordWhUrl := whAction.Extra

	// execute action
	bodyRaw := DiscordWebhook{
		//Content:  "test",
		Username: "TaskTab",
		Embeds: []DiscordWhEmbed{
			{
				//Description: "Issue created/updated/closed/reopened",
				Title: "#" + strconv.FormatInt(webhook.ObjectAttributes.Iid, 10) + " " + webhook.ObjectAttributes.Title,
				Url:   webhook.ObjectAttributes.URL,
				Color: 9287168,
				Author: DiscordWhEmbedAuthor{
					Name:     webhook.User.Name,
					Url:      "https://gitlab.com/" + webhook.User.Username,
					Icon_url: webhook.User.AvatarURL,
				},
				Fields: []DiscordWhFields{
					{
						Name:  "State",
						Value: webhook.ObjectAttributes.State,
						// Value:  webhook.ObjectAttributes.Action,
						Inline: true,
					},
				},
			},
		},
	}
	SendDiscordWebhook(bodyRaw, discordWhUrl, w)
}

// https://discordapp.com/developers/docs/resources/webhook
// https://discordapp.com/developers/docs/resources/channel#embed-object-embed-author-structure
type DiscordWebhook struct {
	Content  string           `json:"content"`
	Username string           `json:"username"`
	Embeds   []DiscordWhEmbed `json:"embeds"`
}
type DiscordWhEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	// https://www.spycolor.com/
	Color  int                  `json:"color"`
	Author DiscordWhEmbedAuthor `json:"author"`
	Fields []DiscordWhFields    `json:"fields"`

	/*
		timestamp?	ISO8601 timestamp	timestamp of embed content
		footer?	embed footer object	footer information
		image?	embed image object	image information
		thumbnail?	embed thumbnail object	thumbnail information
		video?	embed video object	video information
		provider?	embed provider object	provider information
		fields?	array of embed field objects	fields information
	*/
}
type DiscordWhFields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type DiscordWhEmbedAuthor struct {
	Name           string `json:"name"`
	Url            string `json:"url"`
	Icon_url       string `json:"icon_url"`
	Proxy_icon_url string `json:"proxy_icon_url"`
}

func SendDiscordWebhook(bodyRaw DiscordWebhook, discordWhUrl string, w http.ResponseWriter) {
	// prepare JSON
	bodyReady, err := json.Marshal(bodyRaw)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", discordWhUrl, bytes.NewBuffer(bodyReady))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//bodyRes, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(bodyRes))
}

type GitlabIssueWebhook struct {
	ObjectKind       string           `json:"object_kind"`
	User             Assignee         `json:"user"`
	Project          Project          `json:"project"`
	Repository       Repository       `json:"repository"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
	Assignees        []Assignee       `json:"assignees"`
	Assignee         Assignee         `json:"assignee"`
	Labels           []Label          `json:"labels"`
	Changes          Changes          `json:"changes"`
}

type Assignee struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type Changes struct {
	UpdatedByID UpdatedByID `json:"updated_by_id"`
	UpdatedAt   UpdatedAt   `json:"updated_at"`
	Labels      Labels      `json:"labels"`
}

type Labels struct {
	Previous []Label `json:"previous"`
	Current  []Label `json:"current"`
}

type Label struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Color       string `json:"color"`
	ProjectID   int64  `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Template    bool   `json:"template"`
	Description string `json:"description"`
	Type        string `json:"type"`
	GroupID     int64  `json:"group_id"`
}

type UpdatedAt struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type UpdatedByID struct {
	Previous interface{} `json:"previous"`
	Current  int64       `json:"current"`
}

type ObjectAttributes struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	AssigneeIDS []int64     `json:"assignee_ids"`
	AssigneeID  int64       `json:"assignee_id"`
	AuthorID    int64       `json:"author_id"`
	ProjectID   int64       `json:"project_id"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Position    int64       `json:"position"`
	BranchName  interface{} `json:"branch_name"`
	Description string      `json:"description"`
	MilestoneID interface{} `json:"milestone_id"`
	State       string      `json:"state"`
	Iid         int64       `json:"iid"`
	URL         string      `json:"url"`
	Action      string      `json:"action"`
}

type Project struct {
	ID                int64       `json:"id"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	WebURL            string      `json:"web_url"`
	AvatarURL         interface{} `json:"avatar_url"`
	GitSSHURL         string      `json:"git_ssh_url"`
	GitHTTPURL        string      `json:"git_http_url"`
	Namespace         string      `json:"namespace"`
	VisibilityLevel   int64       `json:"visibility_level"`
	PathWithNamespace string      `json:"path_with_namespace"`
	DefaultBranch     string      `json:"default_branch"`
	Homepage          string      `json:"homepage"`
	URL               string      `json:"url"`
	SSHURL            string      `json:"ssh_url"`
	HTTPURL           string      `json:"http_url"`
}

type Repository struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}
