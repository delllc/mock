package domain

import (
	"errors"
	"testing"
)

func TestNewProject(t *testing.T) {
	project, err := NewProject("Shop API", "shop")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if project.Name != "Shop API" {
		t.Errorf("unexpected name: %s", project.Name)
	}

	if project.Slug != "shop" {
		t.Errorf("unexpected slug: %s", project.Slug)
	}

	if project.BasePath != "/api" {
		t.Errorf("unexpected base path: %s", project.BasePath)
	}

	if project.ID == [16]byte{} {
		t.Error("project ID must not be empty")
	}
}

func TestNewProjectInvalidSlug(t *testing.T) {
	_, err := NewProject("Shop API", "Shop API")

	if !errors.Is(err, ErrInvalidProjectSlug) {
		t.Fatalf("expected ErrInvalidProjectSlug, got %v", err)
	}
}

func TestNewProjectEmptyName(t *testing.T) {
	_, err := NewProject("", "shop")

	if !errors.Is(err, ErrProjectNameRequired) {
		t.Fatalf("expected ErrProjectNameRequired, got %v", err)
	}
}
