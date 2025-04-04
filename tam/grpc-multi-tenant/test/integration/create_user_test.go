package integration

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	user_v1 "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/user/v1"
)

func Test_Integration_Create_User(t *testing.T) {
	type fields struct {
	}

	type testCase struct {
		args     *user_v1.CreateUserRequest
		prepare  func(f *fields)
		expected *user_v1.CreateUserResponse
		wantErr  bool
	}
	tests := map[string]testCase{
		"CreateUser_Success": {
			args: &user_v1.CreateUserRequest{
				TenantId: "tenant_0001",
				Name:     "Nguyen Tam",
				Email:    "nguyentam22@gmail.com",
			},
			prepare: func(f *fields) {
			},
			expected: &user_v1.CreateUserResponse{
				Id:        "12345",
				Name:      "Nguyen Tam",
				TenantId:  "tenant_0001",
				Email:     "nguyentam22@gmail.com",
				CreatedAt: "2006-01-02T15:04:05Z",
			},
			wantErr: false,
		},
		"CreateUser_Fail_By_Empty_TenantId": {
			args: &user_v1.CreateUserRequest{
				TenantId: "",
				Name:     "Nguyen Tam",
				Email:    "nguyentam22@gmail.com",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"CreateUser_Fail_By_Empty_Name": {
			args: &user_v1.CreateUserRequest{
				TenantId: "tenant_0001",
				Name:     "",
				Email:    "nguyentam22@gmail.com",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"CreateUser_Fail_By_Empty_Email": {
			args: &user_v1.CreateUserRequest{
				TenantId: "tenant_0001",
				Name:     "Nguyen Tam",
				Email:    "",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"CreateUser_Fail_By_Invalid_Email": {
			args: &user_v1.CreateUserRequest{
				TenantId: "tenant_0001",
				Name:     "Nguyen Tam",
				Email:    "nguyentam22@gmail",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"CreateUser_Fail_By_Tenant_Not_Found": {
			args: &user_v1.CreateUserRequest{
				TenantId: "tenant_00010001",
				Name:     "Nguyen Tam",
				Email:    "nguyentam22@gmail.com",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			h := CreateServiceTestHelper(t)
			resp, err := h.userCli.CreateUser(h.ctx, tc.args)

			f := fields{}
			if tc.prepare != nil {
				tc.prepare(&f)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("Integration_Create_User error = %v, wantErr %v", err, tc.wantErr)
			}

			if resp != nil {
				cmpCmpOpts := batchIgnoreProtoUnexportedFields(
					user_v1.CreateUserResponse{},
				)
				ignoreFieldsOpt := cmpopts.IgnoreFields(
					user_v1.CreateUserResponse{},
					"Id",
					"CreatedAt",
				)
				cmpCmpOpts = append(cmpCmpOpts, ignoreFieldsOpt)

				if diff := cmp.Diff(resp, tc.expected, cmpCmpOpts...); diff != "" {
					t.Errorf("Integration_Create_User value is mismatch (-actual +expected):\n%s", diff)
				}
			}
		})
	}
}
