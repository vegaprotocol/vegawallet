package cmd_test

import (
	"strings"
	"testing"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
	"code.vegaprotocol.io/vegawallet/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Run("Initialising program succeeds", testInitialisingProgramSucceeds)
	t.Run("Forcing program initialisation with force succeeds", testForcingProgramInitialisationSucceeds)
}

func testInitialisingProgramSucceeds(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	f := &cmd.InitFlags{
		Force: false,
	}

	// when
	resp, err := cmd.Init(testDir, f)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, strings.HasPrefix(resp.RSAKeys.PublicKeyFilePath, testDir))
	assert.FileExists(t, resp.RSAKeys.PublicKeyFilePath)
	assert.True(t, strings.HasPrefix(resp.RSAKeys.PublicKeyFilePath, testDir))
	assert.FileExists(t, resp.RSAKeys.PublicKeyFilePath)
}

func testForcingProgramInitialisationSucceeds(t *testing.T) {
	testDir := NewTempDir(t)

	// given
	f := &cmd.InitFlags{
		Force: false,
	}

	// when
	resp, err := cmd.Init(testDir, f)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)

	privRSAKey1, err := vgfs.ReadFile(resp.RSAKeys.PrivateKeyFilePath)
	if err != nil {
		t.Fatalf("couldn't read private RSA key: %v", err)
	}
	pubRSAKey1, err := vgfs.ReadFile(resp.RSAKeys.PublicKeyFilePath)
	if err != nil {
		t.Fatalf("couldn't read public RSA key: %v", err)
	}

	// given
	f = &cmd.InitFlags{
		Force: true,
	}

	// when
	resp, err = cmd.Init(testDir, f)

	// then
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, strings.HasPrefix(resp.RSAKeys.PublicKeyFilePath, testDir))
	assert.FileExists(t, resp.RSAKeys.PublicKeyFilePath)
	assert.True(t, strings.HasPrefix(resp.RSAKeys.PublicKeyFilePath, testDir))
	assert.FileExists(t, resp.RSAKeys.PublicKeyFilePath)

	privRSAKey2, err := vgfs.ReadFile(resp.RSAKeys.PrivateKeyFilePath)
	if err != nil {
		t.Fatalf("couldn't read private RSA key: %v", err)
	}
	pubRSAKey2, err := vgfs.ReadFile(resp.RSAKeys.PublicKeyFilePath)
	if err != nil {
		t.Fatalf("couldn't read public RSA key: %v", err)
	}

	assert.NotEqual(t, privRSAKey1, privRSAKey2)
	assert.NotEqual(t, pubRSAKey1, pubRSAKey2)
}
