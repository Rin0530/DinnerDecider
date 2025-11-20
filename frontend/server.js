import express from 'express';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const app = express();
const PORT = process.env.PORT || 3000;

// 静的ファイルの配信
app.use(express.static(join(__dirname, 'dist')));

// すべてのルートでReactアプリを返す
app.get('*', (req, res) => {
  const indexPath = join(__dirname, 'dist', 'index.html');
  
  // index.htmlの内容を読み込んでカスタマイズ可能
  const html = fs.readFileSync(indexPath, 'utf-8');
  res.send(html);
});

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});
