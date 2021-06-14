package user

import (
	"context"
	"fmt"
	"sort"

	"github.com/osapers/mch-back/internal/constant"
	"github.com/osapers/mch-back/internal/types"
	"github.com/osapers/mch-back/pkg/fakedata"
)

type MyProject struct {
	*types.Project
	Role string `json:"role"`
}

func (s *Service) MyProjects(ctx context.Context) ([]*MyProject, error) {
	userID := constant.GetUserIDFromCtx(ctx)

	user, err := s.storage.getByID(ctx, constant.GetUserIDFromCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	appliedProjects, err := s.storage.getAppliedProjects(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user applied projects: %w", err)
	}

	var appliedIds []string
	for projectID := range appliedProjects {
		appliedIds = append(appliedIds, projectID)
	}

	myProjects, err := s.storage.getMyProjects(ctx, userID, appliedIds)
	if err != nil {
		return nil, fmt.Errorf("get my projects: %w", err)
	}

	var projectOwnerIds []string
	for _, p := range myProjects {
		if p.OwnerID != userID {
			projectOwnerIds = append(projectOwnerIds, p.OwnerID)
		}
	}

	projectOwners, err := s.storage.getOwnersInfo(ctx, projectOwnerIds)
	if err != nil {
		return nil, fmt.Errorf("get projects owners: %w", err)
	}

	var res []*MyProject

	for _, project := range myProjects {
		if project.OwnerID == userID {
			project.Owner = types.ProjectOwner{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Image:     user.Photo,
			}
		} else {
			project.Owner = types.ProjectOwner{
				FirstName: projectOwners[project.OwnerID].FirstName,
				LastName:  projectOwners[project.OwnerID].LastName,
				Image:     projectOwners[project.OwnerID].Photo,
			}
		}

		res = append(res, &MyProject{Project: project.ToJson(), Role: project.GetRole(userID, appliedProjects)})
	}

	return res, nil
}

func (s *Service) SearchProjects(ctx context.Context) ([]*types.Project, error) {
	user, err := s.storage.getByID(ctx, constant.GetUserIDFromCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	viewedProjects, err := s.storage.getViewedProjects(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("list user viewed projects: %w", err)
	}

	projects, err := s.storage.searchProjects(ctx, user.ID, viewedProjects)
	if err != nil {
		return nil, fmt.Errorf("search projects: %w", err)
	}

	if len(projects) == 0 {
		return []*types.Project{}, nil
	}

	var projectOwnerIds []string
	for _, p := range projects {
		projectOwnerIds = append(projectOwnerIds, p.OwnerID)
	}

	projectOwners, err := s.storage.getOwnersInfo(ctx, projectOwnerIds)
	if err != nil {
		return nil, fmt.Errorf("get projects owners: %w", err)
	}

	for i, p := range projects {
		p.TagsIntersection = types.GetTagsIntersection(p.Tags, user.Tags)
		p.Owner = types.ProjectOwner{
			FirstName: projectOwners[p.OwnerID].FirstName,
			LastName:  projectOwners[p.OwnerID].LastName,
			Image:     projectOwners[p.OwnerID].Photo,
		}
		projects[i] = p.ToJson()
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return len(projects[i].TagsIntersection) > len(projects[j].TagsIntersection)
	})

	return projects, nil
}

func (s *Service) GenerateProject(ctx context.Context) error {
	ownerID := constant.GetUserIDFromCtx(ctx)

	project := fakedata.GenerateProject(ownerID)

	keywords, err := s.kwe.Extract(ctx, project.Description)
	if err != nil {
		return fmt.Errorf("extract keywords: %w", err)
	}

	tags, err := s.storage.matchTags(ctx, keywords)
	if err != nil {
		return fmt.Errorf("match tags: %w", err)
	}

	project.Tags = tags

	err = s.storage.createProject(ctx, project)
	if err != nil {
		return fmt.Errorf("generate project: %w", err)
	}

	return nil
}

func (s *Service) ViewProject(ctx context.Context, projectID string) error {
	err := s.storage.viewProject(ctx, constant.GetUserIDFromCtx(ctx), projectID)
	if err != nil {
		return fmt.Errorf("user view project: %w", err)
	}

	return nil
}

func (s *Service) ApplyToProject(ctx context.Context, projectID string) error {
	err := s.storage.applyToProject(ctx, constant.GetUserIDFromCtx(ctx), projectID)
	if err != nil {
		return fmt.Errorf("user apply to project: %w", err)
	}

	return nil
}

func (s *Service) DeleteProject(ctx context.Context, projectID string) error {
	err := s.storage.deleteProject(ctx, constant.GetUserIDFromCtx(ctx), projectID)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}

	return nil
}

func (s *Service) GetProjectCandidates(ctx context.Context, projectID string) ([]*types.User, error) {
	ownerID := constant.GetUserIDFromCtx(ctx)

	tags, err := s.storage.getProjectTags(ctx, ownerID, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project tags: %w", err)
	}

	filteredOutUsers, err := s.storage.getProjectViewers(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project viewers: %w", err)
	}

	users, err := s.storage.getProjectCandidates(ctx, tags, filteredOutUsers)
	if err != nil {
		return nil, fmt.Errorf("get project candidates: %w", err)
	}

	for _, u := range users {
		u.TagsIntersection = types.GetTagsIntersection(u.Tags, tags)
	}

	sort.SliceStable(users, func(i, j int) bool {
		return len(users[i].TagsIntersection) > len(users[j].TagsIntersection)
	})

	return users, nil
}
