import express, { Router } from 'express';
import reportCSPViolation from '../controllers/reportCSPViolation';
import reportTo from '../controllers/reportTo';

const violationRouter = Router()
  .post('/report-csp-violation', express.json({ type: 'application/csp-report' }), reportCSPViolation)
  .post('/report-to', express.json({ type: 'application/reports+json' }), reportTo);

export default violationRouter;
