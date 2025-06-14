Just trying Azure AI with the Go language.

## Overview

This repository demonstrates how to interact with Azure AI services using the Go programming language. It’s a simple experiment to explore the capabilities of Azure AI and Go integration.

## Features

- Connects to Azure AI services from Go
- Example usage and basic API calls
- Serves as a starting point for further Azure AI experiments in Go

## Prerequisites

- Go (version 1.18 or higher recommended)
- An Azure account with access to Azure AI services
- Azure AI resource credentials (endpoint and API key)

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/AbdelrhmanSaid/go-ai.git
   cd go-ai
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**

   Create a `.env` file or export the following variables in your shell:
   ```
   AZURE_AI_ENDPOINT=your_azure_ai_endpoint
   AZURE_AI_KEY=your_azure_ai_api_key
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

## Usage

Update the main.go file (or relevant code) with your specific Azure AI logic and credentials. The current implementation serves as a template for making authenticated requests to Azure AI services.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have suggestions or improvements.

## License

This project is licensed under the MIT License.

---

Feel free to modify this template to better fit your project's specifics! Let me know if you want to include code examples or more details about what the repo does.
