curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN_HERE" \
  -d '{
    "phone": "11987654321",
    "birth_date": "1985-03-15T00:00:00Z",
    "type": "RESIDENT",
    "age_group": "ADULT"
  }'


