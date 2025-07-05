# 📦 SCORM Engine POC

**Prova de Conceito: Motor SCORM em Go + Gin + SQLite**

---

## ✅ Como rodar

# Clonar o repositório
```bash
git clone git@github.com:guido-marinho/poc_scorm.git
```
```bash
cd poc_scorm
```

# Instalar dependências Go
```bash
go mod tidy
```

# Rodar o servidor normalmente
```bash
go run cmd/server/main.go
```

# (Opcional) Instalar Air para hot reload durante o desenvolvimento
```bash
go install github.com/air-verse/air@latest
```

# Rodar com hot reload (requer Air instalado)
```bash
make dev
```

# Endpoints

📂 Upload de Pacote SCORM
**POST /upload
  -Descrição: Recebe um arquivo .zip SCORM.
  -Como usar: Envie via form-data ➜ file = [testzip.zip].

🎮 Servir o Player SCORM
**GET /packages/{package}/index.html
  -Exemplo: /packages/testzip/index.html
  -Descrição: Serve o HTML principal do curso SCORM.
***Observação: O scorm-api.js e demais arquivos (imagens, CSS) são servidos automaticamente pela mesma estrutura /packages/{package}/....

📑 Tracking de Progresso
**POST /track
  Body JSON:
  
  ```json
  {
    "userId": 1,
    "scormId": "curso-exemplo",
    "scoId": "intro",
    "status": "completed",
    "score": 100
  }
  ```
  -Descrição: Registra o progresso do aluno, SCO por SCO.

📊 Consulta de Progresso
**GET /progress/{userId}
  -Descrição: Lista todo o progresso de um usuário.

**GET /progress/{userId}/csv
  -Descrição: Exporta o progresso em CSV.

**GET /progress/{userId}/pdf
  -Descrição: Exporta o progresso em PDF. (⚠️ Para acentuação correta, use fonte TTF no futuro)

📚 Gerenciamento de Cursos
**GET /courses
  -Descrição: Lista todos os cursos cadastrados (metadados do imsmanifest.xml).

**DELETE /courses/{id}
  -Descrição: Remove um curso específico:
  -Exclui metadados do banco (courses).
  -Exclui tracking relacionado (progress).
  -Exclui os arquivos físicos na pasta storage/.

