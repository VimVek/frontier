package user

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/raystack/frontier/core/permission"

	"github.com/raystack/frontier/pkg/utils"

	"github.com/raystack/frontier/core/relation"
	"github.com/raystack/frontier/internal/bootstrap/schema"
	"github.com/raystack/frontier/pkg/errors"
	"github.com/raystack/frontier/pkg/str"
)

type RelationService interface {
	Create(ctx context.Context, rel relation.Relation) (relation.Relation, error)
	CheckPermission(ctx context.Context, rel relation.Relation) (bool, error)
	Delete(ctx context.Context, rel relation.Relation) error
	LookupSubjects(ctx context.Context, rel relation.Relation) ([]string, error)
	LookupResources(ctx context.Context, rel relation.Relation) ([]string, error)
	BatchCheckPermission(ctx context.Context, relations []relation.Relation) ([]relation.CheckPair, error)
}

type PermissionService interface {
	List(ctx context.Context, flt permission.Filter) ([]permission.Permission, error)
}

type Service struct {
	repository        Repository
	relationService   RelationService
	Now               func() time.Time
	permissionService PermissionService
}

func NewService(repository Repository, relationRepo RelationService,
	permissionService PermissionService) *Service {
	return &Service{
		repository:      repository,
		relationService: relationRepo,
		Now: func() time.Time {
			return time.Now().UTC()
		},
		permissionService: permissionService,
	}
}

// GetByID email or slug
func (s Service) GetByID(ctx context.Context, id string) (User, error) {
	if isValidEmail(id) {
		return s.repository.GetByEmail(ctx, id)
	}
	if utils.IsValidUUID(id) {
		return s.repository.GetByID(ctx, id)
	}
	return s.repository.GetByName(ctx, strings.ToLower(id))
}

func (s Service) GetByIDs(ctx context.Context, userIDs []string) ([]User, error) {
	return s.repository.GetByIDs(ctx, userIDs)
}

func (s Service) GetByEmail(ctx context.Context, email string) (User, error) {
	email = strings.ToLower(email)
	return s.repository.GetByEmail(ctx, email)
}

func (s Service) Create(ctx context.Context, user User) (User, error) {
	return s.repository.Create(ctx, User{
		Name:     strings.ToLower(user.Name),
		Email:    strings.ToLower(user.Email),
		Avatar:   user.Avatar,
		Title:    user.Title,
		Metadata: user.Metadata,
	})
}

func (s Service) List(ctx context.Context, flt Filter) ([]User, error) {
	if flt.OrgID != "" {
		return s.ListByOrg(ctx, flt.OrgID, schema.MembershipPermission)
	}
	if flt.GroupID != "" {
		return s.ListByGroup(ctx, flt.GroupID, schema.MembershipPermission)
	}

	// state gets filtered in db
	return s.repository.List(ctx, flt)
}

// Update by user uuid, email or slug
func (s Service) Update(ctx context.Context, toUpdate User) (User, error) {
	id := toUpdate.ID
	toUpdate.Email = strings.ToLower(toUpdate.Email)
	toUpdate.Name = strings.ToLower(toUpdate.Name)
	if isValidEmail(id) {
		return s.UpdateByEmail(ctx, toUpdate)
	}
	if utils.IsValidUUID(id) {
		return s.repository.UpdateByID(ctx, toUpdate)
	}
	return s.repository.UpdateByName(ctx, toUpdate)
}

func (s Service) UpdateByEmail(ctx context.Context, toUpdate User) (User, error) {
	toUpdate.Email = strings.ToLower(toUpdate.Email)
	toUpdate.Name = strings.ToLower(toUpdate.Name)
	return s.repository.UpdateByEmail(ctx, toUpdate)
}

func (s Service) Enable(ctx context.Context, id string) error {
	return s.repository.SetState(ctx, id, Enabled)
}

func (s Service) Disable(ctx context.Context, id string) error {
	return s.repository.SetState(ctx, id, Disabled)
}

func (s Service) Delete(ctx context.Context, id string) error {
	if err := s.relationService.Delete(ctx, relation.Relation{Object: relation.Object{
		ID:        id,
		Namespace: schema.ProjectNamespace,
	}}); err != nil {
		return err
	}
	return s.repository.Delete(ctx, id)
}

func (s Service) ListByOrg(ctx context.Context, orgID string, permissionFilter string) ([]User, error) {
	userIDs, err := s.relationService.LookupSubjects(ctx, relation.Relation{
		Object: relation.Object{
			ID:        orgID,
			Namespace: schema.OrganizationNamespace,
		},
		Subject: relation.Subject{
			Namespace: schema.UserPrincipal,
		},
		RelationName: permissionFilter,
	})
	if err != nil {
		return nil, err
	}
	if len(userIDs) == 0 {
		// no users
		return []User{}, nil
	}
	return s.repository.GetByIDs(ctx, userIDs)
}

