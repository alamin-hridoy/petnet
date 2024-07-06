package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

// GetQuestion return Question matched against org ID & user id & question id
func (s *Storage) GetQuestion(ctx context.Context, req *storage.Question) (*storage.Question, error) {
	if req.OrgID != "" && req.UserID != "" && req.QID != "" {
		const QuestionQuery = `SELECT * FROM risk_assesment_question WHERE org_id = $1 and user_id = $2 and qid = $3`
		var qsn storage.Question
		if err := s.db.Get(&qsn, QuestionQuery, req.OrgID, req.UserID, req.QID); err != nil {
			if err == sql.ErrNoRows {
				return nil, storage.NotFound
			}
			return nil, err
		}
		return &qsn, nil
	}
	return nil, storage.InvalidArgument
}

// GetMlTfQuestion return MLTF Question matched against org ID & user id & question id
func (s *Storage) GetMlTfQuestion(ctx context.Context, req *storage.Question) (*storage.Question, error) {
	if req.OrgID != "" && req.UserID != "" && req.QID != "" {
		const QuestionQuery = `SELECT * FROM mltf_risk_assesment_question WHERE org_id = $1 and user_id = $2 and qid = $3`
		var qsn storage.Question
		if err := s.db.Get(&qsn, QuestionQuery, req.OrgID, req.UserID, req.QID); err != nil {
			if err == sql.ErrNoRows {
				return nil, storage.NotFound
			}
			return nil, err
		}
		return &qsn, nil
	}
	return nil, storage.InvalidArgument
}

// ListQuestion return Question List matched against org ID & user id
func (s *Storage) ListQuestion(ctx context.Context, req *storage.Question) ([]storage.Question, error) {
	if req.OrgID != "" && req.UserID != "" {
		QuestionQuery := fmt.Sprintf(`SELECT * FROM risk_assesment_question WHERE org_id = '%s' AND user_id = '%s'`, req.OrgID, req.UserID)
		var qsn []storage.Question
		if err := s.db.Select(&qsn, QuestionQuery); err != nil {
			if err == sql.ErrNoRows {
				return nil, storage.NotFound
			}
			return nil, err
		}
		return qsn, nil
	}
	return nil, storage.InvalidArgument
}

// ListMlTfQuestion return MLTF Question List matched against org ID & user id
func (s *Storage) ListMlTfQuestion(ctx context.Context, req *storage.Question) ([]storage.Question, error) {
	if req.OrgID != "" && req.UserID != "" {
		QuestionQuery := `SELECT * FROM mltf_risk_assesment_question WHERE org_id = $1 AND user_id = $2`
		var qsn []storage.Question
		if err := s.db.Select(&qsn, QuestionQuery, req.OrgID, req.UserID); err != nil {
			if err == sql.ErrNoRows {
				return nil, storage.NotFound
			}
			return nil, err
		}
		return qsn, nil
	}
	return nil, storage.InvalidArgument
}

// CreateQuestion creates new question and returns the created question
func (s *Storage) CreateQuestion(ctx context.Context, req *storage.Question) (*storage.Question, error) {
	log := logging.FromContext(ctx)
	const insertQuestion = `INSERT INTO risk_assesment_question (user_id, org_id, qid, ans, qtype) VALUES (:user_id, :org_id, :qid, :ans, :qtype) RETURNING *`
	pstmt, err := s.db.PrepareNamedContext(ctx, insertQuestion)
	if err != nil {
		logging.WithError(err, log).Error("insert question")
		return nil, err
	}
	defer pstmt.Close()
	if err := pstmt.Get(req, req); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing question insert: %w", err)
	}
	return req, nil
}

// Create MLTF Question creates new MLTF question and returns the created MLTF question
func (s *Storage) CreateMlTfQuestion(ctx context.Context, req *storage.Question) (*storage.Question, error) {
	log := logging.FromContext(ctx)
	const insertMlTfQuestion = `INSERT INTO mltf_risk_assesment_question (user_id, org_id, qid, qtype, customers_total, hr_total, impact_score) VALUES (:user_id, :org_id, :qid, :qtype, :customers_total, :hr_total, :impact_score) RETURNING *`
	pstmt, err := s.db.PrepareNamedContext(ctx, insertMlTfQuestion)
	if err != nil {
		logging.WithError(err, log).Error("insert mltf question")
		return nil, err
	}
	defer pstmt.Close()
	if err := pstmt.Get(req, req); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing mltf question insert: %w", err)
	}
	return req, nil
}

// Update Question update existing question and returns the question
func (s *Storage) UpdateQuestion(ctx context.Context, req *storage.Question) (*storage.Question, error) {
	const updateQuestion = `UPDATE risk_assesment_question SET ans= COALESCE(NULLIF(:ans, ''), ans) WHERE org_id=:org_id AND user_id=:user_id AND qid=:qid RETURNING *`
	stmt, err := s.db.PrepareNamedContext(ctx, updateQuestion)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(req, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return req, nil
}

// Update MLTF Question update existing MLTF question and returns the MLTF question
func (s *Storage) UpdateMlTfQuestion(ctx context.Context, req *storage.Question) (*storage.Question, error) {
	const updateMlTfQuestion = `UPDATE mltf_risk_assesment_question SET customers_total= COALESCE(NULLIF(:customers_total, ''), customers_total), hr_total= COALESCE(NULLIF(:hr_total, ''), hr_total), impact_score= COALESCE(NULLIF(:impact_score, ''), impact_score) WHERE org_id=:org_id AND user_id=:user_id AND qid=:qid RETURNING *`
	stmt, err := s.db.PrepareNamedContext(ctx, updateMlTfQuestion)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(req, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return req, nil
}

// Upsert Question create and update existing question and returns the question
func (s *Storage) UpsertQuestion(ctx context.Context, req *storage.Question) (*storage.Question, error) {
	if req.OrgID != "" && req.UserID != "" && req.QID != "" {
		if _, err := s.GetQuestion(ctx, req); err != nil {
			res, err := s.CreateQuestion(ctx, req)
			if err != nil {
				return nil, err
			}
			return res, nil
		} else {
			res, err := s.UpdateQuestion(ctx, req)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
	return nil, storage.InvalidArgument
}

// Upsert Question create and update existing question and returns the question
func (s *Storage) UpsertMlTfQuestion(ctx context.Context, req *storage.Question) (*storage.Question, error) {
	if req.OrgID != "" && req.UserID != "" && req.QID != "" {
		if _, err := s.GetMlTfQuestion(ctx, req); err != nil {
			res, err := s.CreateMlTfQuestion(ctx, req)
			if err != nil {
				return nil, err
			}
			return res, nil
		} else {
			res, err := s.UpdateMlTfQuestion(ctx, req)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
	}
	return nil, storage.InvalidArgument
}
