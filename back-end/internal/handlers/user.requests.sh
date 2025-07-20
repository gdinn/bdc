curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "João Silva Santos",
    "email": "joao.silva@email.com",
    "phone": "11987654321",
    "birth_date": "1985-03-15T00:00:00Z",
    "type": "RESIDENT",
    "age_group": "ADULT",
    "is_manager": false,
    "is_advisor": false
  }'


