package pkg

const (
	MemberStatusAccepted = "accepted"
	MemberStatusPending  = "pending"
	MemberStatusRejected = "rejected"
)

type MemberStatus string

const (
	Accepted MemberStatus = MemberStatusAccepted
	Pending  MemberStatus = MemberStatusPending
	Rejected MemberStatus = MemberStatusRejected
)
