## How to run server (Requires Docker)
```bash
docker compose up --build
```

## Endpoints
```bash
# List all funds
curl -X GET http://localhost:8080/funds

# Create a fund
curl -X POST http://localhost:8080/funds \
  -H "Content-Type: application/json" \
  -d '{"name":"Titanbay Growth Fund II","vintage_year":2025,"target_size_usd":500000000.00,"status":"Fundraising"}'

# Update a fund
# Replace {fund_id} with actual UUID string
curl -X PUT http://localhost:8080/funds \
  -H "Content-Type: application/json" \
  -d '{"id":"{fund_id}","name":"Titanbay Growth Fund III","vintage_year":2024,"target_size_usd":300000000.00,"status":"Investing"}'

# Get a specific fund
# Replace {fund_id} with actual UUID string
curl -X GET http://localhost:8080/funds/{fund_id}

# List all investors
curl -X GET http://localhost:8080/investors
# Create investor
curl -X POST http://localhost:8080/investors \
  -H "Content-Type: application/json" \
  -d '{"name":"CalPERS","investor_type":"Institution","email":"privateequity@calpers.ca.gov"}'

# Create an investment
# Replace {fund_id} and {investor_id} with actual UUID strings
curl -X POST http://localhost:8080/funds/{fund_id}/investments \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "{investor_id}",
    "amount_usd": 75000000.00,
    "investment_date": "2024-09-22"
  }'

# List investments for a specific fund
# Replace {fund_id} with actual UUID string
curl -X GET http://localhost:8080/funds/{fund_id}/investments
```