package storage

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Session struct {
	CustomerCode string    `db:"customer_code"`
	Customer     Customer  `db:"session_data"`
	Created      time.Time `db:"created"`
	Updated      time.Time `db:"updated"`
}

type Customer struct {
	FrgnRefNo    string `json:"foreign_reference_no"`
	CustomerCode string `json:"customer_code"`

	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`

	Phone  string `json:"phone"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`

	Birthdate    string `json:"birthdate"`
	BirthCountry string `json:"country_of_birth"`
	Nationality  string `json:"nationality"`
	Gender       string `json:"gender"`

	Occupation   string `json:"occupation"`
	EmployerName string `json:"nameofemployer"`

	Wucardno      string `json:"wucardno"`
	Debitcardno   string `json:"debitcardno"`
	Loyaltycardno string `json:"loyaltycardno"`

	IsTemporaryPassword string `json:"isTemporaryPassword"`
	PasswordExpiration  string `json:"password_expiration"`

	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`

	CardPoints        json.Number `json:"card_points"`
	OneSignalPlayerID string      `json:"one_signal_player_id"`
	Tin               string      `json:"tin"`
	SSS               string      `json:"sss"`

	Secretquestion1 string `json:"-"`
	Answer1         string `json:"-"`
	Secretquestion2 string `json:"-"`
	Answer2         string `json:"-"`
	Secretquestion3 string `json:"-"`
	Answer3         string `json:"-"`

	PresentAddress    string `json:"presentaddress"`
	PresentCity       string `json:"presentcity"`
	PresentState      string `json:"presentstate"`
	PresentProvince   string `json:"presentprovince"`
	PresentRegion     string `json:"presentregion"`
	PresentCountry    string `json:"presentcountry"`
	PresentPostalCode string `json:"presentpostalcode"`

	PermanentAddress    string `json:"permanentaddress"`
	PermanentCity       string `json:"permanentcity"`
	PermanentState      string `json:"permanentstate"`
	PermanentProvince   string `json:"permanentprovince"`
	PermanentRegion     string `json:"permanentregion"`
	PermanentCountry    string `json:"permanentcountry"`
	PermanentPostalCode string `json:"permanentpostalcode"`

	ProfileImage string `json:"ProfileImage"`
	IDImage      string `json:"IdImage"`

	IDType         string `json:"id_type"`
	CustomerIDNo   string `json:"customer_id_no"`
	CountryIDIssue string `json:"country_id_issue"`
	IDIssueDate    string `json:"id_issue_date"`
	IDExpiration   string `json:"id_expiration"`
	Acountry       string `json:"acountry"`

	Idrows []IDEntry `json:"idrows"`
}

type IDEntry struct {
	IDType           string `json:"id_type"`
	IDPhoto          string `json:"id_photo"`
	IDExpirationDate string `json:"id_expiration_date"`
	CustomerIDNumber string `json:"customer_id_number"`
	IDCountryIssue   string `json:"id_country_issue"`
	IDIssueDate      string `json:"id_issue_date"`
	DateAdded        string `json:"date_added"`
}

func (c Customer) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Customer) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}
