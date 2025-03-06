import express from 'express';
import dotenv from 'dotenv';
import axios from 'axios';
import morgan from 'morgan';
import { body, validationResult } from 'express-validator';
import dbConnection from './db/database.js';

dotenv.config();
const app = express();
app.use(express.json());
app.use(morgan('dev'));

// Initialize database
const db = await dbConnection;

// Create transactions table if not exists
await db.exec(`
  CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reference TEXT UNIQUE,
    amount INTEGER,
    currency TEXT,
    status TEXT,
    customer_email TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
  )
`);

// Paystack configuration
const PAYSTACK_SECRET_KEY = process.env.PAYSTACK_SECRET_KEY;
const PAYSTACK_BASE_URL = 'https://api.paystack.co';

// Initialize payment
app.post('/initialize-payment', 
  [
    body('email').isEmail(),
    body('amount').isFloat({ gt: 0 })
  ],
  async (req, res) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) return res.status(400).json({ errors: errors.array() });

    try {
      const { email, amount } = req.body;
      const amountInCents = Math.round(amount * 100);

      // Create Paystack payment
      const paystackResponse = await axios.post(
        `${PAYSTACK_BASE_URL}/transaction/initialize`,
        {
          email,
          amount: amountInCents,
          currency: 'ZAR',
          callback_url: 'http://localhost:3000/verify-payment' // Callback URL
        },
        {
          headers: {
            Authorization: `Bearer ${PAYSTACK_SECRET_KEY}`,
            'Content-Type': 'application/json'
          }
        }
      );

      // Save to database
      await db.run(
        `INSERT INTO transactions 
        (reference, amount, currency, status, customer_email)
        VALUES (?, ?, ?, ?, ?)`,
        [
          paystackResponse.data.data.reference,
          amountInCents,
          'ZAR',
          'pending',
          email
        ]
      );

      res.json({ 
        authorizationUrl: paystackResponse.data.data.authorization_url 
      });

    } catch (error) {
      console.error('Payment error:', error.response?.data || error.message);
      res.status(500).json({ error: 'Payment initialization failed' });
    }
  }
);

// Verify payment callback (called by Paystack)
app.get('/verify-payment', async (req, res) => {
  try {
    const reference = req.query.reference;

    if (!reference) {
      return res.status(400).json({ error: 'Missing transaction reference' });
    }

    // Verify with Paystack
    const verification = await axios.get(
      `${PAYSTACK_BASE_URL}/transaction/verify/${reference}`,
      { headers: { Authorization: `Bearer ${PAYSTACK_SECRET_KEY}` } }
    );

    // Update database
    await db.run(
      `UPDATE transactions SET 
      status = ?, 
      updated_at = CURRENT_TIMESTAMP 
      WHERE reference = ?`,
      [verification.data.data.status, reference]
    );

    res.json(verification.data);

  } catch (error) {
    console.error('Verification error:', error.response?.data || error.message);
    res.status(500).json({ error: 'Payment verification failed' });
  }
});

// Verify payment manually (for testing)
app.get('/verify-payment/:reference', async (req, res) => {
  try {
    const { reference } = req.params;

    // Verify with Paystack
    const verification = await axios.get(
      `${PAYSTACK_BASE_URL}/transaction/verify/${reference}`,
      { headers: { Authorization: `Bearer ${PAYSTACK_SECRET_KEY}` } }
    );

    // Update database
    await db.run(
      `UPDATE transactions SET 
      status = ?, 
      updated_at = CURRENT_TIMESTAMP 
      WHERE reference = ?`,
      [verification.data.data.status, reference]
    );

    res.json(verification.data);

  } catch (error) {
    console.error('Verification error:', error.response?.data || error.message);
    res.status(500).json({ error: 'Payment verification failed' });
  }
});

// Get all transactions
app.get('/transactions', async (req, res) => {
    try {
      const transactions = await db.all('SELECT * FROM transactions');
      res.json(transactions);
    } catch (error) {
      console.error('Database error:', error.message);
      res.status(500).json({ error: 'Failed to fetch transactions' });
    }
  });

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});