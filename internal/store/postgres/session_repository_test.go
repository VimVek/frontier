package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/raystack/frontier/core/authenticate/session"
	frontiersession "github.com/raystack/frontier/core/authenticate/session"
	"github.com/raystack/frontier/core/user"
	"github.com/raystack/frontier/internal/store/postgres"
	"github.com/raystack/frontier/pkg/db"
	"github.com/raystack/frontier/pkg/metadata"
	"github.com/raystack/salt/log"
	"github.com/stretchr/testify/suite"
)

type SessionRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.SessionRepository
	session    []session.Session
	users      []user.User
}

func (s *SessionRepositoryTestSuite) SetupSuite() {
	var err error
	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewSessionRepository(logger, s.client)
	s.users, err = bootstrapUser(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *SessionRepositoryTestSuite) SetupTest() {
	var err error
	s.session, err = bootstrapSession(s.client, s.users)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *SessionRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *SessionRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *SessionRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_SESSIONS),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *SessionRepositoryTestSuite) TestSet() {
	type testCase struct {
		Name             string
		ExpectedSessison session.Session
		ExpectedError    error
	}
	var testCases = []testCase{
		{
			Name: "should set session",
			ExpectedSessison: session.Session{
				UserID: s.users[0].ID,
				Metadata: metadata.Metadata{
					"key": "value",
				},
			},
			ExpectedError: nil,
		},
		{
			Name: "should give error if id is empty",
			ExpectedSessison: session.Session{
				UserID: "",
			},
			ExpectedError: postgres.ErrInvalidID,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.Set(s.ctx, &tc.ExpectedSessison)
			if err != nil && !errors.Is(err, tc.ExpectedError) {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
		})
	}
}

func (s *SessionRepositoryTestSuite) TestGet() {
	type testCase struct {
		Name          string
		ExpextedId    uuid.UUID
		ExpectedResp  *session.Session
		ExpectedError error
	}
	var testCases = []testCase{
		{
			Name:       "should GET session",
			ExpextedId: s.session[0].ID,
			ExpectedResp: &session.Session{
				ID:     s.session[0].ID,
				UserID: s.users[0].ID,
				Metadata: metadata.Metadata{
					"some_key": "some_value",
				},
			},
			ExpectedError: nil,
		},
		{
			Name:          "should give error",
			ExpextedId:    uuid.Nil,
			ExpectedResp:  nil,
			ExpectedError: frontiersession.ErrNoSession,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Get(s.ctx, tc.ExpextedId)
			if err != nil && !errors.Is(err, tc.ExpectedError) {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
			if !cmp.Equal(resp, tc.ExpectedResp, cmpopts.IgnoreFields(session.Session{}, "AuthenticatedAt", "CreatedAt", "ExpiresAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResp)
			}
		})
	}
}

func (s *SessionRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Name          string
		ExpectedId    uuid.UUID
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name:          "should delete session",
			ExpectedId:    s.session[0].ID,
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.Delete(s.ctx, tc.ExpectedId)
			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
		})
	}
}

func (s *SessionRepositoryTestSuite) TestDeleteExpiredSessions() {
	type testCase struct {
		Name          string
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name:          "should delete Expired session",
			ExpectedError: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.DeleteExpiredSessions(s.ctx)
			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
		})
	}
}

func TestSessionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SessionRepositoryTestSuite))
}
