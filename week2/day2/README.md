# Week 2 Day 2 - Go Training

Welcome to Day 2 of Week 2 in the Go Training program. This README provides an overview of the exercises and topics covered today.

## Topics Covered

- Introduction standard Go layout project structure
- Go concurrency patterns
- Error handling in Go
- Testing in Go
- Working with `net/http` package

## Exercises
1. **Introduction standard Go layout project structure**
```sh
.
├── cmd # main applications of the project
│   └── server
├── config # configuration files
├── internal # private application and library code
│   └── pkg # shared packages within the application
│       └── worker # custom worker pool package for the application
│           └── mock # generated mocks
├── pkg # library code that's ok to use by external applications
│   └── worker # worker pool package
└── test # test folders
    ├── e2e # end-to-end tests
    └── integration # integration tests
```

2. **Concurrency Patterns**
    - Implement a worker pool
    - Use channels for communication between goroutines

3. **Error Handling in Go**
    - Handle errors in Go using errors.Join method

4. **Testing in Go**
    - Unit test
    - Integration test
        - Test a simple HTTP server with `httptest` package

5. **Working with net/http package**
    - Initialize a simple HTTP server
    - Handle graceful shutdown of the server

## Getting Started

To get started with the exercises, clone the repository and navigate to the `week2/day2` directory:

```sh
make mockgen # to generate mocks
make test # to run the tests
make run # to start the server
make submit-greeting-jobs # to submit jobs to the server
```

### AsyncGreetingService Specification

The `AsyncGreetingService` is designed to handle greeting requests. It submit incoming greeting jobs to the pool and responds once all jobs are submitted.

#### Features

- **Submit Greeting Job**: Clients can submit greeting jobs with a list names.
- **Worker Pool**: A pool of workers processes the jobs concurrently.

#### Endpoints

1. **Submit Greeting Job**
    - **URL**: `/greeting`
    - **Method**: `POST`
    - **Request Query Params**:
        ```sh
        ?name=Alice&name=Bob&name=Charlie
        ```
    - **Response**:
        ```json
        {
            "message": "OK"
        }
        ```
    - **Response Error**:
        ```json
        {
            "errors": "Failed to submit greeting jobs"
        }
        ```

#### Example Usage

1. **Submit Jobs**:
    ```sh
    curl --location --request POST 'http://localhost:8080/greeting?name=Joe&name=Lily'
    ```
Happy coding!
