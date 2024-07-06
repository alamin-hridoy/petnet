package client

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"google.golang.org/grpc"

	lpb "brank.as/rbac/serviceutil/leaderelex/gunk/v1/lead"
	"brank.as/rbac/serviceutil/mainpkg"
)

type Leader struct {
	cl    lpb.LeaderElexServiceClient
	close func() error
}

func WithElector() mainpkg.Option {
	return func(conf *mainpkg.Config) {
		conf.LeadElector = func(p string) (mainpkg.Leader, error) {
			l, err := NewLeader(p)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize leader election %w", err)
			}
			return l, nil
		}
	}
}

func NewLeader(path string) (*Leader, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	u := &url.URL{Scheme: "unix", Path: p}
	conn, err := grpc.Dial(u.String(), grpc.WithInsecure(),
		grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		return nil, err
	}
	return &Leader{cl: lpb.NewLeaderElexServiceClient(conn), close: conn.Close}, nil
}

// Close the connection.
func (s *Leader) Close() error { return s.close() }

// IsLead helper function reports current leader status.
// If an error occurs, defaults false to avoid split-brain.
func (s *Leader) IsLead(ctx context.Context) bool {
	l, err := s.cl.GetLeader(ctx, &lpb.LeaderRequest{})
	if err != nil {
		return false
	}
	return l.Lead
}

func (s *Leader) GetLead(ctx context.Context) (string, bool, error) {
	l, err := s.cl.GetLeader(ctx, &lpb.LeaderRequest{})
	return l.GetID(), l.GetLead(), err
}
