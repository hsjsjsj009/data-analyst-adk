# Data Analyst ADK

This project is a data analyst agent built with the Go ADK (Agent Development Kit). It uses a Gemini model to analyze data and provides insights through a web interface. The agent has custom tools for OAuth2 authentication and data fetching, and it integrates with the GenAI Toolbox.

## Features

-   **Gemini Model Integration:** Leverages the power of Gemini for data analysis.
-   **OAuth2 Authentication:** Custom tools for secure user authentication.
-   **Data Fetching:** Tools to fetch data from external sources.
-   **Web Interface:** An interactive web UI for interacting with the agent.
-   **GenAI Toolbox Integration:** Extends the agent's capabilities with the GenAI Toolbox.

## Getting Started

### Prerequisites

-   Go 1.25.0 or later
-   Docker (for containerized deployment)
-   Google Cloud SDK

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/hsjsjsj009/data-analyst-adk.git
    cd data-analyst-adk
    ```

2.  **Set up your environment:**

    -   Set the `GOOGLE_API_KEY` environment variable to your Google API key.
    -   Make sure you have authenticated with Google Cloud using `gcloud auth application-default login`.

3.  **Build and run the application:**

    You can run the application locally or using Docker.

    **Local Development:**

    ```bash
    go run main.go web webui api
    ```

    **Docker:**

    ```bash
    docker build -t data-analyst-adk .
    docker run -p 8080:8080 data-analyst-adk
    ```

    The application will be available at `http://localhost:8080`.

## GenAI Toolbox

This project connects to an instance of the GenAI Toolbox. The GenAI Toolbox is an open-source MCP (Model Context Protocol) server for databases that helps in building Gen AI tools that allow agents to access data in databases.

For more information on how to set up and use the GenAI Toolbox, please refer to the [official GitHub repository](https://github.com/googleapis/genai-toolbox).

## Important Note

Please be aware that there are hardcoded URLs in the project that you may need to change for your own implementation.

### `main.go`

The GenAI Toolbox endpoint URL is hardcoded in `main.go`:

```go
Endpoint: "https://mcp-toolbox-280946129258.asia-southeast1.run.app/mcp/sse",
```

**You MUST change this URL** to point to your own GenAI Toolbox instance.

### `Dockerfile`

The `api_server_address` and `webui_address` are hardcoded in the `Dockerfile`:

```dockerfile
CMD ["web", "webui", "-api_server_address", "https://data-analyst-agent-280946129258.asia-southeast1.run.app/api", "api", "-webui_address", "https://data-analyst-agent-280946129258.asia-southeast1.run.app" ]
```

**You MUST change these URLs** to match your deployment environment for the web UI and API server.

For more information on how to use the Go ADK, please refer to the official documentation.