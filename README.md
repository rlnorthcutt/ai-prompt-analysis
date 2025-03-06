# Prompt Analysis API

A Go application that exposes API endpoints to analyze potential LLM prompts using Claude and ChatGPT.

## Features

- Exposes HTTP endpoints to analyze prompts with different LLM providers
- Supports Claude and ChatGPT as analysis providers
- Returns a structured JSON response containing:
  - Token count (estimated)
  - Prompt type categorization (coding, research, content creation, etc.)
  - PII detection (true/false)
  - Jailbreak attempt detection (true/false)
  - Risk assessment score (1-10)
  - Response latency (milliseconds)
- Includes an optional demo UI for testing

## Prerequisites

- Go 1.20+ (if building from source)
- Claude API key and/or OpenAI API key

## Setup

1. Clone the repository
2. Copy the example `.env` file and add your API keys:

```bash
cp .env.example .env
```

3. Edit the `.env` file:

```
CLAUDE_API_KEY=your_claude_api_key_here
OPENAI_API_KEY=your_openai_api_key_here
```

## Running the Application

### Using the pre-built binary

```bash
./ai-prompt-analysis
```

### Building and running from source

1. Install dependencies:

```bash
go mod download
```

2. Build the application:

```bash
go build -o ai-prompt-analysis
```

3. Run the application:

```bash
./ai-prompt-analysis
```

Or run directly without building:

```bash
go run main.go
```

The server will start on port 8080 by default (configurable in `config.yaml`).

## Demo UI

The application includes a browser-based UI for testing the API. To enable it, set `demoui: true` in the `server` section of your `config.yaml`:

```yaml
server:
  port: 8080
  demoui: true
```

When enabled, you can access the demo UI at:

```
http://localhost:8080/analyze
```

The UI provides:

- A form to enter prompts for analysis
- Provider selection (if both Claude and ChatGPT are available)
- Formatted display of results
- Response time tracking

## API Usage

### Analyze a Prompt with Claude

**Endpoint:** `POST /analyze/claude`

**Request:**

```json
{
  "prompt": "Your prompt text here"
}
```

**Response:**

```json
{
  "tokenCount": 42,
  "promptType": "coding",
  "containsPII": false,
  "isSuspicious": false,
  "riskScore": 2,
  "latency": 1250
}
```

### Analyze a Prompt with ChatGPT

**Endpoint:** `POST /analyze/chatgpt`

**Request:**

```json
{
  "prompt": "Your prompt text here"
}
```

**Response:**

```json
{
  "tokenCount": 45,
  "promptType": "research",
  "containsPII": false,
  "isSuspicious": false,
  "riskScore": 1,
  "latency": 890
}
```

## Configuration

The application is configured using `config.yaml`. You can modify:

- Server settings (port, demo UI)
- Claude API settings (API URL, model, tokens, temperature)
- ChatGPT API settings (API URL, model, tokens, temperature)
- Analysis system prompt

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- 400: Bad Request (invalid input)
- 405: Method Not Allowed (non-POST requests)
- 500: Internal Server Error (API errors)
- 503: Service Unavailable (API key not set)

## Security Considerations

- API keys are read from environment variables, not hardcoded
- Input validation is performed before processing
- Proper error handling to avoid leaking sensitive information

## Project Structure

```
.
├── config.yaml         # Application configuration
├── .env                # Environment variables (API keys)
├── go.mod              # Go module file
├── go.sum              # Go module dependencies
├── main.go             # Application entry point
├── README.md           # Project documentation
└── templates/             # HTML templates for Demo UI
    ├── analyze.html       # Main demo UI page
    └── result.html        # Analysis results template partial
└── internal/           # Internal packages
    ├── config/         # Configuration management
    │   └── config.go
    ├── handler/        # HTTP request handlers
    │   ├── handler.go             # Core handler functionality
    │   ├── handlerDemo.go         # Demo UI handlers
    ├── llm/            # LLM interface and implementations
    │   ├── llm.go      # Interface definition
    │   ├── claude.go   # Claude implementation
    │   └── chatgpt.go  # ChatGPT implementation
    └── prompt/         # Prompt processing utilities
        └── prompt.go
```
