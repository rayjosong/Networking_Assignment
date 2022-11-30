package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Appointment struct {
	Id        int    `json:"apptID"`
	Patient   User   `json:"patient"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Dentist   string `json:"dentist"`
	Completed bool   `json:"completed"`
}

type AppointmentsModel []*Appointment

// type AppointmentModel struct {
// 	DB []Username
// }

/*
Methods:
1. View all
2. Insert
*/

func (a *AppointmentsModel) GetAll() ([]Appointment, error) {
	data, err := ioutil.ReadFile("../../internal/models/appts.json")
	if err != nil {
		return []Appointment{}, err
	}
	var appt []Appointment
	err = json.Unmarshal([]byte(data), &appt)
	if err != nil {
		return appt, err
	}

	return appt, nil
}

func (a *AppointmentsModel) GetLast() (Appointment, error) {
	s, err := a.GetAll()
	if err != nil {
		return Appointment{}, err
	}

	return s[len(s)-1], err
}

// TODO: Change string to User
func (a *AppointmentsModel) Insert(patient User, start string, end string, dentist string, completed bool) (string, error) {
	num, err := a.GetLast()
	if err != nil {
		log.Println(err)
	}

	payload := &Appointment{
		Id:        num.Id,
		Patient:   patient,
		StartTime: start,
		EndTime:   end,
		Dentist:   dentist,
		Completed: completed,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Json data added: %s", jsonData), err
}
