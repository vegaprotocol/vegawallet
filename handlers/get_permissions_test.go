package handlers_test

import (
	"fmt"
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/handlers"
	"code.vegaprotocol.io/vegawallet/handlers/mocks"
	"code.vegaprotocol.io/vegawallet/libs/jsonrpc"
	"code.vegaprotocol.io/vegawallet/permissions"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("Getting permissions with invalid params fails", testGettingPermissionsWithInvalidParamsFails)
	t.Run("Getting permissions with with valid params succeeds", testGettingPermissionsWithValidParamsSucceeds)
	t.Run("Getting permissions with with internal error fails", testGettingPermissionsWithInternalErrorFails)
}

func testGettingPermissionsWithInvalidParamsFails(t *testing.T) {
	tcs := []struct {
		name          string
		params        interface{}
		expectedError error
	}{
		{
			name:          "with nil params",
			params:        nil,
			expectedError: handlers.ErrParamsRequired,
		}, {
			name:          "with wrong type of params",
			params:        "test",
			expectedError: handlers.ErrParamsDoNotMatch,
		}, {
			name: "with empty referer",
			params: handlers.GetPermissionsParams{
				Referer: "",
			},
			expectedError: handlers.ErrRefererIsRequired,
		}, {
			name: "with invalid referer",
			params: handlers.GetPermissionsParams{
				Referer: "this is not a valid URL",
			},
			expectedError: handlers.ErrRefererIsNotValidURL,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// setup
			handler := getPermissionsForTest(tt)
			handler.permissionsStore.EXPECT().PermissionsForHostname(gomock.Any()).Times(0)

			// when
			result, errorDetails := handler.handle(t, tc.params)

			// then
			require.Empty(tt, result)
			assertInvalidParams(tt, errorDetails)
			assert.Equal(tt, tc.expectedError, errorDetails.Data)
		})
	}
}

func testGettingPermissionsWithValidParamsSucceeds(t *testing.T) {
	tcs := []struct {
		name                string
		referer             string
		expectedHostname    string
		detailedPermissions permissions.Permissions
		expectedPermissions map[string]string
	}{
		{
			name:             "With known hostname in basic referer",
			referer:          "https://token.vega.xyz",
			expectedHostname: "token.vega.xyz",
			detailedPermissions: permissions.Permissions{
				PublicKeys: &permissions.PublicKeysPermissions{
					Access: "read",
					RestrictedKeys: map[string]string{
						"my-wallet": vgrand.RandomStr(5),
					},
				},
			},
			expectedPermissions: map[string]string{
				"public_keys": "read",
			},
		}, {
			name:             "With known hostname in a rich referer",
			referer:          "https://vegawallet.xyz/auth?focus=modal#wallet",
			expectedHostname: "vegawallet.xyz",
			detailedPermissions: permissions.Permissions{
				PublicKeys: &permissions.PublicKeysPermissions{
					Access: "read",
					RestrictedKeys: map[string]string{
						"my-wallet": vgrand.RandomStr(5),
					},
				},
			},
			expectedPermissions: map[string]string{
				"public_keys": "read",
			},
		}, {
			name:                "With unknown hostname",
			referer:             "https://vegawallet.xyz/auth?focus=modal#wallet",
			expectedHostname:    "vegawallet.xyz",
			detailedPermissions: permissions.Permissions{},
			expectedPermissions: map[string]string{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// given
			params := handlers.GetPermissionsParams{
				Referer: tc.referer,
			}

			// setup
			handler := getPermissionsForTest(tt)
			handler.permissionsStore.EXPECT().PermissionsForHostname(tc.expectedHostname).Times(1).Return(tc.detailedPermissions, nil)

			// when
			result, errorDetails := handler.handle(tt, params)

			// then
			require.Nil(tt, errorDetails)
			assert.Equal(tt, tc.expectedPermissions, result.Permissions)
		})
	}
}

func testGettingPermissionsWithInternalErrorFails(t *testing.T) {
	// given
	params := handlers.GetPermissionsParams{
		Referer: "https://vega.xyz",
	}

	// setup
	handler := getPermissionsForTest(t)
	handler.permissionsStore.EXPECT().PermissionsForHostname("vega.xyz").Times(1).Return(permissions.Permissions{}, assert.AnError)

	// when
	result, errorDetails := handler.handle(t, params)

	// then
	require.Empty(t, result.Permissions)
	assertInternalError(t, errorDetails)
	assert.Equal(t, fmt.Errorf("couldn't retrieve permissions: %w", assert.AnError), errorDetails.Data)
}

func getPermissionsForTest(t *testing.T) *getPermissionsTestHandler {
	t.Helper()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockPermissionsStore(ctrl)

	return &getPermissionsTestHandler{
		GetPermissions:   handlers.NewGetPermissions(store),
		ctrl:             ctrl,
		permissionsStore: store,
	}
}

type getPermissionsTestHandler struct {
	*handlers.GetPermissions
	ctrl             *gomock.Controller
	permissionsStore *mocks.MockPermissionsStore
}

func (h *getPermissionsTestHandler) handle(t *testing.T, params interface{}) (handlers.GetPermissionsResult, *jsonrpc.ErrorDetails) {
	t.Helper()

	rawResult, err := h.Handle(params)
	if rawResult != nil {
		result, ok := rawResult.(handlers.GetPermissionsResult)
		if !ok {
			t.Fatal("GetPermissions handler result is not a GetPermissionsResult")
		}
		return result, err
	}
	return handlers.GetPermissionsResult{}, err
}
