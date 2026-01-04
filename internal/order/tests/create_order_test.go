package tests

import (
	"context"
	"log"
	"testing"

	sw "github.com/SInITRS/gorder/common/client/order"
	"github.com/stretchr/testify/assert"
)

var (
	ctx    = context.Background()
	server = "http://127.0.0.1:8282/api"
)

func TestMain(m *testing.M) {
	before()
	m.Run()
}

func before() {
	log.Printf("server=%s", server)
}

func TestCreateOrder_success(t *testing.T) {
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_TP5IE1MjBdjL5Z",
				Quantity: 1,
			},
			{
				Id:       "prod_TULsHeZBwnCeDH",
				Quantity: 1,
			},
		},
	})
	t.Logf("body=%s", string(response.Body))
	assert.Equal(t, 200, response.StatusCode())

	assert.Equal(t, 0, response.JSON200.Errno)
}

func TestCreateOrder_invalidParams(t *testing.T) {
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items:      nil,
	})
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, 2, response.JSON200.Errno)
}

func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIdOrdersJSONRequestBody) *sw.PostCustomerCustomerIdOrdersResponse {
	t.Helper()
	client, err := sw.NewClientWithResponses(server)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("body=%v", body)

	response, err := client.PostCustomerCustomerIdOrdersWithResponse(ctx, customerID, body)
	if err != nil {
		t.Fatal(err)
	}
	return response
}
