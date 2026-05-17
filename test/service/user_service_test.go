package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/model"
	"github.com/magistraapta/golang-devops/internal/service"
	"github.com/magistraapta/golang-devops/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// ---- fake repository ----

type fakeUserRepository struct {
	createdUser   *model.User
	users         []model.User
	createErr     error
	getByIDErr    error
	getByEmailErr error
	updateErr     error
	deleteErr     error
	getAllErr     error
}

func (r *fakeUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	if r.createErr != nil {
		return r.createErr
	}
	copied := *user
	r.createdUser = &copied
	return nil
}

func (r *fakeUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	for _, u := range r.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *fakeUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if r.getByEmailErr != nil {
		return nil, r.getByEmailErr
	}
	for _, u := range r.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, nil // not found is not an error for email lookup
}

func (r *fakeUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.updateErr
}

func (r *fakeUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.deleteErr
}

func (r *fakeUserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	if r.getAllErr != nil {
		return nil, r.getAllErr
	}
	return r.users, nil
}

// ---- helpers ----

func newUserService(repo *fakeUserRepository) service.UserService {
	return service.NewUserService(repo)
}

func assertCustomError(t *testing.T, err error, wantCode int) {
	t.Helper()
	ce, ok := err.(*utils.CustomError)
	if !ok {
		t.Fatalf("expected *utils.CustomError, got %T", err)
	}
	if ce.Code != wantCode {
		t.Fatalf("expected code %d, got %d", wantCode, ce.Code)
	}
}

// ---- CreateUser ----

func TestCreateUser_SetsIDAndTimestamps(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{}
	svc := newUserService(repo)

	user := &model.User{Username: "alice", Email: "alice@example.com", Password: "secret"}
	if err := svc.CreateUser(ctx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.createdUser.ID == uuid.Nil {
		t.Fatal("expected ID to be set")
	}
	if repo.createdUser.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
	if repo.createdUser.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}
}

func TestCreateUser_HashesPassword(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{}
	svc := newUserService(repo)

	const plain = "secret123"
	user := &model.User{Username: "alice", Email: "alice@example.com", Password: plain}
	if err := svc.CreateUser(ctx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.createdUser.Password == plain {
		t.Fatal("expected password to be hashed, got plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.createdUser.Password), []byte(plain)); err != nil {
		t.Fatalf("hash does not match original password: %v", err)
	}
}

func TestCreateUser_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{createErr: errors.New("db error")}
	svc := newUserService(repo)

	err := svc.CreateUser(ctx, &model.User{Username: "alice", Email: "alice@example.com", Password: "secret"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertCustomError(t, err, 500)
}

// ---- GetUserByID ----

func TestGetUserByID_Success(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	repo := &fakeUserRepository{users: []model.User{{ID: id, Username: "alice", Email: "alice@example.com"}}}
	svc := newUserService(repo)

	user, err := svc.GetUserByID(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != id {
		t.Fatalf("expected ID %v, got %v", id, user.ID)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{} // empty users slice
	svc := newUserService(repo)

	user, err := svc.GetUserByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
	assertCustomError(t, err, 500)
}

func TestGetUserByID_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{getByIDErr: errors.New("db error")}
	svc := newUserService(repo)

	user, err := svc.GetUserByID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
	assertCustomError(t, err, 500)
}

// ---- GetUserByEmail ----

func TestGetUserByEmail_Success(t *testing.T) {
	ctx := context.Background()
	email := "alice@example.com"
	repo := &fakeUserRepository{users: []model.User{{ID: uuid.New(), Username: "alice", Email: email}}}
	svc := newUserService(repo)

	user, err := svc.GetUserByEmail(ctx, email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != email {
		t.Fatalf("expected email %q, got %q", email, user.Email)
	}
}

func TestGetUserByEmail_NotFound_ReturnsNil(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{} // empty
	svc := newUserService(repo)

	user, err := svc.GetUserByEmail(ctx, "missing@example.com")
	if err != nil {
		t.Fatalf("expected nil error for missing email, got %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func TestGetUserByEmail_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{getByEmailErr: errors.New("db error")}
	svc := newUserService(repo)

	user, err := svc.GetUserByEmail(ctx, "alice@example.com")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
	assertCustomError(t, err, 500)
}

// ---- UpdateUser ----

func TestUpdateUser_Success(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{}
	svc := newUserService(repo)

	user := &model.User{ID: uuid.New(), Username: "alice", Email: "alice@example.com"}
	if err := svc.UpdateUser(ctx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateUser_SetsUpdatedAt(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{}
	svc := newUserService(repo)

	user := &model.User{ID: uuid.New(), Username: "alice", Email: "alice@example.com"}
	before := time.Now()

	if err := svc.UpdateUser(ctx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}
	if user.UpdatedAt.Before(before) {
		t.Fatal("expected UpdatedAt to be after the call time")
	}
}

func TestUpdateUser_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{updateErr: errors.New("db error")}
	svc := newUserService(repo)

	err := svc.UpdateUser(ctx, &model.User{ID: uuid.New(), Username: "alice"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertCustomError(t, err, 500)
}

// ---- DeleteUser ----

func TestDeleteUser_Success(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{}
	svc := newUserService(repo)

	if err := svc.DeleteUser(ctx, uuid.New()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteUser_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{deleteErr: errors.New("db error")}
	svc := newUserService(repo)

	err := svc.DeleteUser(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	assertCustomError(t, err, 500)
}

// ---- GetAllUsers ----

func TestGetAllUsers_Success(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{users: []model.User{
		{ID: uuid.New(), Username: "alice", Email: "alice@example.com"},
		{ID: uuid.New(), Username: "bob", Email: "bob@example.com"},
	}}
	svc := newUserService(repo)

	users, err := svc.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestGetAllUsers_Empty(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{users: []model.User{}}
	svc := newUserService(repo)

	users, err := svc.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}

func TestGetAllUsers_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{getAllErr: errors.New("db error")}
	svc := newUserService(repo)

	users, err := svc.GetAllUsers(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if users != nil {
		t.Fatalf("expected nil users, got %v", users)
	}
	assertCustomError(t, err, 500)
}
