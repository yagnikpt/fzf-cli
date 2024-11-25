# Go Fuzzy Finder

## Preview
![Go Fuzzy Finder Demo](output.gif)

## Project Description
This Go project implements a fuzzy searching algorithm and provides a user interface for interacting with the search functionality. Right now it does not perform very fast in large dataset, so run it in directory which doesn't have too many files.

## Installation Instructions
1. Clone the repository:
```
git clone https://github.com/yagnik-patel-47/go-fzf-cli.git
```
2. Navigate to the project directory:
```
cd go-fzf-cli
```
3. Install the dependencies:
```
go mod tidy
```

## Usage Examples
To run the application, use the following command:
```
go run cmd/main.go
```
To build the application, use the following command:
```
go build -o fzf.exe ./cmd/main.go
```
Once the application is running, you can input your search queries, and the application will display the matching results based on the fuzzy search algorithm.

## Contribution Guidelines
1. Fork the repository.
2. Create a new branch for your feature or bug fix:
```
git checkout -b feature/YourFeature
```
3. Make your changes and commit them:
```
git commit -m "Add your message here"
```
4. Push to the branch:
```
git push origin feature/YourFeature
```
5. Create a pull request detailing your changes.
