package permissions

type AllPermissions struct {
	Version     uint32                 `json:"version"`
	Permissions map[string]Permissions `json:"permissions"`
}

func (p *AllPermissions) UpdatePermissions(hostname string, permissions Permissions) {
	p.Permissions[hostname] = permissions
}

func (p *AllPermissions) PermissionsForHostname(hostname string) Permissions {
	return p.Permissions[hostname]
}

func (p *AllPermissions) ExistsForHostname(hostname string) bool {
	_, ok := p.Permissions[hostname]
	return ok
}

type Permissions struct {
	PublicKeys *PublicKeysPermissions `json:"publicKeys,omitempty"`
}

func (p *Permissions) Summary() map[string]string {
	summary := map[string]string{}
	if p.PublicKeys != nil {
		summary["public_keys"] = p.PublicKeys.Access
	}
	return summary
}

type AccessMode string

var ReadAccess AccessMode = "read"

type PublicKeysPermissions struct {
	Access         AccessMode        `json:"access"`
	RestrictedKeys map[string]string `json:"restrictedKeys"`
}
