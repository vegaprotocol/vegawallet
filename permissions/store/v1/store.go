package v1

import (
	"fmt"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/permissions"
)

type Store struct {
	permissionsFile string
}

func (s *Store) HasPermissions(hostname string) bool {
	perms := &permissions.AllPermissions{}
	if err := paths.ReadStructuredFile(s.permissionsFile, perms); err != nil {
		return false
	}

	return perms.ExistsForHostname(hostname)
}

func (s *Store) PermissionsForHostname(hostname string) (permissions.Permissions, error) {
	allPerms, err := s.readPermissionsFile()
	if err != nil {
		return permissions.Permissions{}, err
	}

	return allPerms.PermissionsForHostname(hostname), nil
}

func (s *Store) SavePermissions(hostname string, perms permissions.Permissions) error {
	allPerms, err := s.readPermissionsFile()
	if err != nil {
		return err
	}

	allPerms.UpdatePermissions(hostname, perms)

	if err := paths.WriteStructuredFile(s.permissionsFile, allPerms); err != nil {
		return fmt.Errorf("couldn't write permissions file: %w", err)
	}

	return nil
}

func (s *Store) readPermissionsFile() (permissions.AllPermissions, error) {
	exists, err := vgfs.FileExists(s.permissionsFile)
	if err != nil {
		return permissions.AllPermissions{}, fmt.Errorf("couldn't verify permissions file existence: %w", err)
	}

	if !exists {
		return permissions.AllPermissions{}, nil
	}

	perms := &permissions.AllPermissions{}
	if err := paths.ReadStructuredFile(s.permissionsFile, perms); err != nil {
		return permissions.AllPermissions{}, fmt.Errorf("couldn't read permissions file: %w", err)
	}
	return permissions.AllPermissions{}, nil
}

func InitialiseStore(vegaPaths paths.Paths) (*Store, error) {
	permissionsFile, err := vegaPaths.CreateConfigPathFor(paths.WalletServicePermissionsConfigFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't get config path for %s: %w", paths.WalletServicePermissionsConfigFile, err)
	}

	return &Store{
		permissionsFile: permissionsFile,
	}, nil
}
