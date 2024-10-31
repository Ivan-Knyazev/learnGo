# Go-storage

Дополнительно добавлена внутренняя логика по хранению разных типов данных *(int64, string, bool, float64, complex64)* для большей вариативности сценариев использования и удобства дальнейшего применения данного хранилища.

## App `Storage`

Add ENV vars to `.env`. Example in `.env.example`.

[INPORTANT] Starting app from directory with `.env` and optionally `storage.json`

For run app from root project directory

```bash
go run cmd/storage/main.go
```

For build app and run

```bash
go build cmd/storage/main.go
./main
```

Coverage of tests for `Storage`:
```
ok  	go-storage/internal/pkg/storage	0.011s	coverage: 93.2% of statements
```

<hr>

Is a pet project to study Golang in MISIS ITAM Course.
