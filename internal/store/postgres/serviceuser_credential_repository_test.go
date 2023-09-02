package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/raystack/frontier/core/serviceuser"
	"github.com/raystack/frontier/internal/store/postgres"
	"github.com/raystack/frontier/pkg/db"
	"github.com/raystack/frontier/pkg/metadata"
	"github.com/raystack/salt/log"
	"github.com/stretchr/testify/suite"
)

type ServiceUserCredentialRepositoryTestSuite struct {
	suite.Suite
	ctx             context.Context
	client          *db.Client
	pool            *dockertest.Pool
	resource        *dockertest.Resource
	repository      *postgres.ServiceUserCredentialRepository
	serviceUserCred []serviceuser.Credential
	serviceUser     []serviceuser.ServiceUser
}

func (s *ServiceUserCredentialRepositoryTestSuite) SetupSuite() {
	var err error
	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewServiceUserCredentialRepository(s.client)
	s.serviceUser, err = bootstrapServiceUser(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceUserCredentialRepositoryTestSuite) SetupTest() {
	var err error
	s.serviceUserCred, err = bootstrapServiceUserCredential(s.client, s.serviceUser)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceUserCredentialRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceUserCredentialRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceUserCredentialRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_SERVICEUSERCREDENTIALS),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *ServiceUserCredentialRepositoryTestSuite) TestList() {
	type TestCase struct {
		Name           string
		ExpectedFilter serviceuser.Filter
		ExpectedResp   []serviceuser.Credential
		ExpectedErr    string
	}
	var testCases = []TestCase{
		{
			Name: "should get list of service user credential",
			ExpectedFilter: serviceuser.Filter{
				ServiceUserID: s.serviceUser[0].ID,
			},
			ExpectedResp: []serviceuser.Credential{
				{
					ID:            s.serviceUserCred[0].ID,
					ServiceUserID: s.serviceUser[0].ID,
					SecretHash:    []byte{},
					Title:         "title01",
				},
			},
			ExpectedErr: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.List(s.ctx, tc.ExpectedFilter)

			if tc.ExpectedErr != "" && err.Error() != tc.ExpectedErr {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedErr)
			}

			if !cmp.Equal(resp, tc.ExpectedResp, cmpopts.IgnoreFields(serviceuser.Credential{}, "PrivateKey", "PublicKey", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResp)
			}
		})
	}
}

func (s *ServiceUserCredentialRepositoryTestSuite) TestCreate() {
	testSecretHash := []byte("myrandomsecrethash")
	type TestCase struct {
		Name          string
		ExpectedInput serviceuser.Credential
		ExpectedResp  serviceuser.Credential
		ExpectedErr   string
	}
	var testCases = []TestCase{
		{
			Name: "should create service user credential",
			ExpectedInput: serviceuser.Credential{
				ServiceUserID: s.serviceUser[2].ID,
				Title:         "title_09",
				Metadata: metadata.Metadata{
					"key": "value",
				},
			},
			ExpectedResp: serviceuser.Credential{
				ServiceUserID: s.serviceUser[2].ID,
				Title:         "title_09",
				SecretHash:    []byte{},
				Metadata: metadata.Metadata{
					"key": "value",
				},
			},
			ExpectedErr: "",
		},
		{
			Name: "should create service user credential using secret hash",
			ExpectedInput: serviceuser.Credential{
				ServiceUserID: s.serviceUser[1].ID,
				Title:         "title_10",
				SecretHash:    testSecretHash,
				Metadata: metadata.Metadata{
					"foo": "bar",
				},
			},
			ExpectedResp: serviceuser.Credential{
				ServiceUserID: s.serviceUser[1].ID,
				Title:         "title_10",
				SecretHash:    testSecretHash,
				Metadata: metadata.Metadata{
					"foo": "bar",
				},
			},
			ExpectedErr: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Create(s.ctx, tc.ExpectedInput)
			if tc.ExpectedErr != "" && err.Error() != tc.ExpectedErr {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedErr)
			}
			if !cmp.Equal(resp, tc.ExpectedResp, cmpopts.IgnoreFields(serviceuser.Credential{}, "ID", "PrivateKey", "PublicKey", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResp)
			}
		})
	}
}

func (s *ServiceUserCredentialRepositoryTestSuite) TestGet() {
	type TestCase struct {
		Name         string
		ExpectedId   string
		ExpectedResp serviceuser.Credential
		ExpectedErr  string
	}
	var testCases = []TestCase{
		{
			Name:       "should get service user credential by using ID",
			ExpectedId: s.serviceUserCred[0].ID,
			ExpectedResp: serviceuser.Credential{
				ID:            s.serviceUserCred[0].ID,
				ServiceUserID: s.serviceUser[0].ID,
				SecretHash:    s.serviceUserCred[0].SecretHash,
				Title:         "title01",
			},
			ExpectedErr: "",
		},
		{
			Name:         "should give error if ID is empty",
			ExpectedId:   "",
			ExpectedResp: serviceuser.Credential{},
			ExpectedErr:  serviceuser.ErrInvalidKeyID.Error(),
		},
		{
			Name:         "should give error if service user credential does not exist",
			ExpectedId:   uuid.New().String(),
			ExpectedResp: serviceuser.Credential{},
			ExpectedErr:  serviceuser.ErrCredNotExist.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Get(s.ctx, tc.ExpectedId)

			if tc.ExpectedErr != "" && err.Error() != tc.ExpectedErr {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedErr)
			}

			if !cmp.Equal(resp, tc.ExpectedResp, cmpopts.IgnoreFields(serviceuser.Credential{}, "PrivateKey", "PublicKey", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedResp)
			}
		})
	}
}

func (s *ServiceUserCredentialRepositoryTestSuite) TestDelete() {
	type TestCase struct {
		Name        string
		ExpectedId  string
		ExpectedErr error
	}
	var testCases = []TestCase{
		{
			Name:        "should Delete service user credential",
			ExpectedId:  s.serviceUserCred[0].ID,
			ExpectedErr: nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.repository.Delete(s.ctx, tc.ExpectedId)
			if tc.ExpectedErr != nil {
				if err == nil || err.Error() != tc.ExpectedErr.Error() {
					s.T().Fatalf("got error %v, expected was %s", err, tc.ExpectedErr)
				}
			}
		})
	}
}

func TestServiceUserCredentialRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceUserCredentialRepositoryTestSuite))
}
