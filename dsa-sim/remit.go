package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	mrand "math/rand"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/serviceutil/logging"
	"github.com/gorilla/csrf"
	"google.golang.org/grpc/status"

	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	ct "brank.as/petnet/gunk/drp/v1/quote"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
)

const (
	remTypeNotImplemented = "not-implemented"
	remTypeDisburse       = "disburse"
	remTypeCreate         = "create"
	remitSess             = "remitsess"
	remitRegisterUser     = "register-user"
	remitCreateProfile    = "create-profile"
	remitCreateRecipient  = "create-recipient"
	remitCreateQuote      = "create-quote"
)

type remitForm struct {
	CSRFField          template.HTML
	Partners           map[string]*pnpb.RemitPartner
	Partner            string
	CreateRemitReq     *tpb.CreateRemitRequest
	DisburseRemitReq   *tpb.DisburseRemitRequest
	RegisterUserReq    *ppb.RegisterUserRequest
	CreateProfileReq   *ppb.CreateProfileRequest
	CreateRecipientReq *ppb.CreateRecipientRequest
	CreateQuoteReq     *ct.CreateQuoteRequest
	RequirementsName   string
	RequirementsValue  string
	Gender             string
	KYCVerified        string
	ProofOfAddress     string
	Type               string
	NotImplemented     bool
	Types              []string
	Error              string
	Success            bool
	ControlNo          string
}

type listRemitForm struct {
	CSRFField   template.HTML
	Remittances []*tpb.Remittance
}

