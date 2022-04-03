// Package classification Products API
//
// Documentation for Products API
//
// Schemes: http
// BasePath: /
// Version 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package docs

import "apitask/db"

// swagger:route PATCH /client/{id} client UpdateClient
// Updates client in db
// responses:
// 200
// 500

//swagger:parameters UpdateClient
type ClientIdParam struct {
	// in:path
	ClientId string `json:"id"`
}

// swagger:parameters UpdateClient
type ClientParam struct {
	// in:body
	Client db.Client `json:"client"`
}

// swagger:route PATCH /mailing/{id} mailing UpdateMailing
// Updates mailing in db
// responses:
// 200
// 500

// swagger:parameters UpdateMailing
type MailingIdParam struct {
	// in:path
	MailingId string `json:"id"`
}

// swagger:parameters UpdateMailing
type MailingParam struct {
	// in:body
	Mailing db.Mailing `json:"mailing"`
}

// swagger:route GET /statistics statistics GetStatistics
// Returns statistics about sent messages
// responses:
// 200: StatisticsResponse
// 500

// swagger:response StatisticsResponse
type StatisticsResponse struct {
	// in:body
	Msgs string `json:"msgs"`
}

// swagger:route GET /mailing/statistics mailing GetMailingStatistics
// Returns statistics about mailings
// responses:
// 200: MailingStatistics
// 500

// swagger:response MailingStatistics
type MailingStatistics struct {
	// in:body
	Msgs string `json:"msgs"`
}

// swagger:route POST /client client AddClient
// Adds new client to DB
// responses:
// 200
// 500

// swagger:parameters AddClient
type NewClient struct {
	// in:body
	Client db.Client `json:"client"`
}

// swagger:route POST /mailing mailing AddMailing
// Adds new mailing to DB
// responses:
// 200
// 500

// swagger:parameters AddMailing
type NewMailing struct {
	// in:body
	Mailing db.Mailing `json:"mailing"`
}

// swagger:route POST /mailing/start mailing StartMailing
// Start mailing
// responses:
// 200
// 500

// swagger:route DELETE /client/{id} client DeleteClient
// Deletes client from DB
// responses:
// 200
// 500

// swagger:parameters DeleteClient
type ClientId struct {
	// in:path
	ClientId string `json:"id"`
}

// swagger:route DELETE /mailing/{id} mailing DeleteMailing
// Deletes mailing from DB
// responses:
// 200
// 500

// swagger:parameters DeleteMailing
type MailingId struct {
	// in:path
	MailingId string `json:"id"`
}
