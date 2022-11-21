package dbrepo

import "chatty/models"

// postgres method
func (p *postgresRepo) Testing() {

}

func (p *postgresRepo) InsertMessage(msg models.EventMessage) error {

	return nil
}

func (p *postgresRepo) InsertMessages(msgs []models.EventMessage) error {

	return nil
}

func (p *postgresRepo) GetMessagesFromUser(userId string) ([]models.EventMessage, error) {
	var msgEvents []models.EventMessage

	return msgEvents, nil
}
