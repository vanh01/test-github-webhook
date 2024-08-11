package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v63/github"

	"github-webhooks/tele"
)

func main() {
	ghem := GitHubEventService{
		webhookSecretKey: "please update it",
		teleClient: tele.TelegramClient{
			ChatId: 0,
			ApiKey: "",
		},
	}
	http.HandleFunc("/webhook", ghem.ServeHTTP)

	log.Fatalln(http.ListenAndServe(":8081", nil))
}

type GitHubEventService struct {
	webhookSecretKey string
	teleClient       tele.TelegramClient
}

func (s *GitHubEventService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	payload, err := github.ValidatePayload(r, []byte(s.webhookSecretKey))
	if err != nil {
		log.Printf("Err to validate payload: %s\n", err.Error())
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("Err to parse webhook: %s\n", err.Error())
	}
	switch event := event.(type) {
	case *github.Commit:
		s.processCommitCommentEvent(event)
	case *github.ForkEvent:
		s.processForkEvent(event)
	}
}

func (s *GitHubEventService) processCommitCommentEvent(event *github.Commit) {
	if event == nil {
		s.teleClient.SendMessage("Event is nil")
	}
	senderName := "Someone"
	if event.Committer != nil && event.Committer.Name != nil {
		senderName = *event.Committer.Name
	}
	sha := event.GetSHA()
	cmtMsg := event.GetMessage()
	url := event.GetURL()
	msg := fmt.Sprintf("%s has commited with \"%s\". The SHA: %s\nYou can access to %s", senderName, cmtMsg, sha, url)
	s.teleClient.SendMessage(msg)
}

func (s *GitHubEventService) processForkEvent(event *github.ForkEvent) {
	if event == nil {
		s.teleClient.SendMessage("Event is nil")
	}
	json, err := json.Marshal(event)
	if err == nil {
		fmt.Println(string(json))
	}
	senderName := "Someone"
	if event.Sender != nil && event.Sender.Login != nil {
		senderName = *event.Sender.Login
	}
	repoName := "your repository"
	if event.Repo != nil && event.Repo.Name != nil {
		repoName = *event.Repo.Name
	}
	url := ""
	if event.Forkee != nil && event.Forkee.URL != nil {
		url = *event.Forkee.URL
	}
	msg := fmt.Sprintf("%s has forked %s\nYou can access to %s", senderName, repoName, url)
	s.teleClient.SendMessage(msg)
}
