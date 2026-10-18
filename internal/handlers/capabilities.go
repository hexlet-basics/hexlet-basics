package handlers

import (
	"context"
	"errors"

	"hexletbasics/internal/api"
	"hexletbasics/internal/permission"
	"hexletbasics/internal/permissioncatalog"
)

// AdminGetMyCapabilities lists the capability Permissions the caller holds
// through their Staff Role, so the admin menu can show an entry behind one (the
// audit trail) only to someone who may open it. It asks the same
// Policy.HasCapability the operations are held to, so the menu and the
// operation never disagree, and it carries no administrator shortcut: an
// administrator without the grant holds nothing here. The list is read from the
// contract's catalog, so a new capability appears without a change.
func (s *Server) AdminGetMyCapabilities(ctx context.Context) (*api.MyCapabilities, error) {
	u := MustAuthenticatedUser(ctx)
	policy := permission.NewPolicy(s.db)

	held := []string{}

	for _, entry := range permissioncatalog.Catalog() {
		if !entry.Capability {
			continue
		}

		err := policy.HasCapability(ctx, u, entry.Name)
		if errors.Is(err, permission.ErrForbidden) {
			continue
		}

		if err != nil {
			return nil, err
		}

		held = append(held, entry.Name)
	}

	return &api.MyCapabilities{Held: held}, nil
}
