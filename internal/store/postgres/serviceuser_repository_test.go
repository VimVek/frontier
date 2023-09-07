package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/raystack/frontier/core/role"
	"github.com/raystack/frontier/core/serviceuser"
	"github.com/raystack/frontier/internal/store/postgres"
	"github.com/raystack/frontier/pkg/db"
	"github.com/raystack/salt/log"
	"github.com/stretchr/testify/suite"
)

type ServiceUserRepositoryTestSuite struct {
	suite.Suite
	ctx         context.Context
	client      *db.Client
	pool        *dockertest.Pool
	resource    *dockertest.Resource
	repository  *postgres.ServiceUserRepository
	serviceUser []serviceuser.ServiceUser
}

func (s *ServiceUserRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewServiceUserRepository(s.client)
}

func (s *ServiceUserRepositoryTestSuite) SetupTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
	var err error
	s.serviceUser, err = bootstrapServiceUser(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceUserRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceUserRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceUserRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_SERVICEUSER),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *ServiceUserRepositoryTestSuite) TestList() {
	type TestCase struct {
		Name           string
		ExpectedFilter serviceuser.Filter
		ExpectedResp   []serviceuser.ServiceUser
		ExpectedErr    string
	}
	var testCases = []TestCase{
		{
			Name: "should List service users by using orgid as filter",
			ExpectedFilter: serviceuser.Filter{
				OrgID: s.serviceUser[0].OrgID,
			},
			ExpectedResp: []serviceuser.ServiceUser{
				{
					OrgID: s.serviceUser[0].OrgID,
					State: "enabled",
					Title: "title_01",
				},
			},
			ExpectedErr: "",
		},
		{
			Name: "should List service users by using state as filter",
			ExpectedFilter: serviceuser.Filter{
				State: serviceuser.Enabled,
			},
			ExpectedResp: []serviceuser.ServiceUser{
				{
					OrgID: s.serviceUser[0].OrgID,
					State: serviceuser.Enabled.String(),
					Title: "title_01",
				},
				{
					OrgID: s.serviceUser[1].OrgID,
					State: serviceuser.Enabled.String(),
					Title: "title_02",
				},
				{
					OrgID: s.serviceUser[2].OrgID,
					State: serviceuser.Enabled.String(),
					Title: "title_03",
				},
			},
			ExpectedErr: "",
		},
		{
			Name: "should give nil if id does not exist",
			ExpectedFilter: serviceuser.Filter{
				OrgID: uuid.New().String(),
			},
			ExpectedResp: nil,
			ExpectedErr:  "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.List(s.ctx, tc.ExpectedFilter)

			if tc.ExpectedErr != "" && err.Error() != tc.ExpectedErr {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedErr)
			}
			if !cmp.Equal(resp, tc.ExpectedResp, cmpopts.IgnoreFields(serviceuser.ServiceUser{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResp)
			}
		})
	}
}

func (s *ServiceUserRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Name          string
		SrvUser       serviceuser.ServiceUser
		ExpectedResp  serviceuser.ServiceUser
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name: "should create the service user",
			SrvUser: serviceuser.ServiceUser{
				OrgID: s.serviceUser[0].OrgID,
				State: serviceuser.Enabled.String(),
				Title: "title_03",
			},
			ExpectedResp: serviceuser.ServiceUser{
				OrgID: s.serviceUser[0].OrgID,
				State: serviceuser.Enabled.String(),
				Title: "title_03",
			},
			ExpectedError: "",
		},
		{
			Name:          "should create empty service user if nothing is provided in input",
			SrvUser:       serviceuser.ServiceUser{},
			ExpectedResp:  serviceuser.ServiceUser{},
			ExpectedError: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Create(s.ctx, tc.SrvUser)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedResp, cmpopts.IgnoreFields(serviceuser.ServiceUser{}, "ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResp)
			}
		})
	}
}

func (s *ServiceUserRepositoryTestSuite) TestGetByID() {
	type testCase struct {
		Name          string
		serviceUserId string
		ExpectedResp  serviceuser.ServiceUser
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name:          "should get the service user by its id",
			serviceUserId: s.serviceUser[0].ID,
			ExpectedResp: serviceuser.ServiceUser{
				OrgID: s.serviceUser[0].OrgID,
				State: serviceuser.Enabled.String(),
				Title: "title_01",
			},
			ExpectedError: "",
		},
		{
			Name:          "should give error if input id is empty",
			serviceUserId: "",
			ExpectedResp:  serviceuser.ServiceUser{},
			ExpectedError: role.ErrInvalidID.Error(),
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.GetByID(s.ctx, tc.serviceUserId)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedResp, cmpopts.IgnoreFields(serviceuser.ServiceUser{}, "ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResp)
			}
		})
	}
}

func (s *ServiceUserRepositoryTestSuite) TestGetByIDs() {
	type testCase struct {
		Name          string
		serviceUserId []string
		ExpectedResp  []serviceuser.ServiceUser
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name: "should List all service users whose id is provided as input",
			serviceUserId: []string{
				s.serviceUser[0].ID,
				s.serviceUser[1].ID,
			},
			ExpectedResp: []serviceuser.ServiceUser{
				{
					ID:    s.serviceUser[0].ID,
					OrgID: s.serviceUser[0].OrgID,
					State: serviceuser.Enabled.String(),
					Title: "title_01",
				},
				{
					ID:    s.serviceUser[1].ID,
					OrgID: s.serviceUser[1].OrgID,
					State: serviceuser.Enabled.String(),
					Title: "title_02",
				},
			},
			ExpectedError: "",
		},
		{
			Name:          "should give error if no input id is provided",
			serviceUserId: []string{},
			ExpectedResp:  nil,
			ExpectedError: role.ErrInvalidID.Error(),
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.GetByIDs(s.ctx, tc.serviceUserId)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedResp, cmpopts.IgnoreFields(serviceuser.ServiceUser{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResp)
			}
		})
	}
}

func (s *ServiceUserRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Name          string
		serviceUserId string
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name:          "should delete the service user",
			serviceUserId: s.serviceUser[0].ID,
			ExpectedError: "",
		},
		{
			Name:          "should give no error if service user id does not exist",
			serviceUserId: uuid.New().String(),
			ExpectedError: "",
		},
		{
			Name:          "should give no error if service user id is empty",
			serviceUserId: "",
			ExpectedError: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.Delete(s.ctx, tc.serviceUserId)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
		})
	}
}

func TestServiceUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceUserRepositoryTestSuite))
}
