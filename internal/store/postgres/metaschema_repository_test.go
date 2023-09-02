package postgres_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/raystack/frontier/core/metaschema"
	"github.com/raystack/frontier/internal/store/postgres"
	"github.com/raystack/frontier/pkg/db"
	"github.com/raystack/salt/log"
	"github.com/stretchr/testify/suite"
)

type MetaSchemaRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.MetaSchemaRepository
	metaschema []metaschema.MetaSchema
}

func (s *MetaSchemaRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewMetaSchemaRepository(logger, s.client)
}

func (s *MetaSchemaRepositoryTestSuite) SetupTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
	var err error
	s.metaschema, err = bootstrapMetaSchema(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *MetaSchemaRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *MetaSchemaRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *MetaSchemaRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_METASCHEMA),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *MetaSchemaRepositoryTestSuite) TestGet() {
	type testCase struct {
		Name               string
		SelectedID         string
		ExpectedMetaSchema metaschema.MetaSchema
		ExpectedError      string
	}
	var testCases = []testCase{
		{
			Name:       "should get Meta Schema",
			SelectedID: s.metaschema[0].ID,
			ExpectedMetaSchema: metaschema.MetaSchema{
				ID:     s.metaschema[0].ID,
				Name:   "metaSchema1",
				Schema: "schema1",
			},
			ExpectedError: "",
		},
		{
			Name:               "should give error if Id is empty",
			SelectedID:         "",
			ExpectedMetaSchema: metaschema.MetaSchema{},
			ExpectedError:      metaschema.ErrInvalidID.Error(),
		},
		{
			Name:               "should give error if meta schema does not exist",
			SelectedID:         uuid.New().String(),
			ExpectedMetaSchema: metaschema.MetaSchema{},
			ExpectedError:      metaschema.ErrInvalidID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Get(s.ctx, tc.ExpectedMetaSchema.ID)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedMetaSchema, cmpopts.IgnoreFields(metaschema.MetaSchema{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedMetaSchema)
			}
		})
	}
}

func (s *MetaSchemaRepositoryTestSuite) TestCreate() {
	uid := uuid.New().String()
	type testCase struct {
		Name               string
		InputMetaSchema    metaschema.MetaSchema
		ExpectedMetaSchema metaschema.MetaSchema
		ExpectedError      string
	}
	var testCases = []testCase{
		{
			Name: "should Create meta schema",
			InputMetaSchema: metaschema.MetaSchema{
				ID:     uid,
				Name:   "meta4",
				Schema: "schema4",
			},
			ExpectedMetaSchema: metaschema.MetaSchema{
				ID:     uid,
				Name:   "meta4",
				Schema: "schema4",
			},
			ExpectedError: "",
		},
		{
			Name: "should Give error if name is empty",
			InputMetaSchema: metaschema.MetaSchema{
				ID:     uid,
				Name:   "",
				Schema: "schema5",
			},
			ExpectedMetaSchema: metaschema.MetaSchema{},
			ExpectedError:      metaschema.ErrInvalidID.Error(),
		},
		{
			Name: "should Give error if schema is empty",
			InputMetaSchema: metaschema.MetaSchema{
				ID:     uid,
				Name:   "meta5",
				Schema: "",
			},
			ExpectedMetaSchema: metaschema.MetaSchema{},
			ExpectedError:      metaschema.ErrInvalidDetail.Error(),
		},
		{
			Name: "should Give error if existing meta schema id is used to creat new meta schema",
			InputMetaSchema: metaschema.MetaSchema{
				ID:     s.metaschema[0].ID,
				Name:   "meta4",
				Schema: "schema4",
			},
			ExpectedMetaSchema: metaschema.MetaSchema{},
			ExpectedError:      metaschema.ErrConflict.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Create(s.ctx, tc.InputMetaSchema)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedMetaSchema, cmpopts.IgnoreFields(metaschema.MetaSchema{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedMetaSchema)
			}
		})
	}
}

func (s *MetaSchemaRepositoryTestSuite) TestList() {
	type testCase struct {
		Name               string
		ExpectedMetaSchema []metaschema.MetaSchema
		ExpectedError      string
	}
	var testCases = []testCase{
		{
			Name: "should list all metaschema",
			ExpectedMetaSchema: []metaschema.MetaSchema{
				{
					Name:   "metaSchema1",
					Schema: "schema1",
				},
				{
					Name:   "metaSchema2",
					Schema: "schema2",
				},
			},
			ExpectedError: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.List(s.ctx)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedMetaSchema, cmpopts.IgnoreFields(metaschema.MetaSchema{}, "ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedMetaSchema)
			}
		})
	}
}

func (s *MetaSchemaRepositoryTestSuite) TestDelete() {
	type testCase struct {
		Name               string
		InputMetaSchema    string
		ExpectedMetaSchema string
		ExpectedError      string
	}
	var testCases = []testCase{
		{
			Name:               "should delete meta schema",
			InputMetaSchema:    s.metaschema[0].ID,
			ExpectedMetaSchema: s.metaschema[0].Name,
			ExpectedError:      "",
		},
		{
			Name:               "should give error if meta schema does not exist",
			InputMetaSchema:    uuid.New().String(),
			ExpectedMetaSchema: "",
			ExpectedError:      metaschema.ErrNotExist.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Delete(s.ctx, tc.InputMetaSchema)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedMetaSchema, cmpopts.IgnoreFields(metaschema.MetaSchema{}, "ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedMetaSchema)
			}
		})
	}
}

func (s *MetaSchemaRepositoryTestSuite) TestUpdate() {
	type testCase struct {
		Name               string
		InputMetaSchemaID  string
		InputMetaSchema    metaschema.MetaSchema
		ExpectedMetaSchema metaschema.MetaSchema
		ExpectedError      string
	}
	var testCases = []testCase{
		{
			Name:              "should update meta schema",
			InputMetaSchemaID: s.metaschema[0].ID,
			InputMetaSchema: metaschema.MetaSchema{
				Schema: "schema999",
			},
			ExpectedMetaSchema: metaschema.MetaSchema{
				Name:   "metaSchema1",
				Schema: "schema999",
			},
			ExpectedError: "",
		},
		{
			Name:              "should give error if metaschema with input id does not exist",
			InputMetaSchemaID: uuid.New().String(),
			InputMetaSchema: metaschema.MetaSchema{
				Schema: "schema999",
			},
			ExpectedMetaSchema: metaschema.MetaSchema{},
			ExpectedError:      metaschema.ErrNotExist.Error(),
		},
	}
	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.repository.Update(s.ctx, tc.InputMetaSchemaID, tc.InputMetaSchema)
			if tc.ExpectedError != "" {
				if err.Error() != tc.ExpectedError {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ExpectedError)
				}
			}
			if !cmp.Equal(resp, tc.ExpectedMetaSchema, cmpopts.IgnoreFields(metaschema.MetaSchema{}, "ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", resp, tc.ExpectedMetaSchema)
			}
		})
	}
}

func TestMetaSchemaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MetaSchemaRepositoryTestSuite))
}
