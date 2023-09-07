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
	"github.com/raystack/frontier/core/authenticate"
	"github.com/raystack/frontier/internal/store/postgres"
	"github.com/raystack/frontier/pkg/db"
	"github.com/raystack/frontier/pkg/metadata"
	"github.com/raystack/salt/log"
	"github.com/stretchr/testify/suite"
)

type FlowRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.FlowRepository
	flow       []authenticate.Flow
}

func (s *FlowRepositoryTestSuite) SetupSuite() {
	var err error
	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx = context.TODO()
	s.repository = postgres.NewFlowRepository(logger, s.client)
}

func (s *FlowRepositoryTestSuite) SetupTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
	var err error
	s.flow, err = bootstrapFlow(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *FlowRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *FlowRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *FlowRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_FLOWS),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *FlowRepositoryTestSuite) TestSet() {
	type TestCase struct {
		Name        string
		FlowSet     *authenticate.Flow
		ExpectedErr string
	}
	var testCases = []TestCase{
		{
			Name: "should Set flow",
			FlowSet: &authenticate.Flow{
				Method: "testMethod05",
				Email:  "testmail05@123.com",
				Nonce:  "randomNum005",
			},
			ExpectedErr: "",
		},
		{
			Name: "should Set flow with metadata also given as input",
			FlowSet: &authenticate.Flow{
				Method: authenticate.MailOTPAuthMethod.String(),
				Email:  "testmail07@123.com",
				Nonce:  "randomNum007",
				Metadata: metadata.Metadata{
					"foo": "bar",
				},
			},
			ExpectedErr: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.Set(s.ctx, tc.FlowSet)
			if tc.ExpectedErr != "" {
				if err.Error() != tc.ExpectedErr {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedErr)
				}
			}
		})
	}
}

func (s *FlowRepositoryTestSuite) TestGet() {
	type testCase struct {
		Name           string
		SelectedID     uuid.UUID
		ExpectedInvite *authenticate.Flow
		ExpectedError  string
	}
	var testCases = []testCase{
		{
			Name:       "should Get the flow",
			SelectedID: s.flow[0].ID,
			ExpectedInvite: &authenticate.Flow{
				ID:     s.flow[0].ID,
				Method: authenticate.MailLinkAuthMethod.String(),
				Email:  "testmail@123.com",
				Nonce:  "randomNum01",
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Get(s.ctx, tc.SelectedID)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedInvite, cmpopts.IgnoreFields(authenticate.Flow{}, "Metadata", "CreatedAt", "ExpiresAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedInvite)
			}
		})
	}
}

func (s *FlowRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Name          string
		SelectedID    uuid.UUID
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name:          "should Delete the flow",
			SelectedID:    s.flow[0].ID,
			ExpectedError: "",
		},
		{
			Name:          "should Give error if id not present in database",
			SelectedID:    uuid.New(),
			ExpectedError: errors.New("no entry to delete").Error(),
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.Delete(s.ctx, tc.SelectedID)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
		})
	}
}

func (s *FlowRepositoryTestSuite) TestDeleteExpiredFlows() {
	type testCase struct {
		Name          string
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name:          "should Get the flow",
			ExpectedError: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.DeleteExpiredFlows(s.ctx)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
		})
	}
}

func TestFlowRespository(t *testing.T) {
	suite.Run(t, new(FlowRepositoryTestSuite))
}
