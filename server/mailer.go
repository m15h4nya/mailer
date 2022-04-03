package server

import (
	"apitask/db"
	"fmt"
	"github.com/google/uuid"
	"github.com/thethanos/go-containers/containers"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func StartMailing(ch chan struct{}, dbConn *gorm.DB) {
	var mailings []db.Mailing
	for {
		<-ch
		update(&mailings, dbConn)

		heap := containers.NewHeap(func(a, b db.Mailing) bool { return b.StartDate.After(a.StartDate) })
		heap.PushSlice(mailings)

		for !heap.Empty() {
			mailing := heap.Top()
			var clients []db.Client
			dbConn.Model(&db.Client{}).Find(&clients, "tag IN ?", strings.Split(mailing.Tags, " "))

			fmt.Println(mailing.StartDate, time.Now(), mailing.StartDate.Sub(time.Now()))

			time.Sleep(mailing.StartDate.Sub(time.Now()))

			sendMsg(clients, mailing.Text, mailing.MailingID, dbConn)
			heap.Pop()
		}
	}
}

func update(mailings *[]db.Mailing, dbConn *gorm.DB) {
	dbConn.Model(db.Mailing{}).Find(mailings, "end_date > ? AND status = ?", time.Now(), false)
}

func sendMsg(clients []db.Client, text string, mailingId string, dbConn *gorm.DB) {
	client := &http.Client{}
	for _, v := range clients {
		msgId := uuid.NewString()

		reader := strings.NewReader(fmt.Sprintf(`{"id": %v,"phone": %v,"text": "%v"}`, msgId, v.Phone, text))
		req, err := http.NewRequest("POST", "https://probe.fbrq.cloud/v1", reader)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODAzNDcxODksImlzcyI6ImZhYnJpcXVlIiwibmFtZSI6Ik51cklicmFnaW1vdiJ9.jh9ZzasOMklP1ZQyfR7cU27nyQ_1UtA90ExPKPcz_80")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Message sent")

		msg := &db.Msg{
			MsgID:     msgId,
			Date:      time.Now(),
			Status:    resp.StatusCode,
			MailingID: mailingId,
			UserID:    v.ClientID,
		}
		dbConn.Create(msg)
	}

	dbConn.Model(db.Mailing{}).Where("mailing_id = ?", mailingId).Update("status", true)
}
