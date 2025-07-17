# gen-ai-proxy

**Proxy for Generative AI**


The main purpose of this proxy is to enable tracking of how many tokens were used for a given provider and model, and to expose this data to Prometheus.
It is recommended to use Grafana or any other suitable tool to visualize these metrics.

## Current **Features**
- Creating multiple providers/connections/models for single user
- Support for OpenAI Compatible providers endpoints
  - Support for LLM's /chat/completion
- Support for Ollama provider endpoint - requires authorization api-key, it is not drop-in replacement for ollama client
  - Support for LLM's /api/chat
- Exposing prometheus metrics about total tokens usage per model

### Installation
1. Install docker-compose/podman-compose
2. Copy file ``docker-compose.yml`` to your target directory
3. Update ``ENCRYPTION_KEY`` and ``JWT_SECRET`` env vars. For Encryption_key use ``head -c 32 /dev/urandom | base64 | pbcopy``
4. Run ``docker-compose up -d``
5. Access UI under ``http://localhost:8080/`` and api ``http://localhost:8080/api``

### Usage
You can use Web app for most functions.
To start using app:
1. Go to ``http://localhost:8080/``
2. Register account and login
3. Create Provider
4. Create Connection
5. Create Model
6. Create API key
7. Call API with created api key


Using bash script to achive same functionality
- Create OpenAI proxy
```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
USERNAME="gen-user"
PASSWORD="gen-password"
OPENAI_API_KEY="<your_api_key>"

# Register user
curl -X POST "$BASE_URL/api/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"'"$USERNAME"'","password":"'"$PASSWORD"'"}'

# Login
TOKEN=$(curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"'"$USERNAME"'","password":"'"$PASSWORD"'"}' | jq -r .access_token)

# Create provider
PROVIDER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/providers" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"OpenAI","base_url":"https://api.openai.com/v1","type":"openai"}')

PROVIDER_ID=$(echo "$PROVIDER_RESPONSE" | jq -r .id)

# Create connection
CONNECTION_RESPONSE=$(curl -s -X POST "$BASE_URL/api/connections" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"my-openai-connection","api_key":"'"$OPENAI_API_KEY"'","provider_id":"'"$PROVIDER_ID"'"}')
CONNECTION_ID=$(echo "$CONNECTION_RESPONSE" | jq -r .id)

# Create model
PROXY_MODEL_ID="gpt-4.1-mini-2025-04-14"
PROVIDER_MODEL_ID="gpt-4o-mini"
MODEL_NAME="GPT-4.1 Mini"

curl -X POST "$BASE_URL/api/models" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "proxy_model_id":"'"${PROXY_MODEL_ID}"'",
    "provider_model_id":"'"$PROVIDER_MODEL_ID"'",
    "name":"'"$MODEL_NAME"'",
    "connection_id":"'"$CONNECTION_ID"'",
    "thinking":true,
    "tools_usage":true,
    "price_input":0.0015,
    "price_output":0.0020
  }'

# Create API Key
API_KEY_RESPONSE=$(curl -s -X POST "$BASE_URL/api/api-keys" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"llm-test-key"}')

LLM_TEST_API_KEY=$(echo "$API_KEY_RESPONSE" | jq -r .api_key)

# LLM test call
curl -X POST "$BASE_URL/api/v1/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $LLM_TEST_API_KEY" \
  -d '{
    "model": "'"${PROXY_MODEL_ID}"'",
    "messages": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "Hello, how are you today?"}
    ]
  }' | jq .
```

---
**It is adviced to use some kind of Proxy with SSL for secure usage of app when exposed in network.**
---
Need help to deploy on your own infra? Contact us at sales@motlify.com