func (s *Server) getRemit(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	template := s.templates.Lookup("remit.html")
	if template == nil {
		log.Error("unable to load template")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ss, err := s.sess.Get(r, remitSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	f := remitForm{}
	var fcache bool
	if fc, ok := ss.Values["form-cache"]; ok && fc != "" {
		err := json.Unmarshal([]byte(fc.(string)), &f)
		ss.Values["form-cache"] = ""
		if err := ss.Save(r, w); err != nil {
			logging.WithError(err, log).Error("saving session")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			logging.WithError(err, log).Error("unmarshaling remit cache")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fcache = true
	}

	res, err := s.cl.apiClient.RemitPartners(ctx, &pnpb.RemitPartnersRequest{
		Country: "PH",
	})
	if err != nil {
		logging.WithError(err, log).Error("listing remit partners")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f.Partners = res.GetPartners()
	f.CSRFField = csrf.TemplateField(r)

	if !fcache {
		var ptnr string
		var rt string
		if v, ok := ss.Values["partner"]; ok {
			ptnr = v.(string)
		}
		if v, ok := ss.Values["type"]; ok {
			rt = v.(string)
		}
		// set default to WU and create if no partner and type has been selected
		if ptnr == "" {
			ptnr = "WU"
			rt = remTypeCreate
		}

		f = remitForm{
			Partners:  res.GetPartners(),
			Partner:   ptnr,
			CSRFField: csrf.TemplateField(r),
			Type:      rt,
		}

		switch ptnr {
		case "WU":
			f.Types = []string{remTypeCreate, remTypeDisburse}
			if rt == remTypeCreate {
				f.CreateRemitReq = wuCreateReq
			} else {
				f.DisburseRemitReq = wuDisburseReq
			}
		case "IR":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = irDisburseReq
		case "TF":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = tfDisburseReq
		case "RM":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = rmDisburseReq
		case "RIA":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = riaDisburseReq
		case "MB":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = mbDisburseReq
		case "BPI":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = bpiDisburseReq
		case "JPR":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = jprDisburseReq
		case "USSC":
			f.Types = []string{remTypeCreate, remTypeDisburse}
			if rt == remTypeCreate {
				f.CreateRemitReq = usscCreateReq
			} else {
				f.DisburseRemitReq = usscDisburseReq
			}
		case "IC":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = icDisburseReq
		case "UNT":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = untDisburseReq
		case "CEBINT":
			f.Type = remTypeDisburse
			f.Types = []string{remTypeDisburse}
			f.DisburseRemitReq = cebintDisburseReq
		case "CEB":
			f.Types = []string{remTypeCreate, remTypeDisburse}
			if rt == remTypeCreate {
				f.CreateRemitReq = cebCreateReq
			} else {
				f.DisburseRemitReq = cebDisburseReq
			}
		case "AYA":
			f.Types = []string{remTypeCreate, remTypeDisburse}
			if rt == remTypeCreate {
				f.CreateRemitReq = ayaCreateReq
			} else {
				f.DisburseRemitReq = ayaDisburseReq
			}
		case "WISE":
			f.Types = []string{remitRegisterUser, remitCreateProfile, remitCreateRecipient, remitCreateQuote, remTypeCreate}
			if rt == remTypeCreate {
				f.CreateRemitReq = wiseCreateReq
			} else if rt == remitRegisterUser {
				f.RegisterUserReq = wiseRegisterUser
			} else if rt == remitCreateProfile {
				f.CreateProfileReq = wiseCreateProfileReq
			} else if rt == remitCreateRecipient {
				f.CreateRecipientReq = wiseCreateRecipientReq
				f.RequirementsName = "legalType"
				f.RequirementsValue = "PRIVATE"
			} else if rt == remitCreateQuote {
				f.CreateQuoteReq = wiseCreateQuoteReq
			}
		case "IE":
			f.Types = []string{remTypeCreate, remTypeDisburse}
			if rt == remTypeCreate {
				f.CreateRemitReq = ieCreateReq
			} else {
				f.DisburseRemitReq = ieDisburseReq
			}
		default:
			f.Type = remTypeNotImplemented
			f.Types = []string{remTypeNotImplemented}
		}
	}

	if err := template.Execute(w, f); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) postCreateQuote(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	var f remitForm
	err := r.ParseForm()
	if err != nil {
		log.Error()
		return
	}

	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding post form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.Types = strings.Fields(strings.Trim(f.Types[0], "[]"))
	f.CreateQuoteReq.RemitPartner = f.Partner
	_, err = s.cl.apiClient.CreateQuote(ctx, f.CreateQuoteReq)
	if err != nil {
		logging.WithError(err, log).Error("creating Recipient")
		f.Error = getError(err)
	}
	if err == nil {
		f.Success = true
	}

	ss, err := s.sess.Get(r, remitSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(f)
	if err != nil {
		logging.WithError(err, log).Error("marshaling remit cache")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ss.Values["form-cache"] = string(b)
	if err := ss.Save(r, w); err != nil {
		logging.WithError(err, log).Error("saving session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, remitPath, http.StatusSeeOther)
}

func (s *Server) postChoosePartner(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	var f remitForm
	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding post form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ss, err := s.sess.Get(r, remitSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ss.Values["partner"] = f.Partner
	ss.Values["type"] = f.Type
	if err := ss.Save(r, w); err != nil {
		logging.WithError(err, log).Error("saving session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, remitPath, http.StatusSeeOther)
}

func (s *Server) postCreateRecipient(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	var f remitForm
	err := r.ParseForm()
	if err != nil {
		log.Error()
		return
	}

	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding post form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.Types = strings.Fields(strings.Trim(f.Types[0], "[]"))
	f.CreateRecipientReq.RemitPartner = f.Partner
	var rqr []*ppb.Requirement
	rqr = append(rqr, &ppb.Requirement{
		Name:  f.RequirementsName,
		Value: f.RequirementsValue,
	})
	f.CreateRecipientReq.Requirements = rqr
	_, err = s.cl.apiClient.CreateRecipient(ctx, f.CreateRecipientReq)
	if err != nil {
		logging.WithError(err, log).Error("creating Recipient")
		f.Error = getError(err)
	}
	if err == nil {
		f.Success = true
	}

	ss, err := s.sess.Get(r, remitSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(f)
	if err != nil {
		logging.WithError(err, log).Error("marshaling remit cache")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ss.Values["form-cache"] = string(b)
	if err := ss.Save(r, w); err != nil {
		logging.WithError(err, log).Error("saving session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, remitPath, http.StatusSeeOther)
}

func (s *Server) postCreateProfile(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	var f remitForm
	err := r.ParseForm()
	if err != nil {
		log.Error()
		return
	}

	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding post form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.Types = strings.Fields(strings.Trim(f.Types[0], "[]"))
	f.CreateProfileReq.RemitPartner = f.Partner

	_, err = s.cl.apiClient.CreateProfile(ctx, f.CreateProfileReq)
	if err != nil {
		logging.WithError(err, log).Error("creating user")
		f.Error = getError(err)
	}
	if err == nil {
		f.Success = true
	}

	ss, err := s.sess.Get(r, remitSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(f)
	if err != nil {
		logging.WithError(err, log).Error("marshaling remit cache")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ss.Values["form-cache"] = string(b)
	if err := ss.Save(r, w); err != nil {
		logging.WithError(err, log).Error("saving session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, remitPath, http.StatusSeeOther)
}

func (s *Server) postRegisterUser(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	var f remitForm
	err := r.ParseForm()
	if err != nil {
		log.Error()
		return
	}

	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding post form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.Types = strings.Fields(strings.Trim(f.Types[0], "[]"))
	f.RegisterUserReq.RemitPartner = f.Partner
	_, err = s.cl.apiClient.RegisterUser(ctx, f.RegisterUserReq)
	if err != nil {
		logging.WithError(err, log).Error("registering user")
		f.Error = getError(err)
	}
	if err == nil {
		f.Success = true
	}

	ss, err := s.sess.Get(r, remitSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(f)
	if err != nil {
		logging.WithError(err, log).Error("marshaling remit cache")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ss.Values["form-cache"] = string(b)
	if err := ss.Save(r, w); err != nil {
		logging.WithError(err, log).Error("saving session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, remitPath, http.StatusSeeOther)
}

func (s *Server) postCreateRemit(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	var f remitForm

	err := r.ParseForm()
	if err != nil {
		log.Error()
		return
	}

	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding post form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// adding some values manually, as the decoder doesn't decode those values
	remState := r.PostForm.Get("CreateRemitReq.Remitter.ContactInfo.Address.State")
	recState := r.PostForm.Get("CreateRemitReq.Receiver.ContactInfo.Address.State")
	f.CreateRemitReq.RemitPartner = f.Partner
	f.CreateRemitReq.Remitter.Gender = tpb.Gender(tpb.Gender_value[f.Gender])
	f.CreateRemitReq.Remitter.KYCVerified = tpb.Bool(tpb.Bool_value[f.KYCVerified])
	f.CreateRemitReq.Remitter.ContactInfo.Address.State = remState
	f.CreateRemitReq.Receiver.ContactInfo.Address.State = recState
	// trimming and splitting because its saved as json array
	f.Types = strings.Fields(strings.Trim(f.Types[0], "[]"))

	// values set by the DSA system
	f.CreateRemitReq.OrderID = numberString(20)
	f.CreateRemitReq.Agent = &tpb.Agent{
		UserID:    4456,
		IPAddress: "127.0.0.1",
		DeviceID:  "dc29d6674a776db145af78f5ac20a293409a6c1f807885bbb5",
	}

	res, err := s.cl.apiClient.CreateRemit(ctx, f.CreateRemitReq)
	if err != nil {
		logging.WithError(err, log).Error("creating remit")
		f.Error = getError(err)
	}
	if f.Error == "" {
		res, err := s.cl.apiClient.ConfirmRemit(ctx, &tpb.ConfirmRemitRequest{
			TransactionID: res.TransactionID,
			AuthSource:    "User Review",
			AuthCode:      "Manual",
		})
		if err != nil {
			logging.WithError(err, log).Error("creating remit")
			f.Error = getError(err)
		} else {
			f.Success = true
			f.ControlNo = res.GetControlNumber()
		}
	}

	ss, err := s.sess.Get(r, remitSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(f)
	if err != nil {
		logging.WithError(err, log).Error("marshaling remit cache")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ss.Values["form-cache"] = string(b)
	if err := ss.Save(r, w); err != nil {
		logging.WithError(err, log).Error("saving session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, remitPath, http.StatusSeeOther)
}

func (s *Server) postDisburseRemit(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	var f remitForm

	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding post form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// adding some values manually, as the decoder doesn't decode those values
	f.DisburseRemitReq.RemitPartner = f.Partner
	recProv := r.PostForm.Get("DisburseRemitReq.Receiver.ContactInfo.Address.Province")
	recState := r.PostForm.Get("DisburseRemitReq.Receiver.ContactInfo.Address.State")
	f.DisburseRemitReq.Receiver.ContactInfo.Address.Province = recProv
	f.DisburseRemitReq.Receiver.ContactInfo.Address.State = recState
	f.DisburseRemitReq.Receiver.Gender = tpb.Gender(tpb.Gender_value[f.Gender])
	f.DisburseRemitReq.Receiver.KYCVerified = tpb.Bool(tpb.Bool_value[f.KYCVerified])
	f.DisburseRemitReq.Receiver.ProofOfAddress = tpb.Bool(tpb.Bool_value[f.ProofOfAddress])
	f.Types = strings.Fields(strings.Trim(f.Types[0], "[]"))

	// set by DSA system
	f.DisburseRemitReq.RemitType = "Payout"
	f.DisburseRemitReq.OrderID = numberString(20)
	f.DisburseRemitReq.Agent = &tpb.Agent{
		UserID:    4456,
		IPAddress: "127.0.0.1",
		DeviceID:  "dc29d6674a776db145af78f5ac20a293409a6c1f807885bbb5",
	}
	f.DisburseRemitReq.Transaction = &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	}

	res, err := s.cl.apiClient.DisburseRemit(ctx, f.DisburseRemitReq)
	if err != nil {
		logging.WithError(err, log).Error("disbursing remit")
		f.Error = getError(err)
	}
	if f.Error == "" {
		res, err := s.cl.apiClient.ConfirmRemit(ctx, &tpb.ConfirmRemitRequest{
			TransactionID: res.TransactionID,
			AuthSource:    "User Review",
			AuthCode:      "Manual",
		})
		if err != nil {
			logging.WithError(err, log).Error("creating remit")
			f.Error = getError(err)
		} else {
			f.Success = true
			f.ControlNo = res.GetControlNumber()
		}
	}

	ss, err := s.sess.Get(r, remitSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := json.Marshal(f)
	if err != nil {
		logging.WithError(err, log).Error("marshaling remit cache")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ss.Values["form-cache"] = string(b)
	if err := ss.Save(r, w); err != nil {
		logging.WithError(err, log).Error("saving session")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, remitPath, http.StatusSeeOther)
}

func (s *Server) getListRemit(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	template := s.templates.Lookup("list-remit.html")
	if template == nil {
		log.Error("unable to load template")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	res, err := s.cl.apiClient.ListRemit(ctx, &tpb.ListRemitRequest{
		Limit: 300,
	})
	if err != nil {
		logging.WithError(err, log).Debug("no transactions found")
	}

	f := listRemitForm{
		CSRFField: csrf.TemplateField(r),
	}

	f.Remittances = res.GetRemittances()
	if err != nil {
		f.Remittances = nil
	}

	if err := template.Execute(w, f); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func numberString(n int) string {
	letters := []rune("123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

var wiseCreateQuoteReq = &ct.CreateQuoteRequest{
	RemitPartner: static.WISECode,
	Email:        "test@test.com",
	Amount: &ct.QuoteAmount{
		SourceAmount:        "100",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
	},
}

var wiseCreateRecipientReq = &ppb.CreateRecipientRequest{
	RemitPartner:      static.WISECode,
	Email:             "sender2@brankas.com",
	Currency:          "GBP",
	Type:              "personal",
	OwnedByCustomer:   false,
	AccountHolderName: "jhon doe",
	Requirements: []*ppb.Requirement{
		{
			Name:  "legalType",
			Value: "PRIVATE",
		},
	},
}

var wiseCreateProfileReq = &ppb.CreateProfileRequest{
	RemitPartner: static.WISECode,
	Email:        "sender2@brankas.com",
	Type:         "personal",
	FirstName:    "Brankas",
	LastName:     "Sender",
	BirthDate:    "1990-01-01",
	Phone: &ppb.PhoneNumber{
		CountryCode: "63",
		Number:      "9999999999",
	},
	Address: &ppb.Address{
		Address1:   "East Offices Bldg., 114 Aguirre St.,Legaspi Village,",
		Address2:   "East Offices Bldg",
		City:       "Makati",
		State:      "Makati",
		Province:   "Makati",
		Zone:       "Makati",
		PostalCode: "1229",
		Country:    "PH",
	},
	Occupation: "Software Engineer",
}

var wiseRegisterUser = &ppb.RegisterUserRequest{
	RemitPartner: static.WISECode,
	Email:        "test@test.com",
}

var wuCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: "WU",
	RemitType:    "Send",
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "82273864671",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "82273864672",
			},
			Address: &tpb.Address{
				Address1:   "CALINAN",
				Address2:   "CALINAN",
				City:       "DAVAO CITY",
				State:      "DAVAO DEL SUR",
				PostalCode: "8000A",
				Country:    "PH",
			},
		},
		Identification: &ppb.Identification{
			Type:    "M",
			Number:  "24023497AB0877AAB20000",
			Country: "PH",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Birthdate: &ppb.Date{
			Year:  "2000",
			Month: "12",
			Day:   "12",
		},
		BirthCountry:       "PH",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
		SourceFunds:        "Savings",
		TransactionPurpose: "Gift",
		Email:              "johndoe@mail.com",
		PartnerMemberID:    "12351351361",
		ReceiverRelation:   "Family",
	},
	Amount: &tpb.SendAmount{
		Amount:              "1000",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationAmount:   true,
		DestinationCountry:  "PH",
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "Jane",
			MiddleName: "Emily",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "82273864673",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "82273864674",
			},
			Address: &tpb.Address{
				Address1:   "CALINAN",
				Address2:   "CALINAN",
				City:       "DAVAO CITY",
				State:      "DAVAO DEL SUR",
				PostalCode: "8000A",
				Country:    "PH",
			},
		},
	},
}

var usscCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.USSCCode,
	RemitType:    "Send",
	OrderID:      "5142096205",
	Amount: &tpb.SendAmount{
		Amount:              "1000",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Email:      "emai@email.com",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		KYCVerified:        tpb.Bool_True,
		Gender:             tpb.Gender_Male,
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &ppb.Identification{
			Type:    "M",
			Number:  "24023497AB0877AAB20000",
			Country: "PH",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2031",
				Month: "12",
				Day:   "12",
			},
		},
	},
}

var wuDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     "WU",
	ControlNumber:    "",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "82273864673",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "82273864674",
			},
			Address: &tpb.Address{
				Address1:   "CALINAN",
				Address2:   "CALINAN",
				City:       "DAVAO CITY",
				State:      "DAVAO DEL SUR",
				PostalCode: "8000A",
				Country:    "PH",
			},
		},
		Identification: &ppb.Identification{
			Type:    "M",
			Number:  "24023497AB0877AAB20000",
			Country: "PH",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Birthdate: &ppb.Date{
			Year:  "2000",
			Month: "12",
			Day:   "12",
		},
		BirthCountry:       "PH",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
		SourceFunds:        "Savings",
		TransactionPurpose: "Gift",
		Email:              "johndoe@mail.com",
		PartnerMemberID:    "12351351361",
		ReceiverRelation:   "Family",
	},
}

var irDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  "IR",
	ControlNumber: "",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "12345",
			},
			Address: &tpb.Address{
				Address1:   "CALINAN",
				Address2:   "SOMETHING",
				City:       "MALOLOS",
				Province:   "BULACAN",
				PostalCode: "8000A",
				Country:    "PH",
				Zone:       "BULIHAN",
			},
		},
		Identification: &ppb.Identification{
			Type:   "Voter's ID",
			Number: "76050337ae0666iub200010",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Birthdate: &ppb.Date{
			Year:  "2000",
			Month: "12",
			Day:   "12",
		},
		BirthCountry:       "PH",
		BirthPlace:         "PASAY CITY,PASAY CITY",
		SourceFunds:        "Savings",
		TransactionPurpose: "Gift",
		PartnerMemberID:    "12351351361",
		ReceiverRelation:   "Family",
	},
}

var tfDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  "TF",
	ControlNumber: "",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "12345",
			},
			Address: &tpb.Address{
				Address1:   "CALINAN",
				Address2:   "SOMETHING",
				City:       "MALOLOS",
				Province:   "BULACAN",
				PostalCode: "8000A",
				Country:    "PH",
				Zone:       "BULIHAN",
			},
		},
		Identification: &ppb.Identification{
			Type:   "Voter's ID",
			Number: "76050337ae0666iub200010",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2040",
				Month: "12",
				Day:   "12",
			},
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation:   "Student",
			OccupationID: "1",
		},
		Birthdate: &ppb.Date{
			Year:  "2000",
			Month: "12",
			Day:   "12",
		},
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Savings",
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "6374360",
		ReceiverRelation:   "Friend",
		SendingReasonID:    "1",
		KYCVerified:        tpb.Bool_True,
		ProofOfAddress:     tpb.Bool_True,
		Gender:             tpb.Gender_Male,
	},
}

var rmDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     "RM",
	ControlNumber:    "",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "12345",
			},
			Address: &tpb.Address{
				Address1:   "CALINAN",
				Address2:   "SOMETHING",
				City:       "MALOLOS",
				Province:   "BULACAN",
				State:      "PH-00",
				PostalCode: "8000A",
				Country:    "PH",
				Zone:       "BULIHAN",
			},
		},
		Identification: &ppb.Identification{
			Type:   "GOVERNMENT_ISSUED_ID",
			Number: "B83180608851",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2040",
				Month: "12",
				Day:   "12",
			},
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation:   "OTH",
			OccupationID: "1",
		},
		Birthdate: &ppb.Date{
			Year:  "2000",
			Month: "12",
			Day:   "12",
		},
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Savings",
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "6374360",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
	},
}

var riaDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  "RIA",
	ControlNumber: "",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &ppb.Identification{
			Type:   "GOVERNMENT_ISSUED_ID",
			Number: "B83180608851",
			Expiration: &ppb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &ppb.Date{
				Year:  "2000",
				Month: "03",
				Day:   "01",
			},
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
	},
}

var mbDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  "MB",
	ControlNumber: "",
	OrderID:       "5142096205",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09265454935",
			},
			Address: &tpb.Address{
				Address1:   "18 SITIO PULO",
				Address2:   "",
				City:       "MALOLOS",
				Province:   "BULACAN",
				PostalCode: "3000A",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &ppb.Identification{
			Type:   "Postal ID",
			Number: "PRND32200265569P",
			Expiration: &ppb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Birthdate: &ppb.Date{
			Year:  "2000",
			Month: "12",
			Day:   "15",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
	},
	Remitter: &tpb.Contact{
		FirstName:  "Mercado",
		MiddleName: "Marites",
		LastName:   "Cueto",
	},
}

var bpiDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  "BPI",
	ControlNumber: "",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &ppb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var jprDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     "JPR",
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &ppb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
}

var usscDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	ControlNumber:    "",
	RemitPartner:     "USSC",
	OrderID:          "7885645988",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "MAIN ST",
				Address2:   "addr2",
				City:       "AKLAN CITY",
				State:      "PH-00",
				PostalCode: "36989",
				Country:    "PH",
				Zone:       "TALLAOEN",
				Province:   "ILOCOS SUR",
			},
		},
		Identification: &ppb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
			Issued: &ppb.Date{
				Year:  "2000",
				Month: "01",
				Day:   "03",
			},
			Expiration: &ppb.Date{
				Year:  "2000",
				Month: "01",
				Day:   "03",
			},
		},
		Employment: &tpb.Employment{
			OccupationID: "1",
			Occupation:   "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1990",
			Month: "01",
			Day:   "01",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "6925594",
		ReceiverRelation:   "Family",
		BirthCountry:       "PH",
		BirthPlace:         "TAGUDIN,ILOCOS SUR",
		SourceFunds:        "Savings",
		SendingReasonID:    "1",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
		DeviceID:  "dc29d6674a776db145af78f5ac20a293409a6c1f807885bbb5",
	},
}

var icDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:     "Payout",
	RemitPartner:  "IC",
	ControlNumber: "5142096205",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &ppb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
}

var untDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     "UNT",
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				State:      "state",
				Zone:       "TALLAOEN",
			},
		},
		Identification: &ppb.Identification{
			Type:    "LICENSE",
			Number:  "B83180608851",
			Country: "PH",
			Expiration: &ppb.Date{
				Year:  "2050",
				Month: "12",
				Day:   "12",
			},
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		Gender:             tpb.Gender_Male,
		Nationality:        "PH",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
		DeviceID:  "5500",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var cebCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.CEBCode,
	RemitType:    "Send",
	OrderID:      "5142096205",
	Amount: &tpb.SendAmount{
		Amount:              "1000",
		SourceCountry:       "PH",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Email:      "emai@email.com",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "BULIHAN",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		RecipientID: "12345",
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "BULIHAN",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &ppb.Identification{
			Type:    "M",
			Number:  "24023497AB0877AAB20000",
			Country: "PH",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2031",
				Month: "12",
				Day:   "12",
			},
		},
	},
}

var cebDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.CEBCode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "BULIHAN",
			},
		},
		Identification: &ppb.Identification{
			Type:    "GOVERNMENT_ISSUED_ID",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var wiseCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.WISECode,
	RemitType:    "Send",
	OrderID:      "5142096205",
	Message:      "send money",
	Amount: &tpb.SendAmount{
		Amount:              "1000",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Receiver: &tpb.Receiver{
		AccountHolderName:   "John Michael Doe",
		RecipientID:         "1",
		SourceAccountNumber: "1234",
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Email:      "emai@email.com",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "BULIHAN",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
	},
	Remitter: &tpb.UserKYC{
		Email: "john@test.com",
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "BULIHAN",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12345678",
			},
		},
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &ppb.Identification{
			Type:    "M",
			Number:  "24023497AB0877AAB20000",
			Country: "PH",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2031",
				Month: "12",
				Day:   "12",
			},
		},
	},
}

var cebintDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.CEBINTCode,
	ControlNumber:    "5142096205",
	OrderID:          "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				State:      "state",
				Zone:       "zone",
			},
		},
		Identification: &ppb.Identification{
			Type:    "LICENSE",
			Number:  "B83180608851",
			Country: "PH",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
		DeviceID:  "5500",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var ieDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.IECode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &ppb.Identification{
			Type:    "PASSPORT",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var ieCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.IECode,
	RemitType:    "Send",
	OrderID:      "5142096965",
	Amount: &tpb.SendAmount{
		Amount:              "1000",
		SourceCountry:       "PH",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Receiver: &tpb.Receiver{
		Identification: &ppb.Identification{
			Number: "24023497AB0877AAB895668",
		},
		ContactInfo: &tpb.Contact{
			FirstName:  "lkhds",
			MiddleName: "asdksad",
			LastName:   "sdjgsajdh",
			Email:      "emai@email.com",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "123455673",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "123458983",
			},
		},
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "asdkhkh",
			MiddleName: "alsdkdfb",
			LastName:   "skdhfid",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12340967",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12340891",
			},
		},
		PartnerMemberID:    "77125891",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &ppb.Identification{
			Type:    "PASSPORT",
			Country: "PH",
			Number:  "24023497AB0877AAB895668",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2031",
				Month: "12",
				Day:   "12",
			},
		},
	},
	Agent: &tpb.Agent{
		UserID:    1797,
		IPAddress: "130.211.2.002",
	},
}

