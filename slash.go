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

// Attachment represents a slack message attachment
// Currently, only images are support
type Attachment struct {
	ImageUrl string `json:"image_url"`
	Text     string `json:"text"`
}

// Response represents a response to slash command
type Response struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []Attachment `json:"attachments"`
}

// Handler function for repsonding to slack commands
type HandlerFunc func(req *Request) (*Response, error)

// Returns a new slack request from the given HTTP request
// TODO: error handling?
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

// Adds a new handler for the given command. The path and token should match the values
// set via Slack for the command you wish to handle
func HandleFunc(path, token string, h HandlerFunc) {
	http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		slackReq := newRequestFromHttpRequest(req)
		if slackReq.Token != token {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}
		log.Printf("Handling %s %s (channel=%s, user=%s)\n", path, slackReq.Text, slackReq.ChannelName, slackReq.UserName)

		resp, err := h(slackReq)
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

// Serves all slash commands over HTTP. Note that Slack requires Slash commands to run
// over HTTPS, so generally ListenAndServeTLS should be used
func ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}

// ListenAndServeTLS serves all slash commands over HTTPS
func ListenAndServeTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}
