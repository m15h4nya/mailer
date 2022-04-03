package server

import (
	db2 "apitask/db"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	log *zap.SugaredLogger
	db  *gorm.DB
	sig chan struct{}
}

type id struct {
	ID string `json:"id"`
}

func NewHandler(log *zap.SugaredLogger, db *gorm.DB, ch chan struct{}) *Handler {
	return &Handler{log: log, db: db, sig: ch}
}

func (h *Handler) AddClient(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	client := db2.Client{}
	if err := json.Unmarshal(body, &client); err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	tx := h.db.Create(&client)
	if tx.Error != nil {
		h.log.Error(tx.Error.Error())
		http.Error(rw, tx.Error.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateClient(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	client := db2.Client{}
	if err := json.Unmarshal(body, &client); err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if tx := h.db.Model(&db2.Client{}).Where("id = ?", client.ID).Updates(client); tx.Error != nil {
		h.log.Error(tx.Error.Error())
		http.Error(rw, tx.Error.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteClient(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	client := db2.Client{}
	if err = json.Unmarshal(body, &client); err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = h.db.Delete(client).Error; err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddMailing(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	mailing := db2.Mailing{}
	if err = json.Unmarshal(body, &mailing); err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	if err = h.db.Create(&mailing).Error; err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateMailing(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	mailing := db2.Mailing{}
	if err := json.Unmarshal(body, &mailing); err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if tx := h.db.Model(&db2.Client{}).Where("id = ?", mailing.ID).Updates(mailing); tx.Error != nil {
		h.log.Error(tx.Error.Error())
		http.Error(rw, tx.Error.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	var msgSent []db2.Msg
	var msgDelivered []db2.Msg
	var msgFailed []db2.Msg

	if err := h.db.Model(db2.Msg{}).Where("status = ?", db2.SENT).Find(&msgSent).Error; err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.db.Model(db2.Msg{}).Where("status = ?", db2.DELIVERED).Find(&msgDelivered).Error; err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.db.Model(db2.Msg{}).Where("status = ?", db2.FAILED).Find(&msgFailed).Error; err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	//group status = SENT messages
	sent := fmt.Sprintf("SENT %d messages:\n", len(msgSent))
	for _, v := range msgSent {
		sent += fmt.Sprintf("msg ID: %s, mailing ID: %s\n", v.ID, v.MailingID)
	}
	//group status = DELIVERED messages
	delivered := fmt.Sprintf("DELIVERED %d messages:\n", len(msgDelivered))
	for _, v := range msgSent {
		delivered += fmt.Sprintf("msg ID: %s, mailing ID: %s\n", v.ID, v.MailingID)
	}
	//group status = FAILED messages
	failed := fmt.Sprintf("FAILED %d messages:\n", len(msgSent))
	for _, v := range msgSent {
		failed += fmt.Sprintf("msg ID: %s, mailing ID: %s\n", v.ID, v.MailingID)
	}

	if _, err := fmt.Fprintf(rw, sent+delivered+failed); err != nil {
		h.log.Info(err)
	}
}

func (h *Handler) GetMailingStatistic(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	mailingId := id{}
	if err = json.Unmarshal(body, &mailingId); err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	var msgs []db2.Msg

	if err := h.db.Model(db2.Msg{}).Where("id = ?", mailingId).Find(&msgs).Error; err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	res := fmt.Sprintf("Mailing created %d messages:\n", len(msgs))
	for _, v := range msgs {
		res += fmt.Sprintf("msg ID: %s, mailing ID: %s, user ID: %s, status: %s, date: %s\n", v.ID, v.MailingID, v.UserID, v.Status, v.Date.String())
	}
	if _, err := fmt.Fprintf(rw, res); err != nil {
		h.log.Info(err)
	}
}

func (h *Handler) DeleteMailing(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	mailing := db2.Mailing{}
	if err = json.Unmarshal(body, &mailing); err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = h.db.Delete(mailing).Error; err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) StartMailing(rw http.ResponseWriter, r *http.Request) {
	h.sig <- struct{}{}
}
