package microinsurance

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
	"brank.as/petnet/serviceutil/logging"
)

// Transact ...
func (s *Svc) Transact(ctx context.Context, req *migunk.TransactRequest) (*migunk.Insurance, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if err := validateTransactRequest(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.store.Transact(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

// GetReprint ...
func (s *Svc) GetReprint(ctx context.Context, req *migunk.GetReprintRequest) (*migunk.Insurance, error) {
	if req == nil || req.TraceNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "trace number is required")
	}

	res, err := s.store.GetReprint(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

// RetryTransaction ...
func (s *Svc) RetryTransaction(ctx context.Context, req *migunk.RetryTransactionRequest) (*migunk.Insurance, error) {
	if req == nil || req.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	res, err := s.store.RetryTransaction(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

// GetTransactionList ...
func (s *Svc) GetTransactionList(ctx context.Context, req *migunk.GetTransactionListRequest) (*migunk.TransactionListResult, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "date filter is required")
	}

	err := validation.ValidateStruct(req,
		validation.Field(&req.DateFrom, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&req.DateTo, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&req.OrgID, validation.Required, is.UUID),
	)
	if err != nil {
		return nil, err
	}

	res, err := s.store.GetTransactionList(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func validateTransactRequest(ctx context.Context, r *migunk.TransactRequest) error {
	log := logging.FromContext(ctx)
	err := validation.ValidateStruct(
		r,
		validation.Field(&r.Coy, validation.Required),
		validation.Field(&r.LocationID, validation.Required),
		validation.Field(&r.UserCode, validation.Required),
		validation.Field(&r.TrxDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&r.Amount, validation.Required, is.Float),
		validation.Field(&r.CoverageCount, validation.Required, is.Digit),
		validation.Field(&r.ProductCode, validation.Required, is.Alphanumeric),
		validation.Field(&r.ProcessingBranch, validation.Required),
		validation.Field(&r.ProcessedBy, validation.Required),
		validation.Field(&r.UserEmail, validation.Required, is.Email),
		validation.Field(&r.LastName, validation.Required),
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.Gender, validation.Required),
		validation.Field(&r.Birthdate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&r.MobileNumber, validation.Required, is.Digit),
		validation.Field(&r.ProvinceCode, validation.Required),
		validation.Field(&r.CityCode, validation.Required),
		validation.Field(&r.Address, validation.Required),
		validation.Field(&r.MaritalStatus, validation.Required),
		validation.Field(&r.Occupation, validation.Required),
		validation.Field(&r.NumberUnits, validation.Required, is.Digit),
	)
	if err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if len(r.Beneficiaries) == 0 {
		log.Error("beneficiaries are empty")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if len(r.Dependents) == 0 {
		log.Error("dependents are empty")
		return status.Error(codes.InvalidArgument, err.Error())
	}

	err = validateMicroinsurancePersons(r.Beneficiaries)
	if err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(codes.InvalidArgument, "Beneficiary "+err.Error())
	}

	err = validateMicroinsurancePersons(r.Dependents)
	if err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(codes.InvalidArgument, "Dependent "+err.Error())
	}

	return nil
}

func validateMicroinsurancePersons(persons []*migunk.Person) error {
	for _, ben := range persons {
		err := validation.ValidateStruct(
			ben,
			validation.Field(&ben.FirstName, validation.Required),
			validation.Field(&ben.LastName, validation.Required),
			validation.Field(&ben.ContactNumber, validation.Required, is.Digit),
			validation.Field(&ben.BirthDate, validation.Required, validation.Date("2006-01-02")),
			validation.Field(&ben.Relationship, validation.Required),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
