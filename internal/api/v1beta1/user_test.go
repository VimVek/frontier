package v1beta1

import (
	"context"
	"testing"
	"time"

	"github.com/odpf/shield/core/authenticate/token"

	"github.com/odpf/shield/pkg/utils"

	"github.com/odpf/shield/pkg/errors"

	"github.com/odpf/shield/core/organization"

	"github.com/odpf/shield/core/group"
	"github.com/odpf/shield/core/user"
	"github.com/odpf/shield/internal/api/v1beta1/mocks"
	"github.com/odpf/shield/pkg/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	shieldv1beta1 "github.com/odpf/shield/proto/v1beta1"
)

var (
	testUserID  = "9f256f86-31a3-11ec-8d3d-0242ac130003"
	testUserMap = map[string]user.User{
		"9f256f86-31a3-11ec-8d3d-0242ac130003": {
			ID:    "9f256f86-31a3-11ec-8d3d-0242ac130003",
			Title: "User 1",
			Name:  "user1",
			Email: "test@test.com",
			Metadata: metadata.Metadata{
				"foo":    "bar",
				"age":    21,
				"intern": true,
			},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}
)

func TestListUsers(t *testing.T) {
	table := []struct {
		title string
		setup func(us *mocks.UserService)
		req   *shieldv1beta1.ListUsersRequest
		want  *shieldv1beta1.ListUsersResponse
		err   error
	}{
		{
			title: "should return internal error in if user service return some error",
			setup: func(us *mocks.UserService) {
				us.EXPECT().List(mock.Anything, mock.Anything).Return([]user.User{}, errors.New("some error"))
			},
			req: &shieldv1beta1.ListUsersRequest{
				PageSize: 50,
				PageNum:  1,
				Keyword:  "",
			},
			want: nil,
			err:  status.Errorf(codes.Internal, ErrInternalServer.Error()),
		}, {
			title: "should return all users if user service return all users",
			setup: func(us *mocks.UserService) {
				var testUserList []user.User
				for _, u := range testUserMap {
					testUserList = append(testUserList, u)
				}
				us.EXPECT().List(mock.Anything, mock.Anything).Return(testUserList, nil)
			},
			req: &shieldv1beta1.ListUsersRequest{
				PageSize: 50,
				PageNum:  1,
				Keyword:  "",
			},
			want: &shieldv1beta1.ListUsersResponse{
				Count: 1,
				Users: []*shieldv1beta1.User{
					{
						Id:    "9f256f86-31a3-11ec-8d3d-0242ac130003",
						Title: "User 1",
						Name:  "user1",
						Email: "test@test.com",
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"foo":    structpb.NewStringValue("bar"),
								"age":    structpb.NewNumberValue(21),
								"intern": structpb.NewBoolValue(true),
							},
						},
						CreatedAt: timestamppb.New(time.Time{}),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			if tt.setup != nil {
				tt.setup(mockUserSrv)
			}
			mockDep := Handler{userService: mockUserSrv}
			req := tt.req
			resp, err := mockDep.ListUsers(context.Background(), req)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestCreateUser(t *testing.T) {
	email := "user@odpf.io"
	table := []struct {
		title string
		setup func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context
		req   *shieldv1beta1.CreateUserRequest
		want  *shieldv1beta1.CreateUserResponse
		err   error
	}{
		{
			title: "should return unauthenticated error if no auth email header in context",
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title:    "some user",
				Email:    "abc@test.com",
				Metadata: &structpb.Struct{},
			}},
			want: nil,
			err:  grpcUnauthenticated,
		},
		{
			title: "should return bad request error if metadata is not parsable",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "some user",
				Email: "abc@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewNullValue(),
					},
				},
			}},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return bad request error if email is empty",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), user.User{
					Title: "some user",
				}).Return(user.User{}, user.ErrInvalidEmail)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "some user",
				Email: "",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewNullValue(),
					},
				},
			}},
			want: nil,
			err:  grpcBadBodyError,
		},

		{
			title: "should return already exist error if user service return error conflict",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), user.User{
					Title:    "some user",
					Email:    "abc@test.com",
					Name:     "user-slug",
					Metadata: metadata.Metadata{},
				}).Return(user.User{}, user.ErrConflict)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title:    "some user",
				Email:    "abc@test.com",
				Name:     "user-slug",
				Metadata: &structpb.Struct{},
			}},
			want: nil,
			err:  grpcConflictError,
		},
		{
			title: "should return success if user email contain whitespace but still valid service return nil error",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), user.User{
					Title:    "some user",
					Email:    "abc@test.com",
					Name:     "user-slug",
					Metadata: metadata.Metadata{"foo": "bar"},
				}).Return(
					user.User{
						ID:       "new-abc",
						Title:    "some user",
						Email:    "abc@test.com",
						Name:     "user-slug",
						Metadata: metadata.Metadata{"foo": "bar"},
					}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "some user",
				Email: "  abc@test.com  ",
				Name:  "user-slug",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: &shieldv1beta1.CreateUserResponse{User: &shieldv1beta1.User{
				Id:    "new-abc",
				Title: "some user",
				Email: "abc@test.com",
				Name:  "user-slug",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
		{
			title: "should return success if user service return nil error",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), user.User{
					Title:    "some user",
					Email:    "abc@test.com",
					Name:     "user-slug",
					Metadata: metadata.Metadata{"foo": "bar"},
				}).Return(
					user.User{
						ID:       "new-abc",
						Title:    "some user",
						Email:    "abc@test.com",
						Name:     "user-slug",
						Metadata: metadata.Metadata{"foo": "bar"},
					}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "some user",
				Email: "abc@test.com",
				Name:  "user-slug",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: &shieldv1beta1.CreateUserResponse{User: &shieldv1beta1.User{
				Id:    "new-abc",
				Title: "some user",
				Email: "abc@test.com",
				Name:  "user-slug",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			var resp *shieldv1beta1.CreateUserResponse
			var err error

			ctx := context.Background()
			mockUserSrv := new(mocks.UserService)
			mockMetaSrv := new(mocks.MetaSchemaService)
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockUserSrv, mockMetaSrv)
			}
			mockDep := Handler{userService: mockUserSrv, metaSchemaService: mockMetaSrv}
			resp, err = mockDep.CreateUser(ctx, tt.req)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	randomID := utils.NewString()
	table := []struct {
		title string
		req   *shieldv1beta1.GetUserRequest
		setup func(us *mocks.UserService)
		want  *shieldv1beta1.GetUserResponse
		err   error
	}{
		{
			title: "should return not found error if user does not exist",
			setup: func(us *mocks.UserService) {
				us.EXPECT().GetByID(mock.AnythingOfType("*context.emptyCtx"), randomID).Return(user.User{}, user.ErrNotExist)
			},
			req: &shieldv1beta1.GetUserRequest{
				Id: randomID,
			},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return not found error if user id is not uuid",
			setup: func(us *mocks.UserService) {
				us.EXPECT().GetByID(mock.AnythingOfType("*context.emptyCtx"), "some-id").Return(user.User{}, user.ErrInvalidUUID)
			},
			req: &shieldv1beta1.GetUserRequest{
				Id: "some-id",
			},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return not found error if user id is invalid",
			setup: func(us *mocks.UserService) {
				us.EXPECT().GetByID(mock.AnythingOfType("*context.emptyCtx"), "").Return(user.User{}, user.ErrInvalidID)
			},
			req:  &shieldv1beta1.GetUserRequest{},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return user if user service return nil error",
			setup: func(us *mocks.UserService) {
				us.EXPECT().GetByID(mock.AnythingOfType("*context.emptyCtx"), randomID).Return(
					user.User{
						ID:    randomID,
						Title: "some user",
						Email: "someuser@test.com",
						Metadata: metadata.Metadata{
							"foo": "bar",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)
			},
			req: &shieldv1beta1.GetUserRequest{
				Id: randomID,
			},
			want: &shieldv1beta1.GetUserResponse{User: &shieldv1beta1.User{
				Id:    randomID,
				Title: "some user",
				Email: "someuser@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			if tt.setup != nil {
				tt.setup(mockUserSrv)
			}
			mockDep := Handler{userService: mockUserSrv}
			resp, err := mockDep.GetUser(context.Background(), tt.req)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	email := "user@odpf.io"
	table := []struct {
		title  string
		setup  func(ctx context.Context, us *mocks.UserService, ss *mocks.SessionService) context.Context
		header string
		want   *shieldv1beta1.GetCurrentUserResponse
		err    error
	}{
		{
			title: "should return unauthenticated error if no auth email header in context",
			want:  nil,
			err:   grpcUnauthenticated,
			setup: func(ctx context.Context, us *mocks.UserService, ss *mocks.SessionService) context.Context {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("*context.emptyCtx")).Return(user.User{}, errors.ErrUnauthenticated)
				return ctx
			},
		},
		{
			title: "should return not found error if user does not exist",
			setup: func(ctx context.Context, us *mocks.UserService, ss *mocks.SessionService) context.Context {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("*context.valueCtx")).Return(user.User{}, user.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return error if user service return some error",
			setup: func(ctx context.Context, us *mocks.UserService, ss *mocks.SessionService) context.Context {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("*context.valueCtx")).Return(user.User{}, errors.New("some error"))
				return user.SetContextWithEmail(ctx, email)
			},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "should return user if user service return nil error",
			setup: func(ctx context.Context, us *mocks.UserService, ss *mocks.SessionService) context.Context {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("*context.valueCtx")).Return(
					user.User{
						ID:    "user-id-1",
						Title: "some user",
						Email: "someuser@test.com",
						Metadata: metadata.Metadata{
							"foo": "bar",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			want: &shieldv1beta1.GetCurrentUserResponse{User: &shieldv1beta1.User{
				Id:    "user-id-1",
				Title: "some user",
				Email: "someuser@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			ctx := context.Background()
			mockUserSrv := new(mocks.UserService)

			mockOrgService := new(mocks.OrganizationService)
			mockOrgService.EXPECT().ListByUser(mock.Anything, "user-id-1").Return([]organization.Organization{}, nil)
			registrationService := new(mocks.RegistrationService)
			registrationService.EXPECT().Token(mock.Anything,
				mock.AnythingOfType("user.User"), map[string]string{"orgs": ""}).Return(nil, token.ErrMissingRSADisableToken)

			if tt.setup != nil {
				ctx = tt.setup(ctx, mockUserSrv, nil)
			}
			mockDep := Handler{
				userService:         mockUserSrv,
				orgService:          mockOrgService,
				registrationService: registrationService,
			}
			resp, err := mockDep.GetCurrentUser(ctx, nil)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	someID := utils.NewString()
	table := []struct {
		title  string
		setup  func(us *mocks.UserService, ms *mocks.MetaSchemaService)
		req    *shieldv1beta1.UpdateUserRequest
		header string
		want   *shieldv1beta1.UpdateUserResponse
		err    error
	}{
		{
			title: "should return internal error if user service return some error",
			setup: func(us *mocks.UserService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), user.User{
					ID:    someID,
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, errors.New("some error"))
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				}},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "should return not found error if id is invalid",
			setup: func(us *mocks.UserService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), user.User{
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, user.ErrInvalidID)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Body: &shieldv1beta1.UserRequestBody{
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				}},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return already exist error if user service return error conflict",
			setup: func(us *mocks.UserService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), user.User{
					ID:    someID,
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, user.ErrConflict)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				}},
			want: nil,
			err:  grpcConflictError,
		},
		{
			title: "should return bad request error if email in request empty",
			setup: func(us *mocks.UserService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), user.User{
					ID:    someID,
					Title: "abc user",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, user.ErrInvalidEmail)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Title: "abc user",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return bad request error if empty request body",
			req:   &shieldv1beta1.UpdateUserRequest{Id: someID, Body: nil},
			want:  nil,
			err:   grpcBadBodyError,
		},
		{
			title: "should return success if user service return nil error",
			setup: func(us *mocks.UserService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), user.User{
					ID:    someID,
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(
					user.User{
						ID:    someID,
						Title: "abc user",
						Email: "user@odpf.io",
						Metadata: metadata.Metadata{
							"foo": "bar",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				}},
			want: &shieldv1beta1.UpdateUserResponse{User: &shieldv1beta1.User{
				Id:    someID,
				Title: "abc user",
				Email: "user@odpf.io",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
		{
			title: "should return success even though name is empty",
			setup: func(us *mocks.UserService, ms *mocks.MetaSchemaService) {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), user.User{
					ID:    someID,
					Email: "user@odpf.io",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(
					user.User{
						ID:    someID,
						Email: "user@odpf.io",
						Metadata: metadata.Metadata{
							"foo": "bar",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Email: "user@odpf.io",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				}},
			want: &shieldv1beta1.UpdateUserResponse{User: &shieldv1beta1.User{
				Id:    someID,
				Email: "user@odpf.io",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			mockMetaSrv := new(mocks.MetaSchemaService)
			ctx := context.Background()
			if tt.setup != nil {
				tt.setup(mockUserSrv, mockMetaSrv)
			}
			mockDep := Handler{userService: mockUserSrv, metaSchemaService: mockMetaSrv}
			resp, err := mockDep.UpdateUser(ctx, tt.req)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

func TestUpdateCurrentUser(t *testing.T) {
	email := "user@odpf.io"
	table := []struct {
		title  string
		setup  func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context
		req    *shieldv1beta1.UpdateCurrentUserRequest
		header string
		want   *shieldv1beta1.UpdateCurrentUserResponse
		err    error
	}{
		{
			title: "should return unauthenticated error if auth email header not exist",
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "abc user",
				Email: "abcuser123@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: nil,
			err:  grpcUnauthenticated,
		},
		{
			title: "should return internal error if user service return some error",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().UpdateByEmail(mock.AnythingOfType("*context.valueCtx"), user.User{
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, errors.New("some error"))
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "abc user",
				Email: "user@odpf.io",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "should return not found error if user service return err not exist",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().UpdateByEmail(mock.AnythingOfType("*context.valueCtx"), user.User{
					Title: "abc user",
					Email: "user@odpf.io",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, user.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "abc user",
				Email: "user@odpf.io",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return bad request error if diff emails in header and body",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "abc user",
				Email: "abcuser123@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return bad request error if empty request body",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				return user.SetContextWithEmail(ctx, email)
			},
			req:  &shieldv1beta1.UpdateCurrentUserRequest{Body: nil},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return success if user service return nil error",
			setup: func(ctx context.Context, us *mocks.UserService, ms *mocks.MetaSchemaService) context.Context {
				ms.EXPECT().Validate(mock.AnythingOfType("metadata.Metadata"), userMetaSchema).Return(nil)
				us.EXPECT().UpdateByEmail(mock.Anything, mock.Anything).Return(
					user.User{
						ID:    "user-id-1",
						Title: "abc user",
						Email: "user@odpf.io",
						Metadata: metadata.Metadata{
							"foo": "bar",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Title: "abc user",
				Email: "user@odpf.io",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: &shieldv1beta1.UpdateCurrentUserResponse{User: &shieldv1beta1.User{
				Id:    "user-id-1",
				Title: "abc user",
				Email: "user@odpf.io",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			mockMetaSrv := new(mocks.MetaSchemaService)
			ctx := context.Background()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockUserSrv, mockMetaSrv)
			}
			mockDep := Handler{userService: mockUserSrv, metaSchemaService: mockMetaSrv}
			resp, err := mockDep.UpdateCurrentUser(ctx, tt.req)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestHandler_ListUserGroups(t *testing.T) {
	someUserID := utils.NewString()
	tests := []struct {
		name    string
		setup   func(gs *mocks.GroupService)
		request *shieldv1beta1.ListUserGroupsRequest
		want    *shieldv1beta1.ListUserGroupsResponse
		wantErr error
	}{
		{
			name: "should return internal error if group service return some error",
			setup: func(gs *mocks.GroupService) {
				gs.EXPECT().ListByUser(mock.AnythingOfType("*context.emptyCtx"), someUserID).Return([]group.Group{}, errors.New("some error"))
			},
			request: &shieldv1beta1.ListUserGroupsRequest{
				Id: someUserID,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return empty list if user does not exist",
			setup: func(gs *mocks.GroupService) {
				gs.EXPECT().ListByUser(mock.AnythingOfType("*context.emptyCtx"), someUserID).Return([]group.Group{}, nil)
			},
			request: &shieldv1beta1.ListUserGroupsRequest{
				Id: someUserID,
			},
			want:    &shieldv1beta1.ListUserGroupsResponse{},
			wantErr: nil,
		},
		// if user id empty, it would not go to this handler
		{
			name: "should return groups if group service return not nil",
			setup: func(gs *mocks.GroupService) {
				var testGroupList []group.Group
				for _, g := range testGroupMap {
					testGroupList = append(testGroupList, g)
				}
				gs.EXPECT().ListByUser(mock.AnythingOfType("*context.emptyCtx"), someUserID).Return(testGroupList, nil)
			},
			request: &shieldv1beta1.ListUserGroupsRequest{
				Id: someUserID,
			},
			want: &shieldv1beta1.ListUserGroupsResponse{
				Groups: []*shieldv1beta1.Group{
					{
						Id:    "9f256f86-31a3-11ec-8d3d-0242ac130003",
						Name:  "group-1",
						OrgId: "9f256f86-31a3-11ec-8d3d-0242ac130003",
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"foo": structpb.NewStringValue("bar"),
							},
						},
						CreatedAt: timestamppb.New(time.Time{}),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGroupSvc := new(mocks.GroupService)
			if tt.setup != nil {
				tt.setup(mockGroupSvc)
			}
			h := Handler{
				groupService: mockGroupSvc,
			}
			got, err := h.ListUserGroups(context.Background(), tt.request)
			assert.EqualValues(t, got, tt.want)
			assert.EqualValues(t, err, tt.wantErr)
		})
	}
}
