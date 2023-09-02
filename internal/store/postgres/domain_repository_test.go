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
	"github.com/raystack/frontier/core/domain"
	"github.com/raystack/frontier/core/organization"
	"github.com/raystack/frontier/internal/store/postgres"
	"github.com/raystack/frontier/pkg/db"
	"github.com/raystack/salt/log"
	"github.com/stretchr/testify/suite"
)

type DomainRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.DomainRepository
	domain     []domain.Domain
	orgs       []organization.Organization
}

func (s *DomainRepositoryTestSuite) SetupSuite() {
	var err error
	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewDomainRepository(logger, s.client)
	s.orgs, err = bootstrapOrganization(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *DomainRepositoryTestSuite) SetupTest() {
	var err error
	s.domain, err = bootstrapDomain(s.client, s.orgs)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *DomainRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *DomainRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *DomainRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_DOMAINS),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *DomainRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Name             string
		ExpectedDomain   domain.Domain
		ExpectedResponse domain.Domain
		ExpectedError    string
	}
	var testCases = []testCase{
		{
			Name: "should create a domain",
			ExpectedDomain: domain.Domain{
				Name:  "domain5",
				OrgID: s.orgs[0].ID,
				Token: "token05",
				State: domain.Pending,
			},
			ExpectedResponse: domain.Domain{
				Name:  "domain5",
				OrgID: s.orgs[0].ID,
				Token: "token05",
				State: domain.Pending,
			},
			ExpectedError: "",
		},
		{
			Name: "should give error if same domain name is used in a an org",
			ExpectedDomain: domain.Domain{
				Name:  "domain5",
				OrgID: s.orgs[0].ID,
				Token: "token05",
				State: domain.Pending,
			},
			ExpectedResponse: domain.Domain{},
			ExpectedError:    domain.ErrDuplicateKey.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Create(s.ctx, tc.ExpectedDomain)
			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
			if !cmp.Equal(resp, tc.ExpectedResponse, cmpopts.IgnoreFields(domain.Domain{}, "ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResponse)
			}
		})
	}
}

func (s *DomainRepositoryTestSuite) TestList() {
	type testCase struct {
		Name             string
		ExpectedDomain   domain.Filter
		ExpectedResponse []domain.Domain
		ExpectedError    string
	}
	var testCases = []testCase{
		{
			Name: "should List domains if orgid and state used as filter",
			ExpectedDomain: domain.Filter{
				OrgID: s.orgs[0].ID,
				State: domain.Pending,
			},
			ExpectedResponse: []domain.Domain{
				{
					Name:  "domain1",
					OrgID: s.orgs[0].ID,
					State: domain.Pending,
					Token: "token01",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "should List domains if name and state used as filter",
			ExpectedDomain: domain.Filter{
				Name:  "domain2",
				State: domain.Pending,
			},
			ExpectedResponse: []domain.Domain{
				{
					Name:  "domain2",
					OrgID: s.orgs[1].ID,
					State: domain.Pending,
					Token: "token03",
				},
			},
			ExpectedError: "",
		},
		{
			Name: "should List domains if only orgid is used as filter",
			ExpectedDomain: domain.Filter{
				OrgID: s.orgs[2].ID,
			},
			ExpectedResponse: []domain.Domain{
				{
					Name:  "domain3",
					OrgID: s.orgs[2].ID,
					State: domain.Pending,
					Token: "token04",
				},
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.List(s.ctx, tc.ExpectedDomain)
			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
			if !cmp.Equal(resp, tc.ExpectedResponse, cmpopts.IgnoreFields(domain.Domain{}, "ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResponse)
			}
		})
	}
}

func (s *DomainRepositoryTestSuite) TestGet() {
	type testCase struct {
		Name             string
		ExpectedId       string
		ExpectedResponse domain.Domain
		ExpectedError    error
	}
	var testCases = []testCase{
		{
			Name:       "should Get domain by using domain ID",
			ExpectedId: s.domain[0].ID,
			ExpectedResponse: domain.Domain{
				ID:    s.domain[0].ID,
				Name:  "domain1",
				OrgID: s.orgs[0].ID,
				State: domain.Pending,
				Token: "token01",
			},
			ExpectedError: nil,
		},
		{
			Name:             "should Give error if id is empty",
			ExpectedId:       "",
			ExpectedResponse: domain.Domain{},
			ExpectedError:    postgres.ErrQueryRun,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Get(s.ctx, tc.ExpectedId)
			if err != nil && !errors.Is(err, tc.ExpectedError) {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
			if !cmp.Equal(resp, tc.ExpectedResponse, cmpopts.IgnoreFields(domain.Domain{}, "ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResponse)
			}
		})
	}
}

func (s *DomainRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Name          string
		ExpectedId    string
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name:          "should delete a domain",
			ExpectedId:    s.domain[0].ID,
			ExpectedError: "",
		},
		{
			Name:          "should give error if domain does not exist",
			ExpectedId:    uuid.New().String(),
			ExpectedError: domain.ErrNotExist.Error(),
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

func (s *DomainRepositoryTestSuite) TestUpdate() {
	type testCase struct {
		Name             string
		ExpectedDomain   domain.Domain
		ExpectedResponse domain.Domain
		ExpectedError    string
	}
	var testCases = []testCase{
		{
			Name: "should update a domain",
			ExpectedDomain: domain.Domain{
				ID:    s.domain[0].ID,
				Name:  "domain1",
				OrgID: s.orgs[0].ID,
				Token: "token09",
				State: domain.Verified,
			},
			ExpectedResponse: domain.Domain{
				ID:    s.domain[0].ID,
				Name:  "domain1",
				OrgID: s.orgs[0].ID,
				Token: "token09",
				State: domain.Verified,
			},
			ExpectedError: "",
		},
		{
			Name: "should give error if domain id is not present",
			ExpectedDomain: domain.Domain{
				ID:    "",
				Name:  "domain1",
				OrgID: s.orgs[0].ID,
				Token: "token09",
			},
			ExpectedResponse: domain.Domain{},
			ExpectedError:    domain.ErrInvalidId.Error(),
		},
		{
			Name: "should give error if domain does not exist",
			ExpectedDomain: domain.Domain{
				ID:    uuid.New().String(),
				Name:  "domain.name",
				OrgID: s.orgs[0].ID,
				Token: "test-token",
				State: domain.Pending,
			},
			ExpectedResponse: domain.Domain{},
			ExpectedError:    domain.ErrNotExist.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Update(s.ctx, tc.ExpectedDomain)
			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
			if !cmp.Equal(resp, tc.ExpectedResponse, cmpopts.IgnoreFields(domain.Domain{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResponse)
			}
		})
	}
}

func (s *DomainRepositoryTestSuite) TestDeleteExpiredDomainRequests() {
	type testCase struct {
		Name          string
		ExpectedError string
	}
	var testCases = []testCase{
		{
			Name:          "should delete expired domain",
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.DeleteExpiredDomainRequests(s.ctx)
			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
		})
	}
}

func TestDomainRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(DomainRepositoryTestSuite))
}
