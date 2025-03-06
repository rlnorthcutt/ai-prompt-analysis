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

## Project Structure

```
.
├── config
│   └── config.go       # Configuration loading and management
├── llm
│   ├── llm.go          # LLM interface definition
│   ├── claude.go       # Claude implementation
│   └── chatgpt.go      # ChatGPT implementation
├── prompt
│   └── prompt.go       # Prompt processing utilities
├── config.yaml         # Application configuration
├── .env                # Environment variables (API keys)
├── go.mod              # Go module file
├── go.sum              # Go module dependencies
├── main.go             # Application entry point
└── README.md           # Project documentation
```

## Prerequisites

- Go 1.20+
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

4. Install dependencies:

```bash
go mod download
```

## Running the Application

```bash
go run main.go
```

The server will start on port 8080 by default (configurable in `config.yaml`).

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
  "riskScore": 2
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
  "riskScore": 1
}
```

## Configuration

The application is configured using `config.yaml`. You can modify:

- Server settings (port)
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