var ayaDisburseReq = &tpb.DisburseRemitRequest{
	RemitType:        "Payout",
	RemitPartner:     static.AYACode,
	ControlNumber:    "5142096205",
	DisburseCurrency: "PHP",
	Receiver: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "John",
			MiddleName: "Michael",
			LastName:   "Doe",
			Phone: &ppb.PhoneNumber{
				Number: "09516738640",
			},
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				Province:   "prov",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "Zone",
			},
		},
		Identification: &ppb.Identification{
			Type:    "PASSPORT",
			Number:  "B83180608851",
			Country: "PH",
		},
		Employment: &tpb.Employment{
			Occupation: "OTH",
		},
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		TransactionPurpose: "Family Support/Living Expenses",
		PartnerMemberID:    "7712780",
		BirthCountry:       "PH",
		BirthPlace:         "MALOLOS,BULACAN",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
	},
	Agent: &tpb.Agent{
		UserID:    1893,
		IPAddress: "130.211.2.203",
	},
	Transaction: &tpb.Transaction{
		SourceCountry:      "PH",
		DestinationCountry: "PH",
	},
}

var ayaCreateReq = &tpb.CreateRemitRequest{
	RemitPartner: static.AYACode,
	RemitType:    "Send",
	OrderID:      "5142096965",
	Amount: &tpb.SendAmount{
		Amount:              "1000",
		SourceCountry:       "PH",
		SourceCurrency:      "PHP",
		DestinationCurrency: "PHP",
		DestinationCountry:  "PH",
		DestinationAmount:   true,
	},
	Receiver: &tpb.Receiver{
		ContactInfo: &tpb.Contact{
			FirstName:  "lkhds",
			MiddleName: "asdksad",
			LastName:   "sdjgsajdh",
			Email:      "emai@email.com",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "123455673",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "123458983",
			},
		},
	},
	Remitter: &tpb.UserKYC{
		ContactInfo: &tpb.Contact{
			FirstName:  "asdkhkh",
			MiddleName: "alsdkdfb",
			LastName:   "skdhfid",
			Address: &tpb.Address{
				Address1:   "addr1",
				Address2:   "addr2",
				City:       "city",
				State:      "state",
				Province:   "province",
				PostalCode: "12345",
				Country:    "PH",
				Zone:       "zone",
			},
			Phone: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12340967",
			},
			Mobile: &ppb.PhoneNumber{
				CountryCode: "62",
				Number:      "12340891",
			},
		},
		PartnerMemberID:    "77125891",
		BirthCountry:       "PH",
		BirthPlace:         "bit",
		SourceFunds:        "Salary/Income",
		ReceiverRelation:   "Family",
		TransactionPurpose: "Gift",
		Birthdate: &ppb.Date{
			Year:  "1950",
			Month: "12",
			Day:   "12",
		},
		Employment: &tpb.Employment{
			Occupation: "Unemployed",
		},
		Identification: &ppb.Identification{
			Type:    "PASSPORT",
			Country: "PH",
			Number:  "24023497AB0877AAB895668",
			Issued: &ppb.Date{
				Year:  "2020",
				Month: "12",
				Day:   "12",
			},
			Expiration: &ppb.Date{
				Year:  "2031",
				Month: "12",
				Day:   "12",
			},
		},
	},
	Agent: &tpb.Agent{
		UserID:    1797,
		IPAddress: "130.211.2.002",
	},
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func getError(err error) string {
	st := status.Convert(err)
	for _, detail := range st.Details() {
		e := Error{}
		b, err := json.Marshal(detail)
		if err != nil {
			fmt.Println("marshaling", err)
		}
		if err := json.Unmarshal(b, &e); err != nil {
			fmt.Println("unmarshaling", err)
		}
		return fmt.Sprintf("[partner error] code: %v, msg: %v", e.Code, e.Message)
	}
	return st.Message()
}
