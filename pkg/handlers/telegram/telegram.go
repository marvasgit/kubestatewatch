package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/marvasgit/kubestatewatch/config"
	"github.com/marvasgit/kubestatewatch/pkg/event"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Telegram struct {
	httpClient *http.Client

	url             string
	chatID          int64
	messageThreadId int64
}

type MessageRequest struct {
	ChatID          int64  `json:"chat_id"`
	MessageThreadId int64  `json:"message_thread_id,omitempty"`
	Text            string `json:"text"`
	ParseMode       string `json:"parse_mode,omitempty"`
}

var statusToEmoji = map[string]string{
	"Normal":  "✅",
	"Warning": "⚠️",
	"Danger":  "❗️",
}

const messageTmpl = `%s %s %s
<code>%s</code> → <code>%s</code>

%s
`

func (r *Telegram) Init(c *config.Config) error {
	r.httpClient = http.DefaultClient
	r.httpClient.Timeout = time.Second * 30

	token := c.Handler.Telegram.Token
	chatID := c.Handler.Telegram.ChatID
	messageThreadID := c.Handler.Telegram.MessageThreadID

	if token == "" {
		token = os.Getenv("KW_TELEGRAM_TOKEN")
	}

	var err error
	if chatID == 0 {
		chatID, err = strconv.ParseInt(os.Getenv("KW_TELEGRAM_CHAT_ID"), 10, 64)
		if err != nil {
			return fmt.Errorf("telegram chat id must be an integer: %v", err)
		}
	}

	if messageThreadID == 0 {
		raw := os.Getenv("KW_TELEGRAM_MESSAGE_THREAD_ID")
		if raw != "" {
			messageThreadID, err = strconv.ParseInt(os.Getenv("KW_TELEGRAM_MESSAGE_THREAD_ID"), 10, 64)
			if err != nil {
				return fmt.Errorf("telegram message thread id must be an integer: %v", err)
			}
		}
	}
	if token == "" || chatID == 0 {
		return fmt.Errorf("telegram token and chat_id must be present")
	}

	r.url = fmt.Sprintf("https://api.telegram.org/bot%s", token)
	r.chatID = chatID
	r.messageThreadId = messageThreadID

	return r.makeRequest("/getMe", nil)
}

func makeMessageText(e event.StatemonitorEvent) string {
	diffs := ""
	if len(e.Diff) != 0 {
		var diffsValues = make([]string, 0)
		for idx, op := range e.Diff {
			if op.Value == nil {
				diffsValues = append(diffsValues, fmt.Sprintf("<b>%s</b> in %s\n", op.Type, op.Path))
			} else {
				diffsValues = append(diffsValues, fmt.Sprintf("<b>%s</b> in %s:\n    <code>%v</code>\n", op.Type, op.Path, op.Value))
			}

			// Pass only 3 diffs, cause telegram blocks too big messages
			if idx >= 3 {
				diffsValues = append(diffsValues, "<diff trimmed>")
				break
			}
		}

		diffs = fmt.Sprintf("<blockquote>%s</blockquote>", strings.Join(diffsValues, "\n"))
	}

	return fmt.Sprintf(
		messageTmpl,
		statusToEmoji[e.Status],
		e.Kind,
		strings.ToLower(e.Reason),
		e.Namespace,
		e.Name,
		diffs,
	)
}

func (r *Telegram) Handle(e event.StatemonitorEvent) {
	payload := MessageRequest{
		ChatID:          r.chatID,
		MessageThreadId: r.messageThreadId,
		Text:            makeMessageText(e),
		ParseMode:       "html",
	}

	err := r.makeRequest("/sendMessage", payload)
	if err != nil {
		logrus.Errorf("Failed to make telegram request, attempting with disabled parse mode: %v\n", err)

		payload.ParseMode = ""
		err = r.makeRequest("/sendMessage", payload)
		if err != nil {
			logrus.Errorf("Failed to make telegram retried request: %v\n", err)
			return
		}
		return
	}

	logrus.Printf("Message successfully sent to channel %d", r.chatID)
}

func (r *Telegram) makeRequest(action string, body any) error {
	var (
		resp *http.Response
		err  error
	)
	if body == nil {
		resp, err = r.httpClient.Post(r.url+action, "application/json", nil)
	} else {
		byteBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %s", err)
		}
		resp, err = r.httpClient.Post(r.url+action, "application/json", bytes.NewBuffer(byteBody))
	}
	if err != nil {
		return fmt.Errorf("failed to do request: %s", err)
	}

	return checkRespOk(resp)
}

func checkRespOk(resp *http.Response) error {
	type TGResponse struct {
		Ok          bool   `json:"ok"`
		Description string `json:"description"`
	}

	var (
		responseBody TGResponse
		err          error
	)
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}

	if !responseBody.Ok {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("telegram not found, most likely that you passed wrong bot-token: %d, %s", resp.StatusCode, responseBody.Description)
		}
		return fmt.Errorf("telegram error, %d, %s", resp.StatusCode, responseBody.Description)
	}

	return nil
}
