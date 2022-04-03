package db

import (
	"gorm.io/gorm"
	"time"
)

const (
	SENT = iota
	DELIVERED
	FAILED
)

type Mailing struct {
	gorm.Model
	MailingID string    `json:"id"`
	StartDate time.Time `json:"start_date"`
	Text      string    `json:"text"`
	Tags      string    `json:"tags"`
	EndDate   time.Time `json:"end_date"`
	Status    bool      `json:"status"`
}

type Client struct {
	gorm.Model
	ClientID  string `json:"id"`
	Phone     int    `json:"phone"`
	PhoneCode int    `json:"phone_code"`
	Tag       string `json:"tag"`
	TimeZone  int    `json:"time_zone"`
}

type Msg struct {
	gorm.Model
	MsgID     string    `json:"id"`
	Date      time.Time `json:"date"`
	Status    int       `json:"status"`
	MailingID string    `json:"mailing_id"`
	UserID    string    `json:"user_id"`
}
