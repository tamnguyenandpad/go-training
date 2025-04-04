package integration

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	tenant_v1 "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
)

func Test_Integration_Create_Tenant(t *testing.T) {
	type fields struct {
	}

	type testCase struct {
		args     *tenant_v1.CreateTenantRequest
		prepare  func(f *fields)
		expected *tenant_v1.CreateTenantResponse
		wantErr  bool
	}
	tests := map[string]testCase{
		"CreateTenant_Success": {
			args: &tenant_v1.CreateTenantRequest{
				Name:       "Test Tenant",
				OwnerEmail: "nguyentam@gmail.com",
			},
			prepare: func(f *fields) {
			},
			expected: &tenant_v1.CreateTenantResponse{
				Id:         "12345",
				Name:       "Test Tenant",
				OwnerEmail: "nguyentam@gmail.com",
				CreatedAt:  "2006-01-02T15:04:05Z",
			},
			wantErr: false,
		},
		"CreateTenant_Fail_By_Empty_Name": {
			args: &tenant_v1.CreateTenantRequest{
				Name:       "",
				OwnerEmail: "nguyentam@gmail.com",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"CreateTenant_Fail_By_Empty_OwnerEmail": {
			args: &tenant_v1.CreateTenantRequest{
				Name:       "Test Tenant",
				OwnerEmail: "",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"CreateTenant_Fail_By_Invalid_OwnerEmail": {
			args: &tenant_v1.CreateTenantRequest{
				Name:       "Test Tenant",
				OwnerEmail: "nguyentam@gmail",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			h := CreateServiceTestHelper(t)
			resp, err := h.cli.CreateTenant(h.ctx, tc.args)

			f := fields{}
			if tc.prepare != nil {
				tc.prepare(&f)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("Integration_Create_Tenant error = %v, wantErr %v", err, tc.wantErr)
			}

			if resp != nil {
				cmpCmpOpts := batchIgnoreProtoUnexportedFields(
					tenant_v1.CreateTenantResponse{},
				)
				ignoreFieldsOpt := cmpopts.IgnoreFields(
					tenant_v1.CreateTenantResponse{},
					"Id",
					"CreatedAt",
				)
				cmpCmpOpts = append(cmpCmpOpts, ignoreFieldsOpt)

				if diff := cmp.Diff(resp, tc.expected, cmpCmpOpts...); diff != "" {
					t.Errorf("Integration_Create_Tenant value is mismatch (-actual +expected):\n%s", diff)
				}
			}
		})
	}
}
