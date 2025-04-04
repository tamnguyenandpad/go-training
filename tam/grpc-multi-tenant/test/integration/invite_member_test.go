package integration

import (
	"testing"

	tenant_v1 "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
)

func Test_Integration_Invite_Member(t *testing.T) {
	type fields struct {
	}

	type testCase struct {
		args     *tenant_v1.InviteMemberRequest
		prepare  func(f *fields)
		expected *tenant_v1.InviteMemberResponse
		wantErr  bool
	}
	tests := map[string]testCase{
		"Invite_Member_Success": {
			args: &tenant_v1.InviteMemberRequest{
				TenantId: "tenant_0001",
				UserId:   "user_0003",
			},
			prepare: func(f *fields) {
			},
			expected: &tenant_v1.InviteMemberResponse{
				MemberId: "123",
			},
			wantErr: false,
		},
		"Invite_Member_Fail_By_Empty_TenantId": {
			args: &tenant_v1.InviteMemberRequest{
				TenantId: "",
				UserId:   "user_0002",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"Invite_Member_Fail_By_Empty_UserId": {
			args: &tenant_v1.InviteMemberRequest{
				TenantId: "tenant_0001",
				UserId:   "",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"Invite_Member_Fail_By_Tenant_Not_Found": {
			args: &tenant_v1.InviteMemberRequest{
				TenantId: "tenant_00010001",
				UserId:   "user_0002",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"Invite_Member_Fail_By_User_Not_Found": {
			args: &tenant_v1.InviteMemberRequest{
				TenantId: "tenant_0001",
				UserId:   "user_00020001",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"Invite_Member_Fail_By_User_Already_Joined": {
			args: &tenant_v1.InviteMemberRequest{
				TenantId: "tenant_0001",
				UserId:   "user_0001",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"Invite_Member_Fail_By_User_Already_Pending": {
			args: &tenant_v1.InviteMemberRequest{
				TenantId: "tenant_0005",
				UserId:   "user_0006",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			h := CreateServiceTestHelper(t)
			_, err := h.cli.InviteMember(h.ctx, tc.args)

			f := fields{}
			if tc.prepare != nil {
				tc.prepare(&f)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("Integration_Invite_Member error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
