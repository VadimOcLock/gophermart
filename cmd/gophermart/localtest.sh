gophermarttest -test.v -test.run=^TestGophermart$ \
    -gophermart-binary-path=cmd/gophermart/gophermarttest \
    -gophermart-host=localhost \
    -gophermart-port=8080 \
    -gophermart-database-uri="postgres://postgres_user:postgres_password@localhost:6430/postgres_db?sslmode=disable" \
    -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
    -accrual-host=localhost \
    -accrual-port=$(random unused-port) \
    -accrual-database-uri="***postgres/praktikum?sslmode=disable"