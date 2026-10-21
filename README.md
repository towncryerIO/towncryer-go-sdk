# Towncryer SDK

The Towncryer SDK is a Go library that provides a convenient way to interact with the Towncryer API. It offers a set of API clients for various Towncryer services, such as Customers, Events, and Messages, allowing you to easily manage your application's interactions with the Towncryer platform.

## Installation

To use the Towncryer SDK in your Go project, you can install it using the following command:

```
go get github.com/volvlabs/towncryer-sdk
```

## Usage

To get started with the Towncryer SDK, you need to create a new instance of the `TowncryerSdk` and configure it with the necessary details, such as the API credentials and the base URL.

```go
import (
    "github.com/volvlabs/towncryer-sdk/apis"
    "github.com/volvlabs/towncryer-sdk/config"
)

// Create a new configuration
cfg := &config.Configuration{
    BaseURL:      "https://api.towncryer.com",
    APIKey:       "your-api-key",
    APISecret:    "your-api-secret",
}

// Create a new instance of the Towncryer SDK
sdk := towncryer.New(cfg)
```

Once you have the SDK instance, you can access the various API clients and start interacting with the Towncryer API.

### Customers API

The Customers API client provides methods for managing customer data, such as creating, retrieving, and updating customers.

```go
// Create a new customer
customerPayload := &models.CreateCustomerPayload{
    CustomerId:  "CUST-12345",
    FirstName:   "John",
    LastName:    "Doe",
    Email:       "john.doe@example.com",
    PhoneNumber: "123-456-7890",
}

response, err := sdk.CustomersApi().CreateCustomer(ctx, customerPayload)
if err != nil {
    // Handle error
}

// Retrieve a customer
customer, err := sdk.CustomersApi().GetCustomer(ctx, "CUST-12345")
if err != nil {
    // Handle error
}
```

### Events API

The Events API client provides methods for publishing and managing events within the Towncryer platform.

```go
// Publish an event
eventPayload := &models.PublishEventPayload{
    Name: "UserRegisteredEvent",
    Customer: models.CreateCustomerPayload{
        CustomerId: "CUST-1234",
        FirstName:  "John",
        LastName:   "Doe",
    },
    Data: map[string]any{
        "testKey": "testValue",
    },
}

response, err := sdk.EventsApi().PublishEvent(ctx, eventPayload)
if err != nil {
    // Handle error
}
```

### Messages API

The Messages API client provides methods for sending emails, SMS, and push notifications in a single bulk call.

```go
// Send a bulk batch of messages
messagesPayload := &models.SendMessagesPayload{
    Emails: []models.SendEmailPayload{
        {
            From:  models.EmailFrom{Name: "Acme", Email: "no-reply@acme.com"},
            Recipients: []models.EmailRecipient{
                {Name: "John Doe", Email: "john.doe@example.com"},
            },
            Title:   "Welcome!",
            Content: "Thanks for signing up.",
        },
    },
}

scheduleInfo, err := sdk.MessagesApi().SendBulkMessages(ctx, messagesPayload)
if err != nil {
    // Handle error
}

fmt.Println(scheduleInfo.JobId, scheduleInfo.Status)
```

## Error Handling

Every SDK call returns a plain `error`, but errors coming back from the Towncryer API itself are always wrapped in a `*towncryer.APIError`, which carries the HTTP status code, the API's own error code and message, and any field-level validation errors from the response body. This lets you branch on the failure without string-matching the error message:

```go
import (
    "errors"

    "github.com/volvlabs/towncryer-sdk"
)

response, err := sdk.CustomersApi().CreateCustomer(ctx, customerPayload)
if err != nil {
    var apiErr *towncryer.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.StatusCode {
        case 401:
            // bad credentials
        case 400:
            // validation failure — apiErr.Errors has the field-level detail
        default:
            // other API-level failure
        }
        return
    }

    // err is not an *APIError, so it's a network-level failure
    // (e.g. connection refused, timeout, DNS failure) rather than
    // a response from the Towncryer API.
}
```

`errors.As` returns `false` for network-level failures (the request never reached the API, or never got a response), so those remain distinguishable from API-level failures without any extra type on your end.

