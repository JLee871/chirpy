package auth

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

const secretKey = "secret-key"

func TestValidateJWT(t *testing.T) {
	type args struct {
		tokenString string
		tokenSecret string
	}

	//Making a valid token and an expired token
	uuid1 := uuid.New()
	uuid2 := uuid.New()
	token1, err := MakeJWT(uuid1, secretKey, 5*time.Minute)
	if err != nil {
		t.Errorf("MakeJWT error, failed to create a token")
	}
	token2, _ := MakeJWT(uuid2, secretKey, -5*time.Minute)

	tests := []struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		// Add test cases here
		{
			name:    "Valid Token",
			args:    args{tokenString: token1, tokenSecret: secretKey},
			want:    uuid1,
			wantErr: false,
		},
		{
			name:    "Incorrect Key",
			args:    args{tokenString: token1, tokenSecret: "wrong-key"},
			want:    uuid.Nil,
			wantErr: true,
		},
		{
			name:    "Expired Token",
			args:    args{tokenString: token2, tokenSecret: secretKey},
			want:    uuid.Nil,
			wantErr: true,
		},
		{
			name:    "Incorrect Token",
			args:    args{tokenString: "this-is-an-incorrect-token", tokenSecret: secretKey},
			want:    uuid.Nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJWT(tt.args.tokenString, tt.args.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	type args struct {
		headers http.Header
	}

	//Making a http.header
	testHeader := http.Header{}
	testHeader.Add("Authorization", "Bearer MyToken")

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// Add test cases here
		{
			name:    "Correct",
			args:    args{headers: testHeader},
			want:    "MyToken",
			wantErr: false,
		},
		{
			name:    "Empty Case",
			args:    args{},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
