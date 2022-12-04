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
	return (reflect.DeepEqual(a.Patient, User{}) && len(a.Dentist) != 0 && !a.Completed)
}

func (a *AppointmentsModel) Get(user User) ([]Appointment, error) {
	appointments, err := a.GetAll()
	if err != nil {
		return []Appointment{}, err
	}

	// search for records where user exists
	var records []Appointment
	for _, record := range appointments {
		if record.Patient.Username == user.Username {
			records = append(records, record)
		}
	}

	return records, nil
}

// func (a *AppointmentsModel) GetApptById(id int) (Appointment, error) {
// 	appointments, err := a.GetAll()
// 	if err != nil {
// 		return Appointment{}, err
// 	}

// 	// search for records where apptId = id
// 	var appt Appointment
// 	for _, record := range appointments {
// 		if record.Id == id {
// 			appt = record
// 		}
// 	}

// 	return appt, nil
// }

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

// Insert one record into database
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

// Edit one record on the json database
func (a *AppointmentsModel) Update(apptId int, patient User, completed bool) (*Appointment, error) {

	// read data from file
	fileLines, err := a.GetAll()
	if err != nil {
		return &Appointment{}, err
	}

	var newFileLines []Appointment
	var updatedPayload Appointment
	// find appointment
	for _, appt := range fileLines {
		if appt.Id == apptId {
			appt.Patient = patient
			updatedPayload = appt
		}
		newFileLines = append(newFileLines, appt)
	}

	jsonData, err := json.Marshal(newFileLines)
	if err != nil {
		return &Appointment{}, err
	}

	err = os.WriteFile("../../internal/models/appts.json", jsonData, 0644)
	if err != nil {
		return &Appointment{}, err
	}

	return &updatedPayload, nil
}

// Helper func: Retrieve last appointment record
func (a *AppointmentsModel) GetLast() (Appointment, error) {
	s, err := a.GetAll()
	if err != nil {
		return Appointment{}, err
	}

	return s[len(s)-1], err
}

// Retrieve all appointments that are available for booking
func (a *AppointmentsModel) GetAvailable() ([]Appointment, error) {
	all, err := a.GetAll()
	if err != nil {
		return []Appointment{}, err
	}

	var availAppts []Appointment

	for _, appt := range all {
		if appt.CheckAvailability() {
			availAppts = append(availAppts, appt)
		}
	}
	return availAppts, nil
}

// Given the appointmentID, delete the recod from the json database
func (a *AppointmentsModel) Delete(apptID int) error {

	// 1. read data from file
	data, err := os.ReadFile("../../internal/models/appts.json")
	if err != nil {
		return err
	}

	var sliceAppts []Appointment
	err = json.Unmarshal([]byte(data), &sliceAppts)
	if err != nil {
		return err
	}

	// Find the record
	indexToDel, err := FindIndexFromSlice(apptID, sliceAppts)
	if err != nil {
		return err
	}

	fmt.Println("Index: ", indexToDel)

	sliceAppts = func(s []Appointment, index int) []Appointment {
		// Delete the record &&
		back := s[index+1:]
		newBack := []Appointment{}
		for _, record := range back {
			// Auto decrement the records below
			record.Id = record.Id - 1
			newBack = append(newBack, record)
		}
		return append(s[:index], newBack...)
	}(sliceAppts, indexToDel)

	jsonData, err := json.Marshal(sliceAppts)
	if err != nil {
		return err
	}

	err = os.WriteFile("../../internal/models/appts.json", jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func FindIndexFromSlice(apptID int, a []Appointment) (int, error) {
	for index, appt := range a {
		if appt.Id == apptID {
			return index, nil
		}
	}

	return 0, fmt.Errorf("the record of id = %d is not found", apptID)
}

// type AppointmentModel struct {
// 	DB []Username
// }

/*
	Methods:
	1. View all
	2. Insert
*/

// HELPER

// Format date-time to the "YY-MM-DD HH-MM" format
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
