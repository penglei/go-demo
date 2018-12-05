package service_test

import (
	"github.com/qcloud2018/go-demo/service"
	"github.com/qcloud2018/go-demo/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// Test_AddContact shows how to structure a basic service test. Notice the 'SETUP', 'TEST', and 'VERIFY' steps that
// nearly all tests should have.
func Test_AddContact(t *testing.T) {
	// SETUP:
	// A standard Env. defer is used to ensure the env is cleaned up after the test.
	env := test.SetupEnv(t)
	defer env.Close()

	// TEST: Adding a contact via the API.
	contact, err := env.Client.AddContact(service.AddContactRequest{
		Email: "alice@example.xyz",
		Name:  "Alice Zulu",
	})

	// VERIFY: Response contains the contact
	require.NoError(t, err, "Unable to get contact via API")
	require.NotEmpty(t, contact, "Contact not found")
	assert.True(t, contact.ID > 0, "Contact ID is missing")
	assert.Equal(t, contact.Email, "alice@example.xyz")
	assert.Equal(t, contact.Name, "Alice Zulu")

	// VERIFY: Contact is added to the database properly.
	dbContact := env.ReadContactWithEmail("alice@example.xyz")
	require.NotEmpty(t, dbContact, "Contact not found")
	assert.Equal(t, dbContact.Email, "alice@example.xyz")
	assert.Equal(t, dbContact.Name, "Alice Zulu")
}

func Test_GetContactByEmail(t *testing.T) {
	env := test.SetupEnv(t)
	defer env.Close()

	var email = "alice2@example.xyz"
	// SETUP:
	env.SetupContact(email, "Alice Zulu")

	// -------------------------------------------------------------------------------------------------------------
	// TEST: when contact exists
	{
		// Using braces like this can help isolate different test cases.

		contact, err := env.Client.GetContactByEmail(email)

		// VERIFY: Response contains the contact
		require.NoError(t, err, "Unable to get contact via API")
		require.NotEmpty(t, contact, "Contact not found")
		assert.True(t, contact.ID > 0, "Contact ID is missing")
		assert.Equal(t, contact.Email, email)
		assert.Equal(t, contact.Name, "Alice Zulu")
	}

	// -------------------------------------------------------------------------------------------------------------
	// TEST: when contact doesn't exist
	{
		contact, err := env.Client.GetContactByEmail("bob@example.xyz")

		// VERIFY: 404 Not Found returned
		require.Error(t, err)
		require.IsType(t, service.ErrorResponse{}, err)

		assert.Equal(t, http.StatusNotFound, err.(service.ErrorResponse).StatusCode)

		assert.Nil(t, contact)
	}

}
