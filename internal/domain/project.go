package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"uuid"
)

var (
	ErrProjectNameRequired = errors.New("project name is required")
	ErrProjectSlugRequired = errors.New("slug is required")
	ErrInvalidProjectSlug  = errors.New("invalig project slug")
	ErrInvalidBasePath     = errors.New("invalid project base path")
)

var projectSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Project struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	BasePath  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProject(name, slug string) (*Project, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(strings.ToLower(slug))

	if name == "" {
		return nil, ErrProjectNameRequired
	}

	if slug == "" {
		return nil, ErrProjectSlugRequired
	}

	if !projectSlugPattern.MatchString(slug) {
		return nil, fmt.Errorf(
			"%w: %q", ErrInvalidProjectSlug, slug,
		)
	}

	now := time.Now().UTC()

	return &Project{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		BasePath:  "/api",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil

}

func (p *Project) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrProjectNameRequired
	}

	if strings.TrimSpace(p.Slug) == "" {
		return ErrProjectSlugRequired
	}

	if !projectSlugPattern.MatchString(p.Slug) {
		return fmt.Errorf(
			"%w: %q", ErrInvalidProjectSlug, p.Slug,
		)
	}

	if !strings.HasPrefix(p.BasePath, "/") {
		return ErrInvalidBasePath
	}

	if p.BasePath != "/" {
		p.BasePath = strings.TrimRight(p.BasePath, "/")
	}

	return nil
}
