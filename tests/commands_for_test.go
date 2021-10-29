package tests_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"code.vegaprotocol.io/vegawallet/cmd"
	"github.com/stretchr/testify/assert"
)

func ExecuteCmd(t *testing.T, args []string) ([]byte, error) {
	t.Helper()
	var output bytes.Buffer
	w := bufio.NewWriter(&output)
	c := cmd.NewCmdRoot(w)
	c.SetArgs(args)
	execErr := c.Execute()
	if err := w.Flush(); err != nil {
		t.Fatalf("couldn't flush data out of command writer: %v", err)
	}
	return output.Bytes(), execErr
}

func Command(t *testing.T, args []string) error {
	t.Helper()
	argsWithCmd := []string{"command"}
	argsWithCmd = append(argsWithCmd, args...)
	_, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return err
	}
	return nil
}

type InitResponse struct {
	RSAKeys struct {
		PublicKeyFilePath  string `json:"publicKeyFilePath"`
		PrivateKeyFilePath string `json:"privateKeyFilePath"`
	} `json:"rsaKeys"`
	NetworksHome string `json:"networksHome"`
}

func Init(t *testing.T, args []string) (*InitResponse, error) {
	t.Helper()
	argsWithCmd := []string{"init"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &InitResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

func KeyAnnotate(t *testing.T, args []string) error {
	t.Helper()
	argsWithCmd := []string{"key", "annotate"}
	argsWithCmd = append(argsWithCmd, args...)
	_, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return err
	}
	return nil
}

type GenerateKeyResponse struct {
	Wallet struct {
		Name     string `json:"name"`
		FilePath string `json:"filePath"`
		Mnemonic string `json:"mnemonic,omitempty"`
	} `json:"wallet"`
	Key struct {
		KeyPair struct {
			PrivateKey string `json:"privateKey"`
			PublicKey  string `json:"publicKey"`
		} `json:"keyPair"`
		Algorithm struct {
			Name    string `json:"name"`
			Version uint32 `json:"version"`
		} `json:"algorithm"`
		Meta []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"meta"`
	} `json:"key"`
}

func KeyGenerate(t *testing.T, args []string) (*GenerateKeyResponse, error) {
	t.Helper()
	argsWithCmd := []string{"key", "generate"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &GenerateKeyResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type GenerateKeyAssertion struct {
	t    *testing.T
	resp *GenerateKeyResponse
}

func AssertGenerateKey(t *testing.T, resp *GenerateKeyResponse) *GenerateKeyAssertion {
	t.Helper()

	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Wallet.Name)
	assert.NotEmpty(t, resp.Wallet.FilePath)
	assert.FileExists(t, resp.Wallet.FilePath)
	assert.NotEmpty(t, resp.Key.KeyPair.PublicKey)
	assert.NotEmpty(t, resp.Key.KeyPair.PrivateKey)
	assert.NotEmpty(t, resp.Key.Algorithm.Name)
	assert.NotEmpty(t, resp.Key.Algorithm.Version)

	return &GenerateKeyAssertion{
		t:    t,
		resp: resp,
	}
}

func (a *GenerateKeyAssertion) WithWalletCreation() *GenerateKeyAssertion {
	assert.NotEmpty(a.t, a.resp.Wallet.Mnemonic)
	return a
}

func (a *GenerateKeyAssertion) WithoutWalletCreation() *GenerateKeyAssertion {
	assert.Empty(a.t, a.resp.Wallet.Mnemonic)
	return a
}

func (a *GenerateKeyAssertion) WithName(expected string) *GenerateKeyAssertion {
	assert.Equal(a.t, expected, a.resp.Wallet.Name)
	return a
}

func (a *GenerateKeyAssertion) WithMeta(expected map[string]string) *GenerateKeyAssertion {
	meta := map[string]string{}
	for _, m := range a.resp.Key.Meta {
		meta[m.Key] = m.Value
	}
	assert.Equal(a.t, expected, meta)
	return a
}

func (a *GenerateKeyAssertion) WithPublicKey(expected string) *GenerateKeyAssertion {
	assert.Equal(a.t, expected, a.resp.Key.KeyPair.PublicKey)
	return a
}

func (a *GenerateKeyAssertion) WithPrivateKey(expected string) *GenerateKeyAssertion {
	assert.Equal(a.t, expected, a.resp.Key.KeyPair.PrivateKey)
	return a
}

func (a *GenerateKeyAssertion) LocatedUnder(home string) *GenerateKeyAssertion {
	assert.True(a.t, strings.HasPrefix(a.resp.Wallet.FilePath, home), "wallet has not been generated under home directory")
	return a
}

type ListKeysResponse struct {
	Keys []struct {
		Name      string `json:"name"`
		PublicKey string `json:"publicKey"`
	} `json:"keys"`
}

func KeyList(t *testing.T, args []string) (*ListKeysResponse, error) {
	t.Helper()
	argsWithCmd := []string{"key", "list"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &ListKeysResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type IsolateKeyResponse struct {
	Wallet   string `json:"wallet"`
	FilePath string `json:"filePath"`
}

func KeyIsolate(t *testing.T, args []string) (*IsolateKeyResponse, error) {
	t.Helper()
	argsWithCmd := []string{"key", "isolate"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &IsolateKeyResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type IsolateKeyAssertion struct {
	t    *testing.T
	resp *IsolateKeyResponse
}

func AssertIsolateKey(t *testing.T, resp *IsolateKeyResponse) *IsolateKeyAssertion {
	t.Helper()

	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Wallet)
	assert.NotEmpty(t, resp.FilePath)
	assert.FileExists(t, resp.FilePath)

	return &IsolateKeyAssertion{
		t:    t,
		resp: resp,
	}
}

func (a *IsolateKeyAssertion) WithSpecialName(wallet, pubkey string) *IsolateKeyAssertion {
	assert.Equal(a.t, fmt.Sprintf("%s.%s.isolated", wallet, pubkey[0:8]), a.resp.Wallet)
	return a
}

func (a *IsolateKeyAssertion) LocatedUnder(home string) *IsolateKeyAssertion {
	assert.True(a.t, strings.HasPrefix(a.resp.FilePath, home), "wallet has not been imported under home directory")
	return a
}

func KeyTaint(t *testing.T, args []string) error {
	t.Helper()
	argsWithCmd := []string{"key", "taint"}
	argsWithCmd = append(argsWithCmd, args...)
	_, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return err
	}
	return nil
}

func KeyUntaint(t *testing.T, args []string) error {
	t.Helper()
	argsWithCmd := []string{"key", "untaint"}
	argsWithCmd = append(argsWithCmd, args...)
	_, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return err
	}
	return nil
}

type ImportNetworkResponse struct {
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
}

func NetworkImport(t *testing.T, args []string) (*ImportNetworkResponse, error) {
	t.Helper()
	argsWithCmd := []string{"network", "import"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &ImportNetworkResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type ImportNetworkAssertion struct {
	t    *testing.T
	resp *ImportNetworkResponse
}

func AssertImportNetwork(t *testing.T, resp *ImportNetworkResponse) *ImportNetworkAssertion {
	t.Helper()

	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Name)
	assert.NotEmpty(t, resp.FilePath)
	assert.FileExists(t, resp.FilePath)

	return &ImportNetworkAssertion{
		t:    t,
		resp: resp,
	}
}

func (a *ImportNetworkAssertion) WithName(expected string) *ImportNetworkAssertion {
	assert.Equal(a.t, expected, a.resp.Name)
	return a
}

func (a *ImportNetworkAssertion) LocatedUnder(home string) *ImportNetworkAssertion {
	assert.True(a.t, strings.HasPrefix(a.resp.FilePath, home), "wallet has not been imported under home directory")
	return a
}

type ListNetworksResponse struct {
	Networks []string `json:"networks"`
}

func NetworkList(t *testing.T, args []string) (*ListNetworksResponse, error) {
	t.Helper()
	argsWithCmd := []string{"network", "list"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &ListNetworksResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type SignMessageResponse struct {
	Signature string `json:"signature"`
}

func Sign(t *testing.T, args []string) (*SignMessageResponse, error) {
	t.Helper()
	argsWithCmd := []string{"sign"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &SignMessageResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type SignAssertion struct {
	t    *testing.T
	resp *SignMessageResponse
}

func AssertSign(t *testing.T, resp *SignMessageResponse) *SignAssertion {
	t.Helper()

	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Signature)

	return &SignAssertion{
		t:    t,
		resp: resp,
	}
}

func (a *SignAssertion) WithSignature(expected string) *SignAssertion {
	assert.Equal(a.t, expected, a.resp.Signature)
	return a
}

type VerifyMessageResponse struct {
	IsValid bool `json:"isValid"`
}

func Verify(t *testing.T, args []string) (*VerifyMessageResponse, error) {
	t.Helper()
	argsWithCmd := []string{"verify"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &VerifyMessageResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type VerifyAssertion struct {
	t    *testing.T
	resp *VerifyMessageResponse
}

func AssertVerify(t *testing.T, resp *VerifyMessageResponse) *VerifyAssertion {
	t.Helper()

	assert.NotNil(t, resp)

	return &VerifyAssertion{
		t:    t,
		resp: resp,
	}
}

func (a *VerifyAssertion) IsValid() *VerifyAssertion {
	assert.True(a.t, a.resp.IsValid)
	return a
}

type ImportWalletResponse struct {
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
}

func WalletImport(t *testing.T, args []string) (*ImportWalletResponse, error) {
	t.Helper()
	argsWithCmd := []string{"import"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &ImportWalletResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type ImportWalletAssertion struct {
	t    *testing.T
	resp *ImportWalletResponse
}

func AssertImportWallet(t *testing.T, resp *ImportWalletResponse) *ImportWalletAssertion {
	t.Helper()

	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Name)
	assert.NotEmpty(t, resp.FilePath)
	assert.FileExists(t, resp.FilePath)

	return &ImportWalletAssertion{
		t:    t,
		resp: resp,
	}
}

func (a *ImportWalletAssertion) WithName(expected string) *ImportWalletAssertion {
	assert.Equal(a.t, expected, a.resp.Name)
	return a
}

func (a *ImportWalletAssertion) LocatedUnder(home string) *ImportWalletAssertion {
	assert.True(a.t, strings.HasPrefix(a.resp.FilePath, home), "wallet has not been imported under home directory")
	return a
}

type GetWalletInfoResponse struct {
	Type    string `json:"type"`
	Version uint32 `json:"version"`
	ID      string `json:"id"`
}

func WalletInfo(t *testing.T, args []string) (*GetWalletInfoResponse, error) {
	t.Helper()
	argsWithCmd := []string{"info"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &GetWalletInfoResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}

type GetWalletInfoAssertion struct {
	t    *testing.T
	resp *GetWalletInfoResponse
}

func AssertWalletInfo(t *testing.T, resp *GetWalletInfoResponse) *GetWalletInfoAssertion {
	t.Helper()

	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Type)
	assert.NotEmpty(t, resp.Version)
	assert.NotEmpty(t, resp.ID)

	return &GetWalletInfoAssertion{
		t:    t,
		resp: resp,
	}
}

func (a *GetWalletInfoAssertion) IsHDWallet() *GetWalletInfoAssertion {
	assert.Equal(a.t, "HD wallet", a.resp.Type)
	return a
}

func (a *GetWalletInfoAssertion) IsIsolatedHDWallet() *GetWalletInfoAssertion {
	assert.Equal(a.t, "HD wallet (isolated)", a.resp.Type)
	return a
}

func (a *GetWalletInfoAssertion) WithLatestVersion() *GetWalletInfoAssertion {
	assert.Equal(a.t, uint32(2), a.resp.Version)
	return a
}

func (a *GetWalletInfoAssertion) WithVersion(i int) *GetWalletInfoAssertion {
	assert.Equal(a.t, uint32(i), a.resp.Version)
	return a
}

type ListWalletsResponse struct {
	Wallets []string `json:"wallets"`
}

func WalletList(t *testing.T, args []string) (*ListWalletsResponse, error) {
	t.Helper()
	argsWithCmd := []string{"list"}
	argsWithCmd = append(argsWithCmd, args...)
	output, err := ExecuteCmd(t, argsWithCmd)
	if err != nil {
		return nil, err
	}
	resp := &ListWalletsResponse{}
	if err := json.Unmarshal(output, resp); err != nil {
		t.Fatalf("couldn't unmarshal command output: %v", err)
	}
	return resp, nil
}
