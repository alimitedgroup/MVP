package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRoleToString(t *testing.T) {
	cases := []struct {
		UserRole
		string
	}{
		{RoleClient, "client"},
		{RoleLocalAdmin, "local_admin"},
		{RoleGlobalAdmin, "global_admin"},
	}

	for _, c := range cases {
		t.Run(c.string, func(t *testing.T) {
			res := c.UserRole.String()
			require.Equal(t, c.string, res)
		})
	}

	require.Panics(t, func() {
		_ = RoleNone.String()
	})
}

func TestRoleFromString(t *testing.T) {
	cases := []struct {
		UserRole
		string
	}{
		{RoleClient, "client"},
		{RoleLocalAdmin, "local_admin"},
		{RoleGlobalAdmin, "global_admin"},
		{RoleNone, "non_existent_role"},
	}

	for _, c := range cases {
		t.Run(c.string, func(t *testing.T) {
			res := RoleFromString(c.string)
			require.Equal(t, c.UserRole, res)
		})
	}
}
