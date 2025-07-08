# 📦 Baixa e instala o Air na versão definida
install:
	go install github.com/air-verse/air@v1.62.0

# 🚀 Roda o Air para dev com hot reload
dev:
	air

# 🏃 Roda o servidor sem hot reload
run:
	go run cmd/server/main.go

# 🔨 Compila o projeto para produção (cria binário em ./bin/app)
build:
	go build -o bin/app main.go

# 🧹 Remove binários antigos
clean:
	rm -rf bin/

# ✅ Roda todos os testes (boa prática)
test:
	go test ./...

# ℹ️ Ajuda: mostra os targets disponíveis
help:
	@echo "Comandos disponíveis:"
	@echo "  make install  - Instala o Air na versão correta"
	@echo "  make dev      - Sobe o servidor com hot reload"
	@echo "  make build    - Compila o binário para produção"
	@echo "  make clean    - Remove binários gerados"
	@echo "  make test     - Executa todos os testes"