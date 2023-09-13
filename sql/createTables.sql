CREATE TABLE IF NOT EXISTS categories (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT
);

CREATE TABLE IF NOT EXISTS questions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  question TEXT NOT NULL,
  category_id INTEGER NOT NULL,
  correct_answer_id INTEGER NOT NULL,
  FOREIGN KEY (category_id) REFERENCES categories (id),
  FOREIGN KEY (correct_answer_id) REFERENCES options (id)
);

CREATE TABLE IF NOT EXISTS options (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  text TEXT NOT NULL,
  question_id INTEGER NOT NULL,
  FOREIGN KEY (question_id) REFERENCES questions (id)
);