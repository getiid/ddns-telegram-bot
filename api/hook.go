package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type IP struct {
	Result string `json:"result"`
	Addr   string `json:"addr"`
	Domain string `json:"domain"`
}

type DDNSRequest struct {
	IPv4 IP `json:"ipv4,omitempty"`
	IPv6 IP `json:"ipv6,omitempty"`
}

func HookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	var req DDNSRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Fatal("Error to parse DDNS request", err)
	}

	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Error to create a Telegram bot", err)
	}

	chatId, _ := strconv.ParseInt(r.URL.Query().Get("chatId"), 10, 64)

	text := ""

	if req.IPv4.Result != "未改变" {
		text += fmt.Sprintf("IPv4: %s\n%s\n%s\n", req.IPv4.Result, req.IPv4.Addr, req.IPv4.Domain)
	}

	if req.IPv6.Result != "未改变" {
		text += fmt.Sprintf("IPv6: %s\n%s\n%s\n", req.IPv4.Result, req.IPv6.Addr, req.IPv6.Domain)
	}

	w.Header().Add("Content-Type", "application/json")
	if text != "" {
		msg := tgbotapi.NewMessage(chatId, text)
		if _, err = bot.Send(msg); err != nil {
			fmt.Fprintf(w, `
			{
				status: "false"
				msg: "%s"
			}
			`, err)
		} else {
			fmt.Fprintf(w, `
			{
				status: "true"
			}
			`)
		}
	} else {
		fmt.Fprint(w, `
		{
			status: "false",
			msg: "IPv4 and IPv6 didn't change"
		}
		`)
	}

}
