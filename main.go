package main

import (
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

	http.ListenAndServe(":8080", nil)
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
}

func (s *GitHubEventService) processForkEvent(event *github.ForkEvent) {
	msg := fmt.Sprintf("%s has forked %s\nYou can access to %s", *event.Sender.Name, *event.Repo.Name, *event.Forkee.CloneURL)
	s.teleClient.SendMessage(msg)
}
