-- Cria tabela de cursos SCORM
CREATE TABLE IF NOT EXISTS courses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  identifier TEXT NOT NULL,
  version TEXT NOT NULL,
  manifest_json TEXT NOT NULL,
  path TEXT NOT NULL
);

-- Cria tabela de progresso
CREATE TABLE IF NOT EXISTS progress (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  course_id INTEGER NOT NULL,
  sco_id TEXT,
  status TEXT,
  score INTEGER,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
