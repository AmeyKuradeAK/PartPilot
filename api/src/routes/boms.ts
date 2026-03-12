import { Router, Response } from 'express';
import multer from 'multer';
import path from 'path';
import { AuthRequest, authenticateToken } from '../middleware/auth';
import { query } from '../db';
import fs from 'fs';

const router = Router();

// Ensure uploads dir exists in monorepo root
const uploadsDir = path.join(__dirname, '../../../../uploads');
if (!fs.existsSync(uploadsDir)) {
  fs.mkdirSync(uploadsDir, { recursive: true });
}

const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    cb(null, uploadsDir);
  },
  filename: (req, file, cb) => {
    const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9);
    cb(null, uniqueSuffix + '-' + file.originalname);
  }
});
const upload = multer({ storage: storage });

router.use(authenticateToken);

router.post('/upload', upload.single('bom'), async (req: AuthRequest, res: Response) => {
  if (!req.file) {
    res.status(400).json({ error: 'No BOM file uploaded' });
    return;
  }

  const userId = req.user?.id;
  const filePath = req.file.path;

  try {
    // 1. Create boms record
    const bomResult = await query(
      'INSERT INTO boms (user_id, filename) VALUES ($1, $2) RETURNING id',
      [userId, filePath]
    );
    const bomId = bomResult.rows[0].id;

    // 2. Create job record for the Go engine to pick up
    // Note: Go engine parse the file and creates bom_parts
    const jobResult = await query(
      "INSERT INTO jobs (bom_id, status) VALUES ($1, 'pending') RETURNING id",
      [bomId]
    );
    const jobId = jobResult.rows[0].id;

    res.status(202).json({ 
      message: 'BOM received and processing started',
      jobId, 
      bomId 
    });
  } catch (err) {
    console.error('BOM upload error:', err);
    res.status(500).json({ error: 'Failed to process BOM upload' });
  }
});

router.post('/:id/column-mapping', async (req: AuthRequest, res: Response) => {
  const bomId = req.params.id;
  const { mapping } = req.body; // e.g., { partNumberIdx: 1, quantityIdx: 2 }

  if (!mapping) {
    res.status(400).json({ error: 'Column mapping required' });
    return;
  }

  try {
    await query(
      'UPDATE boms SET column_mapping = $1 WHERE id = $2 AND user_id = $3',
      [mapping, bomId, req.user?.id]
    );

    // If there's a failed job attached to this BOM due to missing mapping, set it back to pending
    await query(
      "UPDATE jobs SET status = 'pending', error = NULL WHERE bom_id = $1 AND status = 'failed'",
      [bomId]
    );

    res.json({ message: 'Column mapping saved, job resumed' });
  } catch (err) {
    console.error('Column mapping error:', err);
    res.status(500).json({ error: 'Failed to update mapping' });
  }
});

export default router;
