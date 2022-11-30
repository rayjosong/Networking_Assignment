package main

import (
	"dental-clinic/internal/models"
)

type templateData struct {
	Appointment  *models.Appointment
	Appointments []models.Appointment
	CurrentUser  models.User
}

func (app *application) newTemplateData() *templateData {
	return &templateData{}
}
