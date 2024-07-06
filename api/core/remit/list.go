package remit

import (
	"context"
	"fmt"

	"brank.as/petnet/api/core"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
)

// ListRemit historical transactions.
func (s *Svc) ListRemit(ctx context.Context, f core.FilterList) (*core.SearchRemitResponse, error) {
	log := logging.FromContext(ctx)
	var err error
	fl := storage.LRHFilter{
		TxnStep:        string(storage.ConfirmStep),
		TxnStatus:      string(storage.SuccessStatus),
		From:           f.From,
		Until:          f.Until,
		Limit:          f.Limit,
		Offset:         f.Offset,
		SortOrder:      storage.SortOrder(f.SortOrder),
		SortByColumn:   storage.Column(f.SortByColumn),
		ControlNo:      f.ControlNumbers,
		ExcludePartner: f.ExcludePartner,
		ExcludeType:    f.ExcludeType,
	}
	fl.DsaOrgID = getDSAOrgID(ctx)
	if fl.DsaOrgID == "" {
		log.Error("missing orgid")
		return nil, fmt.Errorf("processing")
	}

	rhs, err := s.st.ListRemitHistory(ctx, fl)
	if err != nil {
		log.Error(err)
	}

	lst := []core.SearchRemit{}
	for _, rh := range rhs {
		rm := rh.Remittance
		rmt := core.SearchRemit{
			RemitPartner: rh.RemcoID,
			RemitType:    rh.RemType,
			DestCurrency: "",
			ControlNo:    rh.RemcoControlNo,
			Remitter: core.Contact{
				FirstName:  rm.Remitter.FirstName,
				MiddleName: rm.Remitter.MiddleName,
				LastName:   rm.Remitter.LastName,
				Email:      rm.Remitter.Email,
				Address: core.Address{
					Address1:   rm.Remitter.Address1,
					Address2:   rm.Remitter.Address2,
					City:       rm.Remitter.City,
					State:      rm.Remitter.State,
					Province:   rm.Remitter.Province,
					PostalCode: rm.Remitter.PostalCode,
					Country:    rm.Remitter.Country,
					Zone:       rm.Remitter.Zone,
				},
				Phone: core.PhoneNumber{
					CtyCode: rm.Remitter.PhoneCty,
					Number:  rm.Remitter.Phone,
				},
				Mobile: core.PhoneNumber{
					CtyCode: rm.Remitter.MobileCty,
					Number:  rm.Remitter.Mobile,
				},
			},
			Receiver: core.Contact{
				FirstName:  rm.Receiver.FirstName,
				MiddleName: rm.Receiver.MiddleName,
				LastName:   rm.Receiver.LastName,
				Email:      rm.Receiver.Email,
				Address: core.Address{
					Address1:   rm.Receiver.Address1,
					Address2:   rm.Receiver.Address2,
					City:       rm.Receiver.City,
					State:      rm.Receiver.State,
					Province:   rm.Receiver.Province,
					PostalCode: rm.Receiver.PostalCode,
					Country:    rm.Receiver.Country,
					Zone:       rm.Receiver.Zone,
				},
				Phone: core.PhoneNumber{
					CtyCode: rm.Receiver.PhoneCty,
					Number:  rm.Receiver.Phone,
				},
				Mobile: core.PhoneNumber{
					CtyCode: rm.Receiver.MobileCty,
					Number:  rm.Receiver.Mobile,
				},
			},
			RemitAmount:      rm.GrossTotal.Minor,
			DisburseAmount:   rm.SourceAmt,
			Tax:              rm.Tax,
			Charges:          rm.Charges,
			Charge:           rm.Charge,
			TxnStagedTime:    rh.TxnStagedTime.Time,
			TxnCompletedTime: rh.TxnCompletedTime.Time,
		}
		lst = append(lst, rmt)
	}

	tot := 0
	if len(rhs) > 0 {
		tot = rhs[0].Total
	}

	return &core.SearchRemitResponse{
		Total:        tot,
		SearchRemits: lst,
	}, nil
}

func getDSAOrgID(ctx context.Context) string {
	ot := phmw.GetOrgType(ctx)
	switch ot {
	// this happens when API is used internally and means that either the dsa
	// or admin platform is used to get transactions
	case ppb.OrgType_PetNet.String(), ppb.OrgType_DSA.String():
		return phmw.GetDSAOrgID(ctx)
	}
	// this happens when API is used externally and means the user authenticated
	// with api client credentials to get token
	return hydra.OrgID(ctx)
}
