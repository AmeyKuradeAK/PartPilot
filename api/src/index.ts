import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import authRouter from './routes/auth';

dotenv.config();

const app = express();
const port = process.env.PORT || 3001;

app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Routes
import bomsRouter from './routes/boms';
import jobsRouter from './routes/jobs';

app.use('/auth', authRouter);
app.use('/boms', bomsRouter);
app.use('/jobs', jobsRouter);

// Basic health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'partpilot-api', time: new Date().toISOString() });
});

app.listen(port, () => {
  console.log(`PartPilot API layer running on port ${port}`);
});
