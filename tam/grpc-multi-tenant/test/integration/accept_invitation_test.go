package integration

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	tenant_v1 "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
)

func Test_Integration_Accept_Invitation(t *testing.T) {
	type fields struct {
	}

	type testCase struct {
		args     *tenant_v1.AcceptInvitationRequest
		prepare  func(f *fields)
		expected *tenant_v1.AcceptInvitationResponse
		wantErr  bool
	}
	tests := map[string]testCase{
		"Accept_Invitation_Success": {
			args: &tenant_v1.AcceptInvitationRequest{
				MemberId: "member_0002",
			},
			prepare: func(f *fields) {
			},
			expected: &tenant_v1.AcceptInvitationResponse{
				Status: "accepted",
			},
			wantErr: false,
		},
		"Accept_Invitation_Fail_By_Empty_MemberId": {
			args: &tenant_v1.AcceptInvitationRequest{
				MemberId: "",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"Accept_Invitation_Fail_By_Member_Not_Found": {
			args: &tenant_v1.AcceptInvitationRequest{
				MemberId: "member_00010001",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
		"Accept_Invitation_Fail_By_Member_Accepted": {
			args: &tenant_v1.AcceptInvitationRequest{
				MemberId: "member_0001",
			},
			prepare:  func(f *fields) {},
			expected: nil,
			wantErr:  true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			h := CreateServiceTestHelper(t)
			resp, err := h.cli.AcceptInvitation(h.ctx, tc.args)

			f := fields{}
			if tc.prepare != nil {
				tc.prepare(&f)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("Integration_Accept_Invitation error = %v, wantErr %v", err, tc.wantErr)
			}

			if resp != nil {
				cmpCmpOpts := batchIgnoreProtoUnexportedFields(
					tenant_v1.AcceptInvitationResponse{},
				)
				ignoreFieldsOpt := cmpopts.IgnoreFields(
					tenant_v1.AcceptInvitationResponse{},
				)
				cmpCmpOpts = append(cmpCmpOpts, ignoreFieldsOpt)

				if diff := cmp.Diff(resp, tc.expected, cmpCmpOpts...); diff != "" {
					t.Errorf("Integration_Accept_Invitation value is mismatch (-actual +expected):\n%s", diff)
				}
			}
		})
	}
}
