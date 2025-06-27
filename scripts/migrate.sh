#!/bin/bash

docker exec -i rent-postgres-1 psql -U postgres -d auth_db << EOF
$(cat services/auth-service/migrations/001_init.sql)
EOF
