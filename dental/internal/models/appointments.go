package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Appointment struct {
	Id        int       `json:"apptID"`
	Patient   User      `json:"patient"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Dentist   string    `json:"dentist"`
	Completed bool      `json:"completed"`
}

type AppointmentsModel []Appointment

// Check if timeslot is still available (def: true if no patient but got dentist)
func (a *Appointment) CheckAvailability() bool {
	return (reflect.DeepEqual(a.Patient, User{}) && len(a.Dentist) != 0)
}

// type AppointmentModel struct {
// 	DB []Username
// }

/*
Methods:
1. View all
2. Insert
*/

// Retrieve all Appointment records
func (a *AppointmentsModel) GetAll() ([]Appointment, error) {
	data, err := os.ReadFile("../../internal/models/appts.json")
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

// TODO: Change string to User
func (a *AppointmentsModel) Insert(patient User, start time.Time, end time.Time, dentist string, completed bool) (string, error) {
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

	// read data from file
	fileLines, err := a.GetAll()
	if err != nil {
		return "", err
	}
	newPayload := append(fileLines, *payload)

	jsonData, err := json.Marshal(newPayload)
	if err != nil {
		return "", err
	}

	err = os.WriteFile("../../internal/models/appts.json", jsonData, 0644)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Json data added: %s", jsonData), nil
}

// TODO: Change string to User
func (a *AppointmentsModel) Update(patient User, start time.Time, end time.Time, dentist string, completed bool) (string, error) {
	payload := &Appointment{
		// TODO: Figure out how to get ID
		Patient:   patient,
		StartTime: start,
		EndTime:   end,
		Dentist:   dentist,
		Completed: completed,
	}

	// read data from file
	fileLines, err := a.GetAll()
	if err != nil {
		return "", err
	}
	newPayload := append(fileLines, *payload)

	jsonData, err := json.Marshal(newPayload)
	if err != nil {
		return "", err
	}

	err = os.WriteFile("../../internal/models/appts.json", jsonData, 0644)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("JSON data Updated: %s", jsonData), nil
}

// Helper func: Retrieve last appointment record
func (a *AppointmentsModel) GetLast() (Appointment, error) {
	s, err := a.GetAll()
	if err != nil {
		return Appointment{}, err
	}

	return s[len(s)-1], err
}

func (a *AppointmentsModel) FormatDateTime(timeStamp time.Time) string {
	year, month, day := timeStamp.Local().Date()

	hour := timeStamp.Hour()
	minute := timeStamp.Minute()

	var newHour string
	var newMin string

	if hour < 10 {
		newHour = fmt.Sprintf("0%d", hour)
	} else {
		newHour = strconv.Itoa(hour)
	}

	if minute < 10 {
		newMin = fmt.Sprintf("0%d", minute)
	} else {
		newMin = strconv.Itoa(minute)
	}

	return fmt.Sprintf("%d-%d-%d %s:%s SGT", year, month, day, newHour, newMin)
}
