package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/raystack/frontier/core/namespace"
	"github.com/raystack/frontier/core/preference"
	"github.com/raystack/frontier/internal/store/postgres"
	"github.com/raystack/frontier/pkg/db"
	"github.com/raystack/salt/log"
	"github.com/stretchr/testify/suite"
)

type PreferenceRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.PreferenceRepository
	preference []preference.Preference
	namespace  []namespace.Namespace
}

func (s *PreferenceRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewPreferenceRepository(s.client)

	s.namespace, err = bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *PreferenceRepositoryTestSuite) SetupTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
	var err error
	s.preference, err = bootstrapPreference(s.client, s.namespace)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *PreferenceRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *PreferenceRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *PreferenceRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_PREFERENCES),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *PreferenceRepositoryTestSuite) TestSet() {
	uid := uuid.New()
	ruid := uuid.New()
	type testCase struct {
		Name               string
		SelectedPreference preference.Preference
		ExpectedPreference preference.Preference
		ExpectedError      string
	}
	var testCases = []testCase{
		{
			Name: "should Set Preferece",
			SelectedPreference: preference.Preference{
				ID:           uid.String(),
				Name:         "some_name",
				Value:        "some_value",
				ResourceID:   ruid.String(),
				ResourceType: s.namespace[0].Name,
			},
			ExpectedPreference: preference.Preference{
				ID:           uid.String(),
				Name:         "some_name",
				Value:        "some_value",
				ResourceID:   ruid.String(),
				ResourceType: s.namespace[0].Name,
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Set(s.ctx, tc.SelectedPreference)

			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
			if !cmp.Equal(resp, tc.ExpectedPreference, cmpopts.IgnoreFields(preference.Preference{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedPreference)
			}
		})
	}
}

func (s *PreferenceRepositoryTestSuite) TestGet() {
	uuidString := s.preference[0].ID
	uid, err := uuid.Parse(uuidString)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		return
	}
	type testCase struct {
		Name               string
		SelectedID         uuid.UUID
		ExpectedPreference preference.Preference
		ExpectedError      string
	}
	var testCases = []testCase{
		{
			Name:       "should get Preferece by id",
			SelectedID: uid,
			ExpectedPreference: preference.Preference{
				ID:           s.preference[0].ID,
				ResourceID:   s.preference[0].ResourceID,
				Name:         "some_name",
				Value:        "some_value",
				ResourceType: s.preference[0].ResourceType,
			},
			ExpectedError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Get(s.ctx, tc.SelectedID)

			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
			if !cmp.Equal(resp, tc.ExpectedPreference, cmpopts.IgnoreFields(preference.Preference{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedPreference)
			}
		})
	}
}

func (s *PreferenceRepositoryTestSuite) TestList() {
	type testCase struct {
		Name               string
		SelectedFilter     preference.Filter
		ExpectedPreference []preference.Preference
		ExpectedError      string
	}
	var testCases = []testCase{
		{
			Name: "should List Preferece",
			SelectedFilter: preference.Filter{
				OrgID: s.preference[0].ResourceID,
			},
			ExpectedPreference: []preference.Preference{
				{
					ID:           s.preference[0].ID,
					Name:         "some_name",
					Value:        "some_value",
					ResourceType: s.preference[0].ResourceType,
					ResourceID:   s.preference[0].ResourceID,
				},
			},
			ExpectedError: "",
		},
		{
			Name: "should give error is resource name is not given with resource id and vice-versa",
			SelectedFilter: preference.Filter{
				ResourceID: s.namespace[0].ID,
			},
			ExpectedPreference: nil,
			ExpectedError:      preference.ErrInvalidFilter.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.List(s.ctx, tc.SelectedFilter)
			if err != nil && err.Error() != tc.ExpectedError {
				s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
			}
			if !cmp.Equal(resp, tc.ExpectedPreference, cmpopts.IgnoreFields(preference.Preference{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedPreference)
			}
		})
	}
}

func TestPreferenceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PreferenceRepositoryTestSuite))
}
