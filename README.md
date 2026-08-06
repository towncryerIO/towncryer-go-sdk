# Towncryer SDK

The Towncryer SDK is a Go library that provides a convenient way to interact with the Towncryer API. It offers a set of API clients for various Towncryer services, such as Customers and Events, allowing you to easily manage your application's interactions with the Towncryer platform.

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
    BaseUrl:   "https://api.towncryer.com",
    ApiKey:    "your-api-key",
    ApiSecret: "your-api-secret",
}

// Create a new instance of the Towncryer SDK
sdk := towncryer.New(cfg)
```

Once you have the SDK instance, you can access the various API clients and start interacting with the Towncryer API.

### Customers API

The Customers API client provides methods for managing customer data, such as creating customers.

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

## Contributing

If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request on the [Towncryer SDK repository](https://github.com/volvlabs/towncryer-sdk).

## License

The Towncryer SDK is licensed under the [MIT License](LICENSE).
