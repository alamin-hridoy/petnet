package revenue_commission

import (
	"context"
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	revcom_int "brank.as/petnet/api/integration/revenue-commission"
	"brank.as/petnet/api/util"
	"brank.as/petnet/gunk/drp/v1/dsa"
	"brank.as/petnet/serviceutil/logging"
)

var errUnImplemented = status.New(codes.Unimplemented, "not implemented").Err()

// CreateDSA Creates DSA record.
func (s *Svc) CreateDSA(ctx context.Context, req *dsa.CreateDSARequest) (*dsa.DSA, error) {
	log := logging.FromContext(ctx)
	createReq, err := createDSAValidate(req)
	if err != nil {
		logging.WithError(err, log).Error("failed validating request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.dsaStore.CreateDSA(ctx, createReq)

	return toDSAResponse(res, err, log)
}

// UpdateDSA Updates DSA record by id.
func (s *Svc) UpdateDSA(ctx context.Context, req *dsa.UpdateDSARequest) (*dsa.DSA, error) {
	log := logging.FromContext(ctx)
	updateReq, err := updateDSAValidate(req)
	if err != nil {
		logging.WithError(err, log).Error("failed validating request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.dsaStore.UpdateDSA(ctx, updateReq)

	return toDSAResponse(res, err, log)
}

// GetDSAByID Gets DSA record by ID.
func (s *Svc) GetDSAByID(ctx context.Context, req *dsa.GetDSAByIDRequest) (*dsa.DSA, error) {
	log := logging.FromContext(ctx)
	if req.DsaID == 0 {
		log.Error("dsa id required")
		return nil, util.HandleServiceErr(status.New(codes.InvalidArgument, "dsa id required").Err())
	}

	res, err := s.dsaStore.GetDSAByID(ctx, req.DsaID)

	return toDSAResponse(res, err, log)
}

// ListDSA Gets all DSA records
func (s *Svc) ListDSA(ctx context.Context, req *emptypb.Empty) (*dsa.ListDSAResponse, error) {
	log := logging.FromContext(ctx)

	res, err := s.dsaStore.ListDSA(ctx)

	return toListDSAResponse(res, err, log)
}

// DeleteDSAByID Deletes DSA record by ID.
func (s *Svc) DeleteDSAByID(ctx context.Context, req *dsa.DeleteDSAByIDRequest) (*dsa.DSA, error) {
	log := logging.FromContext(ctx)
	if req.DsaID == 0 {
		log.Error("dsa id required")
		return nil, util.HandleServiceErr(status.New(codes.InvalidArgument, "dsa id required").Err())
	}

	res, err := s.dsaStore.DeleteDSA(ctx, req.DsaID)

	return toDSAResponse(res, err, log)
}

func createDSAValidate(req *dsa.CreateDSARequest) (*revcom_int.SaveDSARequest, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.DsaCode, validation.Required),
		validation.Field(&req.DsaName, validation.Required),
		validation.Field(&req.EmailAddress, validation.Required, is.Email),
		validation.Field(&req.Vatable, validation.Required, is.Digit),
		validation.Field(&req.Address, validation.Required),
		validation.Field(&req.Tin, validation.Required),
		validation.Field(&req.UpdatedBy, validation.Required),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &revcom_int.SaveDSARequest{
		DsaCode:        req.DsaCode,
		DsaName:        req.DsaName,
		EmailAddress:   req.EmailAddress,
		Vatable:        json.Number(req.Vatable),
		Address:        req.Address,
		Tin:            req.Tin,
		UpdatedBy:      req.UpdatedBy,
		ContactPerson:  req.ContactPerson,
		City:           req.City,
		Province:       req.Province,
		Zipcode:        req.Zipcode,
		President:      req.President,
		GeneralManager: req.GeneralManager,
	}, nil
}

func updateDSAValidate(req *dsa.UpdateDSARequest) (*revcom_int.SaveDSARequest, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.DsaID, validation.Required),
		validation.Field(&req.DsaCode, validation.Required),
		validation.Field(&req.DsaName, validation.Required),
		validation.Field(&req.EmailAddress, validation.Required, is.Email),
		validation.Field(&req.Vatable, validation.Required, is.Digit),
		validation.Field(&req.Address, validation.Required),
		validation.Field(&req.Tin, validation.Required),
		validation.Field(&req.UpdatedBy, validation.Required),
		validation.Field(&req.ContactPerson, validation.Required),
		validation.Field(&req.City, validation.Required),
		validation.Field(&req.Province, validation.Required),
		validation.Field(&req.Zipcode, validation.Required),
		validation.Field(&req.President, validation.Required),
		validation.Field(&req.GeneralManager, validation.Required),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &revcom_int.SaveDSARequest{
		DsaID:          req.DsaID,
		DsaCode:        req.DsaCode,
		DsaName:        req.DsaName,
		EmailAddress:   req.EmailAddress,
		Vatable:        json.Number(req.Vatable),
		Address:        req.Address,
		Tin:            req.Tin,
		UpdatedBy:      req.UpdatedBy,
		ContactPerson:  req.ContactPerson,
		City:           req.City,
		Province:       req.Province,
		Zipcode:        req.Zipcode,
		President:      req.President,
		GeneralManager: req.GeneralManager,
	}, nil
}

func toDSAResponse(d *revcom_int.DSA, err error, log *logrus.Entry) (*dsa.DSA, error) {
	if err != nil {
		logging.WithError(err, log).Error("dsa store error")
		return nil, util.HandleServiceErr(err)
	}

	if d == nil {
		return nil, util.HandleServiceErr(status.New(codes.Unknown, "empty dsa response").Err())
	}

	var (
		dsaID, status, vatable          int64
		createdAt, updatedAt, deletedAt *timestamppb.Timestamp
	)

	dsaID, err = d.DsaID.Int64()
	if err != nil {
		logging.WithError(err, log).Error("dsa id is not numeric")
		return nil, util.HandleServiceErr(err)
	}

	status, _ = d.Status.Int64()
	vatable, _ = d.Vatable.Int64()

	if d.CreatedAt != nil {
		createdAt = timestamppb.New(*d.CreatedAt)
	}

	if d.UpdatedAt != nil {
		updatedAt = timestamppb.New(*d.UpdatedAt)
	}

	if d.DeletedAt != nil {
		deletedAt = timestamppb.New(*d.DeletedAt)
	}

	return &dsa.DSA{
		DsaID:          uint32(dsaID),
		DsaCode:        d.DsaCode,
		DsaName:        d.DsaName,
		EmailAddress:   d.EmailAddress,
		Status:         uint32(status),
		Vatable:        uint32(vatable),
		Address:        d.Address,
		Tin:            d.Tin,
		UpdatedBy:      d.UpdatedBy,
		ContactPerson:  d.ContactPerson,
		City:           d.City,
		Province:       d.Province,
		Zipcode:        d.Zipcode,
		President:      d.President,
		GeneralManager: d.GeneralManager,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		DeletedAt:      deletedAt,
	}, nil
}

func toListDSAResponse(dl []revcom_int.DSA, err error, log *logrus.Entry) (*dsa.ListDSAResponse, error) {
	if err != nil {
		logging.WithError(err, log).Error("dsa store error")
		return nil, util.HandleServiceErr(err)
	}

	if len(dl) == 0 {
		return &dsa.ListDSAResponse{DSAList: []*dsa.DSA{}}, nil
	}

	dsaList := make([]*dsa.DSA, 0, len(dl))
	for _, d := range dl {
		dsa, _ := toDSAResponse(&d, nil, log)
		dsaList = append(dsaList, dsa)
	}

	return &dsa.ListDSAResponse{
		DSAList: dsaList,
	}, nil
}
