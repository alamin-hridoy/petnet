package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testWUTHBody struct {
	Module  string      `json:"module"`
	Request string      `json:"request"`
	Param   WUTHRequest `json:"param"`
}

type testWUTH struct {
	Header    RequestHeader `json:"header"`
	Body      testWUTHBody  `json:"body"`
	Signature string        `json:"signature"`
}

type testWUTHRequest struct {
	WU testWUTH `json:"uspwuapi"`
}

var thWUReq = WUTHRequest{
	DateStart: "2017-07-22 00:00:00",
	DateEnd:   "2017-09-20 23:59:59",
}

func TestWUTransactionHistory(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          WUTHRequest
		expectedReq testWUTHRequest
		want        []WUTransaction
		wantErr     bool
	}{
		{
			name: "Success",
			in:   thWUReq,
			expectedReq: testWUTHRequest{
				WU: testWUTH{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testWUTHBody{
						Module:  "wurpt",
						Request: "report",
						Param:   thWUReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: []WUTransaction{
				{
					TxnDate:      "2021-03-01 16:01:50",
					MTCN:         "9429067484",
					Principal:    "2127",
					SvcFee:       "0",
					Currency:     "PHP",
					TxnType:      "PO",
					DateClaimed:  "2021-03-01 03:01:00",
					CustomerCode: "1083182       ",
					OrderID:      "WUB8060601202103",
				},
				{
					TxnDate:      "2021-03-01 16:56:32",
					MTCN:         "9409183168",
					Principal:    "2127",
					SvcFee:       "0",
					Currency:     "PHP",
					TxnType:      "PO",
					DateClaimed:  "2021-03-01 03:56:00",
					CustomerCode: "1083182       ",
					OrderID:      "WU2DB0ED01202103",
				},
				{
					TxnDate:      "2021-03-02 15:51:37",
					MTCN:         "9709183168",
					Principal:    "2127",
					SvcFee:       "0",
					Currency:     "PHP",
					TxnType:      "PO",
					DateClaimed:  "2021-03-02 02:51:00",
					CustomerCode: "1083182       ",
					OrderID:      "WUBF6C4C02202103",
				},
				{
					TxnDate:      "2021-03-04 16:57:34",
					MTCN:         "9489067484",
					Principal:    "2133",
					SvcFee:       "0",
					Currency:     "PHP",
					TxnType:      "PO",
					DateClaimed:  "2021-03-04 03:57:00",
					CustomerCode: "1083182       ",
					OrderID:      "WUD1181C04202103",
				},
				{
					TxnDate:      "2021-03-08 16:37:43",
					MTCN:         "1159106049",
					Principal:    "2127",
					SvcFee:       "0",
					Currency:     "PHP",
					TxnType:      "PO",
					DateClaimed:  "2021-03-08 03:37:00",
					CustomerCode: "1086960       ",
					OrderID:      "WU6C5B3F08202103",
				},
				{
					TxnDate:      "2021-03-09 18:37:57",
					MTCN:         "9119226553",
					Principal:    "400",
					SvcFee:       "0",
					Currency:     "PHP",
					TxnType:      "PO",
					DateClaimed:  "2021-03-09 05:37:00",
					CustomerCode: "1076911       ",
					OrderID:      "WUB7E99009202103",
				},
				{
					TxnDate:      "2021-03-09 19:04:46",
					MTCN:         "6639167456",
					Principal:    "500",
					SvcFee:       "0",
					Currency:     "PHP",
					TxnType:      "PO",
					DateClaimed:  "2021-03-09 06:04:00",
					CustomerCode: "1076911       ",
					OrderID:      "WU02D82509202103",
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successWUTHHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.WUTransactionHistory(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("WUTransactionHistory() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func successWUTHHandler(t *testing.T, expectedReq testWUTHRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wu-rpt" {
			t.Errorf("expected request to '/1.1/wu-rpt', got '%s'", req.URL.EscapedPath())
		}

		var newReq testWUTHRequest
		if err := json.NewDecoder(req.Body).Decode(&newReq); err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(expectedReq, newReq) {
			t.Error(cmp.Diff(expectedReq, newReq))
		}

		res.WriteHeader(200)
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, `{
				"status":"1",
				"errmsg":"Success",
				"data":
				[{
			      "transactiondate": "2021-03-01 16:01:50",
			      "MTCN": "9429067484",
			      "principalamount": "2127",
			      "servicefee": 0,
			      "currency": "PHP",
			      "TransactionType": "PO",
			      "DateClaimed": "2021-03-01 03:01:00",
			      "CustomerCode": "1083182       ",
			      "OrderID": "WUB8060601202103"
			    },
			    {
			      "transactiondate": "2021-03-01 16:56:32",
			      "MTCN": "9409183168",
			      "principalamount": "2127",
			      "servicefee": 0,
			      "currency": "PHP",
			      "TransactionType": "PO",
			      "DateClaimed": "2021-03-01 03:56:00",
			      "CustomerCode": "1083182       ",
			      "OrderID": "WU2DB0ED01202103"
			    },
			    {
			      "transactiondate": "2021-03-02 15:51:37",
			      "MTCN": "9709183168",
			      "principalamount": "2127",
			      "servicefee": 0,
			      "currency": "PHP",
			      "TransactionType": "PO",
			      "DateClaimed": "2021-03-02 02:51:00",
			      "CustomerCode": "1083182       ",
			      "OrderID": "WUBF6C4C02202103"
			    },
			    {
			      "transactiondate": "2021-03-04 16:57:34",
			      "MTCN": "9489067484",
			      "principalamount": "2133",
			      "servicefee": 0,
			      "currency": "PHP",
			      "TransactionType": "PO",
			      "DateClaimed": "2021-03-04 03:57:00",
			      "CustomerCode": "1083182       ",
			      "OrderID": "WUD1181C04202103"
			    },
			    {
			      "transactiondate": "2021-03-08 16:37:43",
			      "MTCN": "1159106049",
			      "principalamount": "2127",
			      "servicefee": 0,
			      "currency": "PHP",
			      "TransactionType": "PO",
			      "DateClaimed": "2021-03-08 03:37:00",
			      "CustomerCode": "1086960       ",
			      "OrderID": "WU6C5B3F08202103"
			    },
			    {
			      "transactiondate": "2021-03-09 18:37:57",
			      "MTCN": "9119226553",
			      "principalamount": "400",
			      "servicefee": 0,
			      "currency": "PHP",
			      "TransactionType": "PO",
			      "DateClaimed": "2021-03-09 05:37:00",
			      "CustomerCode": "1076911       ",
			      "OrderID": "WUB7E99009202103"
			    },
			    {
			      "transactiondate": "2021-03-09 19:04:46",
			      "MTCN": "6639167456",
			      "principalamount": "500",
			      "servicefee": 0,
			      "currency": "PHP",
			      "TransactionType": "PO",
			      "DateClaimed": "2021-03-09 06:04:00",
			      "CustomerCode": "1076911       ",
			      "OrderID": "WU02D82509202103"
			    }]
		}`)
	}
}
