# Go Temperature from CEP API Server

This project is a simple API server built using Go (Golang) to search for CEP and returning the information of the weather temperature in: Celsius, Fahrenheit and Kelvin.

## Getting Started

To run the API server, follow these steps:

1. Clone the repository:
   ```
   git clone <repository-url>
   cd pos-fullcycle-deploy-cloud-run
   ```

2. Install the dependencies:
   ```
   go mod tidy
   ```

3. Run the server:
   ```
   cd api-server/cmd
   WEATHER_API_KEY=8dc0a8c0 go run cmd/main.go
   ```

4. Run the tests:
   ```
   docker-compose run api-server-tests
   ```

## Não incluido no Cloud Run pois o mesmo está cobrando um pré-pagamento

## API Endpoints

- **GET /tempForCep/:cep**: Return the Temperature for the informed CEP if available.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.