func (s Service) ListByGroup(ctx context.Context, groupID string, permissionFilter string) ([]User, error) {
	userIDs, err := s.relationService.LookupSubjects(ctx, relation.Relation{
		Object: relation.Object{
			ID:        groupID,
			Namespace: schema.GroupPrincipal,
		},
		Subject: relation.Subject{
			Namespace: schema.UserPrincipal,
		},
		RelationName: permissionFilter,
	})
	if err != nil {
		return nil, err
	}
	if len(userIDs) == 0 {
		// no users
		return []User{}, nil
	}

	return s.repository.GetByIDs(ctx, userIDs)
}

func (s Service) ListByGroupWithAccessPairs(ctx context.Context, groupID string, permissions []string) ([]AccessPair, error) {
	users, err := s.ListByGroup(ctx, groupID, schema.MembershipPermission)
	if err != nil {
		return nil, fmt.Errorf("fetching users: %w", err)
	}

	// ensure all permissions exist
	permModels, err := s.permissionService.List(ctx, permission.Filter{Namespace: schema.GroupNamespace})
	if err != nil {
		return nil, fmt.Errorf("fetching permission: %w", err)
	}
	for _, perm := range permissions {
		if !utils.ContainsFunc(permModels, func(p permission.Permission) bool {
			return p.Name == perm
		}) {
			return nil, fmt.Errorf("invalid %s: %w", perm, permission.ErrNotExist)
		}
	}

	var accessPairs []AccessPair
	var checks []relation.Relation
	// fetch access pairs permissions
	for _, user := range users {
		relSubject := relation.Subject{
			ID:        user.ID,
			Namespace: schema.UserPrincipal,
		}
		relObj := relation.Object{
			ID:        groupID,
			Namespace: schema.GroupNamespace,
		}
		for _, permission := range permissions {
			checks = append(checks, relation.Relation{
				Subject:      relSubject,
				Object:       relObj,
				RelationName: permission,
			})
		}

		accessPairs = append(accessPairs, AccessPair{
			User: user,
			Can:  []string{},
			On:   schema.JoinNamespaceAndResourceID(schema.GroupNamespace, groupID),
		})
	}

	// create permission check requests
	checkPairs, err := s.relationService.BatchCheckPermission(ctx, checks)
	if err != nil {
		return nil, fmt.Errorf("checking permissions: %w", err)
	}
	successChecks := utils.Filter(checkPairs, func(pair relation.CheckPair) bool {
		return pair.Status
	})

	for _, checkPair := range successChecks {
		accessPairs = utils.Map(accessPairs, func(pair AccessPair) AccessPair {
			if pair.User.ID == checkPair.Relation.Subject.ID {
				pair.Can = append(pair.Can, checkPair.Relation.RelationName)
			}
			return pair
		})
	}
	return accessPairs, nil
}

func (s Service) Sudo(ctx context.Context, id string) error {
	currentUser, err := s.GetByID(ctx, id)
	if errors.Is(err, ErrNotExist) {
		if isValidEmail(id) {
			// create a new user
			currentUser, err = s.Create(ctx, User{
				Email: id,
				Name:  str.GenerateUserSlug(id),
			})
			if err != nil {
				return err
			}
		} else {
			// skip
			return nil
		}
	}
	if err != nil {
		return err
	}

	// check if already su
	if ok, err := s.IsSudo(ctx, currentUser.ID); err != nil {
		return err
	} else if ok {
		return nil
	}

	// mark su
	_, err = s.relationService.Create(ctx, relation.Relation{
		Object: relation.Object{
			ID:        schema.PlatformID,
			Namespace: schema.PlatformNamespace,
		},
		Subject: relation.Subject{
			ID:        currentUser.ID,
			Namespace: schema.UserPrincipal,
		},
		RelationName: schema.AdminRelationName,
	})
	return err
}

func (s Service) IsSudo(ctx context.Context, id string) (bool, error) {
	status, err := s.relationService.CheckPermission(ctx, relation.Relation{
		Subject: relation.Subject{
			ID:        id,
			Namespace: schema.UserPrincipal,
		},
		Object: relation.Object{
			ID:        schema.PlatformID,
			Namespace: schema.PlatformNamespace,
		},
		RelationName: schema.SudoPermission,
	})
	if err != nil {
		return false, err
	}
	return status, nil
}

func isValidEmail(str string) bool {
	_, err := mail.ParseAddress(str)
	return err == nil
}
