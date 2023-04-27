// Package chunktest provides a test server for testing the chunk package.
package chunktest

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/slack-go/slack"

	"github.com/rusq/slackdump/v2/internal/chunk"
)

var lg = log.New(os.Stderr, "chunktest: ", log.LstdFlags)

// Server is a test server for testing the chunk package, that serves API
// from a single chunk file.
type Server struct {
	baseServer
	p *chunk.Player
}

// NewServer returns a new Server, it requires the chunk file handle in rs, and
// an ID of the user that will be returned by AuthTest in currentUserID.
func NewServer(rs io.ReadSeeker, currentUserID string) *Server {
	p, err := chunk.NewPlayer(rs)
	if err != nil {
		panic(err)
	}
	return &Server{
		baseServer: baseServer{Server: httptest.NewServer(router(p, currentUserID))},
		p:          p,
	}
}

type GetConversationRepliesResponse struct {
	slack.SlackResponse
	HasMore          bool             `json:"has_more"`
	ResponseMetaData responseMetaData `json:"response_metadata"`
	Messages         []slack.Message  `json:"messages"`
}

type responseMetaData struct {
	NextCursor string `json:"next_cursor"`
}

func router(p *chunk.Player, userID string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/api/auth.test", authHandler{userID})

	mux.HandleFunc("/api/conversations.info", handleConversationInfo(p))
	mux.HandleFunc("/api/conversations.history", handleConversationsHistory(p))
	mux.HandleFunc("/api/conversations.replies", handleConversationsReplies(p))
	mux.HandleFunc("/api/conversations.list", handleConversationsList(p))
	mux.HandleFunc("/api/users.list", handleUsersList(p))
	return mux
}
