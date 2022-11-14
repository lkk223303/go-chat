package repository

import "chatty/models"

type DatabaseRepo interface {
	// databases methods here
	Testing()
	InsertMessage(msg models.EventMessage) error
	InsertMessages(msgs []models.EventMessage) error
	GetMessagesFromUser(userId string) ([]models.EventMessage, error)
}
