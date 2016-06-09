// Copyright 2016 Matthew Fonda. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// slash package wraps net/http for simple creation of slack commands
// For more info on slash commands, see https://api.slack.com/slash-commands
package slash

import (
	"encoding/json"
	"log"
	"net/http"
)

// Attachment represents a slack message attachment
// Currently, only images are support
type Attachment struct {
	ImageUrl string `json:"image_url"`
	Text     string `json:"text"`
}

// Request represents an incoming slash command request
type Request struct {
	Token       string
	TeamId      string
	TeamDomain  string
	ChannelId   string
	ChannelName string
	UserId      string
	UserName    string
	Command     string
	Text        string
	ResponseUrl string
}

// Response represents a response to slash command
type Response struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

// Command represents a slash command and its handler
type Command struct {
	Command string
	Token   string
	Handler Handler
}

// Handler function for repsonding to slack commands
type Handler func(req *Request) (*Response, error)

// NewCommand returns a new command for the given command (e.g. "/foo"), token, and handler
func NewCommand(command, token string, h Handler) *Command {
	c := &Command{}
	c.Command = command
	c.Token = token
	c.Handler = h
	return c
}

func newRequestFromHttpRequest(req *http.Request) *Request {
	r := &Request{}
	r.Token = req.FormValue("token")
	r.TeamId = req.FormValue("team_id")
	r.TeamDomain = req.FormValue("team_domain")
	r.ChannelId = req.FormValue("channel_id")
	r.ChannelName = req.FormValue("channel_name")
	r.UserId = req.FormValue("user_id")
	r.UserName = req.FormValue("user_name")
	r.Command = req.FormValue("command")
	r.Text = req.FormValue("text")
	r.ResponseUrl = req.FormValue("response_url")
	return r
}

// NewInChannelResponse returns a Response to be sent to an entire channel
// in response to a slash command
func NewInChannelResponse(text string, attachments []Attachment) *Response {
	r := &Response{}
	r.ResponseType = "in_channel"
	r.Text = text
	r.Attachments = attachments
	return r
}

func HandleCommand(c *Command) {
	http.HandleFunc(c.Command, func(w http.ResponseWriter, req *http.Request) {
		slackReq := newRequestFromHttpRequest(req)
		if slackReq.Token != c.Token {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		log.Printf("Handling %s %s (channel=%s, user=%s)\n", c.Command, slackReq.Text, slackReq.ChannelName, slackReq.UserName)

		resp, err := c.Handler(slackReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(json)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// ListenAndServeTLS serves all slash commands
func ListenAndServeTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}
