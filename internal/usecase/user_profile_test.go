package usecase

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adesupraptolaia/user_login/internal/entity"
	repoMocks "github.com/adesupraptolaia/user_login/internal/repo/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_userProfile_GetUser(t *testing.T) {
	repo := repoMocks.NewUserProfilesRepo(t)

	data := entity.UserProfile{UserKsuid: "ksuid", Name: "user", DateOfBirth: "2019-01-01", Address: "Perawang"}

	repo.On("GetUserProfile", "ksuid").
		Return(&data, nil).
		Once()

	repo.On("GetUserProfile", "wrongKsuid").
		Return(nil, fmt.Errorf("user not found")).
		Once()

	type fields struct {
		repo repoMocks.UserProfilesRepo
		auth repoMocks.AuthRepo
	}
	type args struct {
		userksuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.UserProfile
		wantErr bool
	}{
		{
			name:    "Success Get User",
			fields:  fields{repo: *repo},
			args:    args{"ksuid"},
			want:    &data,
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
			uc := &userProfile{
				repo: &tt.fields.repo,
				auth: &tt.fields.auth,
			}
			got, err := uc.GetUserProfile(tt.args.userksuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("userProfile.GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userProfile.GetUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userProfile_CreateUser(t *testing.T) {
	repo := repoMocks.NewUserProfilesRepo(t)
	authRepo := repoMocks.NewAuthRepo(t)

	reqSuccess := entity.CreateUserRequest{
		Username: "user", Password: "user", UserProfile: entity.UserProfile{
			Name: "user", DateOfBirth: "2019-01-01", Address: "Perawang",
		}}

	reqFailed := entity.CreateUserRequest{
		Username: "wrongUser", Password: "wrongUser",
	}

	data := entity.UserProfile{
		UserKsuid: "ksuid", Name: "user", DateOfBirth: "2019-01-01", Address: "Perawang",
	}

	repo.On("CreateUserProfile", data).
		Return(&data, nil)

	authRepo.On("CreateUser", entity.User{Username: "user", Password: "user"}).
		Return(&entity.User{Ksuid: "ksuid", Username: "user", Password: "user"}, nil)

	authRepo.On("CreateUser", entity.User{Username: "wrongUser", Password: "wrongUser"}).
		Return(nil, fmt.Errorf("user not found"))

	type fields struct {
		repo repoMocks.UserProfilesRepo
		auth repoMocks.AuthRepo
	}
	type args struct {
		userProfileReq entity.CreateUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.UserProfile
		wantErr bool
	}{
		{
			name:    "Success Create User",
			fields:  fields{repo: *repo, auth: *authRepo},
			args:    args{reqSuccess},
			want:    &data,
			wantErr: false,
		},
		{
			name:    "Success Create User",
			fields:  fields{repo: *repo, auth: *authRepo},
			args:    args{reqFailed},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &userProfile{
				repo: &tt.fields.repo,
				auth: &tt.fields.auth,
			}
			got, err := uc.CreateUserProfile(tt.args.userProfileReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("userProfile.CreateUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userProfile.CreateUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userProfile_UpdateUser(t *testing.T) {
	repo := repoMocks.NewUserProfilesRepo(t)

	successReq := entity.UserProfile{UserKsuid: "ksuid"}
	failedReq := entity.UserProfile{UserKsuid: "wrongKsuid"}

	data := entity.UserProfile{
		UserKsuid: "ksuid", Name: "user", DateOfBirth: "2019-01-01", Address: "Perawang",
	}

	repo.On("GetUserProfile", mock.AnythingOfType("string")).
		Return(&data, nil)

	repo.On("UpdateUserProfile", successReq).
		Return(&data, nil).
		Once()

	repo.On("UpdateUserProfile", failedReq).
		Return(nil, fmt.Errorf("user not found")).
		Once()

	type fields struct {
		repo repoMocks.UserProfilesRepo
		auth repoMocks.AuthRepo
	}
	type args struct {
		userKsuid   string
		userProfile entity.UserProfile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.UserProfile
		wantErr bool
	}{
		{
			name:    "success update user",
			fields:  fields{repo: *repo},
			args:    args{"ksuid", successReq},
			want:    &data,
			wantErr: false,
		},
		{
			name:    "failed update user",
			fields:  fields{repo: *repo},
			args:    args{"wrongKsuid", failedReq},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &userProfile{
				repo: &tt.fields.repo,
				auth: &tt.fields.auth,
			}
			got, err := uc.UpdateUserProfile(tt.args.userKsuid, tt.args.userProfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("userProfile.UpdateUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userProfile.UpdateUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userProfile_DeleteUser(t *testing.T) {
	repo := repoMocks.NewUserProfilesRepo(t)
	authRepo := repoMocks.NewAuthRepo(t)

	data := entity.UserProfile{
		UserKsuid: "ksuid", Name: "user", DateOfBirth: "2019-01-01", Address: "Perawang",
	}

	repo.On("GetUserProfile", mock.AnythingOfType("string")).
		Return(&data, nil)

	repo.On("DeleteUserProfile", "ksuid").
		Return(&data, nil)

	repo.On("DeleteUserProfile", "wrongKsuid").
		Return(nil, fmt.Errorf("error not found"))

	authRepo.On("DeleteUser", mock.AnythingOfType("string")).
		Return(&entity.User{Ksuid: "ksuid", Username: "user", Password: "user"}, nil)

	type fields struct {
		repo repoMocks.UserProfilesRepo
		auth repoMocks.AuthRepo
	}
	type args struct {
		userKsuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.UserProfile
		wantErr bool
	}{
		{
			name:    "Success Delete User",
			fields:  fields{repo: *repo, auth: *authRepo},
			args:    args{"ksuid"},
			want:    &data,
			wantErr: false,
		},
		{
			name:    "Failed Delete User",
			fields:  fields{repo: *repo, auth: *authRepo},
			args:    args{"wrongKsuid"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &userProfile{
				repo: &tt.fields.repo,
				auth: &tt.fields.auth,
			}
			got, err := uc.DeleteUserProfile(tt.args.userKsuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("userProfile.DeleteUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userProfile.DeleteUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
