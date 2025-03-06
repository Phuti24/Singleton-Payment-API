// src/config/database.js
import sqlite3 from 'sqlite3';
import { open } from 'sqlite';

// Initialize SQLite database
const dbConnection = open({
  filename: './db/payments.db',
  driver: sqlite3.Database
});

export default dbConnection;