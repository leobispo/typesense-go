package typesense

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/v-byte-cpu/typesense-go/typesense/api"
	"github.com/v-byte-cpu/typesense-go/typesense/api/mocks"
)

func createNewKeySchema() *api.ApiKeySchema {
	return &api.ApiKeySchema{
		Description: "Search-only key.",
		Actions:     []string{"documents:search"},
		Collections: []string{"companies"},
	}
}

func createNewKey(id int64) *api.ApiKey {
	return &api.ApiKey{
		ApiKeySchema: *createNewKeySchema(),
		Id:           id,
		ValuePrefix:  "vxpx",
	}
}

func TestKeyCreate(t *testing.T) {
	newKey := createNewKeySchema()
	expectedResult := &api.ApiKey{
		ApiKeySchema: *createNewKeySchema(),
		Id:           1,
		Value:        "k8pX5hD0793d8YQC5aD1aEPd7VleSuGP",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	mockedResult := &api.ApiKey{}
	assert.Nil(t, copier.Copy(mockedResult, expectedResult))

	mockAPIClient.EXPECT().
		CreateKeyWithResponse(gomock.Not(gomock.Nil()),
			api.CreateKeyJSONRequestBody(*newKey)).
		Return(&api.CreateKeyResponse{
			JSON201: mockedResult,
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Keys().Create(newKey)

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestKeyCreateOnApiClientErrorReturnsError(t *testing.T) {
	newKey := createNewKeySchema()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	mockAPIClient.EXPECT().
		CreateKeyWithResponse(gomock.Not(gomock.Nil()),
			api.CreateKeyJSONRequestBody(*newKey)).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Keys().Create(newKey)
	assert.NotNil(t, err)
}

func TestKeyCreateOnHttpStatusErrorCodeReturnsError(t *testing.T) {
	newKey := createNewKeySchema()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	mockAPIClient.EXPECT().
		CreateKeyWithResponse(gomock.Not(gomock.Nil()),
			api.CreateKeyJSONRequestBody(*newKey)).
		Return(&api.CreateKeyResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Keys().Create(newKey)
	assert.NotNil(t, err)
}

func TestKeysRetrieve(t *testing.T) {
	expectedResult := []*api.ApiKey{
		createNewKey(1),
		createNewKey(2),
		createNewKey(3),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockClientWithResponsesInterface(ctrl)
	mockedResult := []*api.ApiKey{}
	assert.Nil(t, copier.Copy(&mockedResult, &expectedResult))

	mockAPIClient.EXPECT().
		GetKeysWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetKeysResponse{
			JSON200: &api.ApiKeysResponse{
				Keys: mockedResult,
			},
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	result, err := client.Keys().Retrieve()

	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestKeysRetrieveOnApiClientErrorReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	mockAPIClient.EXPECT().
		GetKeysWithResponse(gomock.Not(gomock.Nil())).
		Return(nil, errors.New("failed request")).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Keys().Retrieve()
	assert.NotNil(t, err)
}

func TestKeysRetrieveOnHttpStatusErrorCodeReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAPIClient := mocks.NewMockClientWithResponsesInterface(ctrl)

	mockAPIClient.EXPECT().
		GetKeysWithResponse(gomock.Not(gomock.Nil())).
		Return(&api.GetKeysResponse{
			HTTPResponse: &http.Response{
				StatusCode: 500,
			},
			Body: []byte("Internal Server error"),
		}, nil).
		Times(1)

	client := NewClient(WithAPIClient(mockAPIClient))
	_, err := client.Keys().Retrieve()
	assert.NotNil(t, err)
}
