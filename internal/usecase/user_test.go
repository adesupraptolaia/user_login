package usecase

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adesupraptolaia/user_login/internal/entity"
	repoMocks "github.com/adesupraptolaia/user_login/internal/repo/mocks"
	"github.com/adesupraptolaia/user_login/internal/utils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/mock"
)

func Test_user_GetUserByUsername(t *testing.T) {
	repo := repoMocks.NewUsersRepo(t)

	data := &entity.User{Ksuid: ksuid.New().String(), Username: "user", Password: "user", Role: entity.USER}

	repo.On("GetUserByUsername", "user").
		Return(data, nil).
		Once()

	repo.On("GetUserByUsername", "toni").
		Return(nil, fmt.Errorf("not found")).
		Once()

	type fields struct {
		repo repoMocks.UsersRepo
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name:    "Success Get user",
			fields:  fields{repo: *repo},
			args:    args{"user"},
			want:    data,
			wantErr: false,
		},
		{
			name:    "Faield Get user",
			fields:  fields{repo: *repo},
			args:    args{"toni"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &user{
				repo: &tt.fields.repo,
			}
			got, err := uc.GetUserByUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.GetUserByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_GetUserByKsuid(t *testing.T) {
	repo := repoMocks.NewUsersRepo(t)

	data := &entity.User{Ksuid: "ksuid", Username: "user", Password: "user", Role: entity.USER}

	repo.On("GetUserByKsuid", "ksuid").
		Return(data, nil).
		Once()

	repo.On("GetUserByKsuid", "wrongKsuid").
		Return(nil, fmt.Errorf("not found")).
		Once()

	type fields struct {
		repo repoMocks.UsersRepo
	}
	type args struct {
		ksuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name:    "Success Get User",
			fields:  fields{repo: *repo},
			args:    args{"ksuid"},
			want:    data,
			wantErr: false,
		},
		{
			name:    "Failed Get User",
			fields:  fields{repo: *repo},
			args:    args{"wrongKsuid"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &user{
				repo: &tt.fields.repo,
			}
			got, err := uc.GetUserByKsuid(tt.args.ksuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.GetUserByKsuid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.GetUserByKsuid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_CreateUser(t *testing.T) {
	repo := repoMocks.NewUsersRepo(t)

	data := &entity.User{Username: "user", Password: utils.HashPassword("user"), Role: entity.USER}

	repo.On("CreateUser", mock.AnythingOfType("entity.User")).
		Return(data, nil).
		Once()

	repo.On("GetUserByUsername", "user").
		Return(nil, fmt.Errorf("user not found")).
		Once()

	type fields struct {
		repo repoMocks.UsersRepo
	}
	type args struct {
		userReq entity.UserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name:    "Success Get User",
			fields:  fields{repo: *repo},
			args:    args{entity.UserRequest{Username: "user", Password: "user"}},
			want:    data,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &user{
				repo: &tt.fields.repo,
			}
			got, err := uc.CreateUser(tt.args.userReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_UpdateUser(t *testing.T) {
	repo := repoMocks.NewUsersRepo(t)

	data := &entity.User{Ksuid: "ksuid", Username: "user", Password: utils.HashPassword("user"), Role: entity.USER}

	repo.On("UpdateUser", "ksuid", *data).
		Return(data, nil).
		Once()

	repo.On("GetUserByKsuid", "ksuid").
		Return(data, nil).
		Once()

	repo.On("GetUserByKsuid", "wrongKsuid").
		Return(nil, fmt.Errorf("user not found")).
		Once()

	type fields struct {
		repo repoMocks.UsersRepo
	}
	type args struct {
		ksuid string
		user  entity.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name:    "Success Get User",
			fields:  fields{repo: *repo},
			args:    args{"ksuid", *data},
			want:    data,
			wantErr: false,
		},
		{
			name:    "Success Get User",
			fields:  fields{repo: *repo},
			args:    args{"wrongKsuid", *data},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &user{
				repo: &tt.fields.repo,
			}
			got, err := uc.UpdateUser(tt.args.ksuid, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.UpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_DeleteUser(t *testing.T) {
	repo := repoMocks.NewUsersRepo(t)

	data := &entity.User{Ksuid: "ksuid", Username: "user", Password: utils.HashPassword("user"), Role: entity.USER}

	repo.On("DeleteUser", "ksuid").
		Return(data, nil).
		Once()

	repo.On("GetUserByKsuid", "ksuid").
		Return(data, nil).
		Once()

	repo.On("GetUserByKsuid", "wrongKsuid").
		Return(nil, fmt.Errorf("user not found")).
		Once()

	type fields struct {
		repo repoMocks.UsersRepo
	}
	type args struct {
		ksuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name:    "Success Get User",
			fields:  fields{repo: *repo},
			args:    args{"ksuid"},
			want:    data,
			wantErr: false,
		},
		{
			name:    "Success Get User",
			fields:  fields{repo: *repo},
			args:    args{"wrongKsuid"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &user{
				repo: &tt.fields.repo,
			}
			got, err := uc.DeleteUser(tt.args.ksuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
