# Simple REST service in Go

Two envs can be used here:

- `RANDOM_INTEGERS_GENERATOR_TIMEOUT` with default of `10s` defines a timeout on connection with random generator
  service
- `RANDOM_INTEGERS_GENERATOR_URL` panics if not provided
    - for [random.org](https://www.random.org/integers/) set it to `https://www.random.org/integers`
    - if you wish to use `mockRandomOrg` set it to `http://localhost:8081/integers`

- `RANDOM_ORG_MOCK_SLEEP` defaults to `0`, this one is really useful for testing timeout handling

To make things easier a Makefile is there to help :)