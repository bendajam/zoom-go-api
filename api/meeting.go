package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
   API Documentation
   https://marketplace.zoom.us/docs/api-reference/zoom-api/meetings/meetingregistrantcreate
*/
func (client Client) AddMeetingRegistrant(meetingID int,
	email,
	firstName,
	lastName,
	address,
	city,
	country,
	zip,
	state,
	phone,
	industry,
	org,
	jobTitle,
	purchasingTimeFrame,
	roleInPurchaseProcess,
	noOfEmployees,
	comments string,
	customQuestions []CustomQuestion) (addMeetingRegistrantResponse AddMeetingRegistrantResponse, err error) {

	addMeetingRegistrantRequest := AddMeetingRegistrantRequest{
		Email:                 email,
		FirstName:             firstName,
		LastName:              lastName,
		Address:               address,
		City:                  city,
		Country:               country,
		Zip:                   zip,
		State:                 state,
		Phone:                 phone,
		Industry:              industry,
		Org:                   org,
		JobTitle:              jobTitle,
		PurchasingTimeFrame:   purchasingTimeFrame,
		RoleInPurchaseProcess: roleInPurchaseProcess,
		NoOfEmployees:         noOfEmployees,
		Comments:              comments,
		CustomQuestions:       customQuestions,
	}

	endpoint := fmt.Sprintf("/meetings/%d/registrants", meetingID)
	httpMethod := http.MethodPost

	reqBody, err := json.Marshal(addMeetingRegistrantRequest)
	if err != nil {
		return
	}
	err = client.doRequestJSON(endpoint, httpMethod, *bytes.NewBuffer(reqBody), &addMeetingRegistrantResponse)

	return
}

func (client Client) CreateMeeting(
	userID,
	topic,
	startTime,
	scheduledFor,
	timezone,
	password,
	agenda string,
	meetingType,
	duration int,
	recurrence *Recurrence,
	settings *Settings) (createMeetingResponse CreateMeetingResponse, err error) {

	if recurrence == nil {
		recurrence = &Recurrence{
			Type:           1,
			RepeatInterval: 0,
			WeeklyDays:     "",
			MonthlyDay:     0,
			MonthlyWeek:    0,
			MonthlyWeekDay: 0,
			EndTimes:       0,
			EndDateTime:    "",
		}
	}

	if settings == nil {
		settings = &Settings{
			HostVideo:                    false,
			ParticipantVideo:             false,
			CnMeeting:                    false,
			InMeeting:                    false,
			JoinBeforeHost:               false,
			MuteUponEntry:                false,
			Watermark:                    false,
			UsePmi:                       false,
			ApprovalType:                 0,
			RegistrationType:             0,
			Audio:                        "",
			AutoRecording:                "",
			EnforceLogin:                 false,
			EnforceLoginDomains:          "",
			AlternativeHosts:             "",
			GlobalDialInCountries:        nil,
			RegistrantsEmailNotification: false,
		}
	}

	createMeetingRequest := CreateMeetingRequest{
		Topic:       topic,
		Type:        meetingType,
		StartTime:   startTime,
		Duration:    duration,
		ScheduleFor: scheduledFor,
		Timezone:    timezone,
		Password:    password,
		Agenda:      agenda,
		Recurrence:  *recurrence,
		Settings:    *settings,
	}

	var reqBody []byte
	reqBody, err = json.Marshal(createMeetingRequest)
	if err != nil {
		return
	}

	endpoint := fmt.Sprintf("/users/%s/meetings", userID)
	httpMethod := http.MethodPost

	err = client.doRequestJSON(endpoint, httpMethod, *bytes.NewBuffer(reqBody), &createMeetingResponse)

	return
}

func (client Client) DeleteMeeting(meetingID int) (err error) {

	endpoint := fmt.Sprintf("/meetings/%d", meetingID)
	httpMethod := http.MethodDelete

	var b []byte
	_, err = client.doRequest(endpoint, httpMethod, *bytes.NewBuffer(b))
	return
}

func (client Client) GetMeeting(meetingID int) (getMeetingResponse GetMeetingResponse, err error) {

	endpoint := fmt.Sprintf("/meetings/%d", meetingID)
	httpMethod := http.MethodGet

	var b []byte
	err = client.doRequestJSON(endpoint, httpMethod, *bytes.NewBuffer(b), &getMeetingResponse)
	return
}

func (client Client) GetMeetingInvitation(meetingID int) (getMeetingInvitationResponse GetMeetingInvitationResponse, err error) {

	endpoint := fmt.Sprintf("/meetings/%d/invitation", meetingID)
	httpMethod := http.MethodGet

	var b []byte
	err = client.doRequestJSON(endpoint, httpMethod, *bytes.NewBuffer(b), &getMeetingInvitationResponse)
	return
}

func (client Client) ListMeetingRegistrants(meetingID int) (apiResponse ListMeetingRegistrantsResponse, err error) {

	endpoint := fmt.Sprintf("/meetings/%d/registrants", meetingID)
	httpMethod := http.MethodGet

	var b []byte
	err = client.doRequestJSON(endpoint, httpMethod, *bytes.NewBuffer(b), &apiResponse)
	return
}

func (client Client) ListMeetings(userID string) (apiResponse ListMeetingsAPIResponse, err error) {

	endpoint := fmt.Sprintf("/users/%s/meetings", userID)
	httpMethod := http.MethodGet

	var b []byte
	err = client.doRequestJSON(endpoint, httpMethod, *bytes.NewBuffer(b), &apiResponse)
	return
}

func (client Client) UpdateMeetingStatus(meetingID int, status string) (err error) {

	updateMeetingStatusRequest := UpdateMeetingStatusRequest{
		Action: status,
	}

	var reqBody []byte
	reqBody, err = json.Marshal(updateMeetingStatusRequest)
	if err != nil {
		return
	}

	endpoint := fmt.Sprintf("/meetings/%d/status", meetingID)
	httpMethod := http.MethodPut

	_, err = client.doRequest(endpoint, httpMethod, *bytes.NewBuffer(reqBody))
	return
}