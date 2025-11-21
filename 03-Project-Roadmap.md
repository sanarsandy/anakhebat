# Project Roadmap & Implementation Plan
**Project Name:** Aplikasi Pemantauan Tumbuh Kembang Anak (Tukem)
**Version:** 1.0
**Date:** 20 November 2025

---

## Phase 1: Foundation & MVP Core (Weeks 1-3)
**Goal:** User can login, add child, and input growth data.

### Week 1: Setup & Auth
- [ ] Setup Project Repo (Nuxt.js + Go).
- [ ] Setup Database (PostgreSQL) & Docker.
- [ ] Implement Authentication (Register, Login, JWT).
- [ ] Design Database Schema (Users, Children).

### Week 2: Child Management & Anthropometry Logic
- [ ] Implement "Add Child" feature (with Premature correction logic).
- [ ] Implement "Add Measurement" feature.
- [ ] **Core Logic:** Implement Z-Score Calculation (LMS Formula) in Backend.
- [ ] Import WHO Standard Data into Database.

### Week 3: Growth Visualization
- [ ] Integrate Chart.js in Frontend.
- [ ] Render Growth Charts (Weight/Age, Height/Age).
- [ ] Display Z-Score Status (Normal, Stunting, etc.).

---

## Phase 2: Development Tracker (Weeks 4-5)
**Goal:** User can track mental/motoric milestones.

### Week 4: Milestone Database & Checklist
- [ ] Import CDC/KPSP Milestones into Database.
- [ ] Implement "Checklist Interface" grouped by Learning Pyramid.
- [ ] Logic to determine "Red Flags".

### Week 5: Assessment Logic & History
- [ ] Save assessment results.
- [ ] Display history of development assessments.
- [ ] Implement "Warning System" for skipped levels (e.g. Cognitive > Sensory).

---

## Phase 3: Reporting & Polish (Weeks 6-7)
**Goal:** Generate Reports and improve UX.

### Week 6: Reporting Engine
- [ ] Implement PDF Generation (GoPDF).
- [ ] Create "Doctor's Report" layout (Charts + Red Flags).
- [ ] Dashboard Summary Widgets.

### Week 7: UI/UX Polish & Testing
- [ ] Responsive Design Check (Mobile View).
- [ ] Performance Optimization (Load time < 2s).
- [ ] Security Audit (Input validation, Auth checks).
- [ ] User Acceptance Testing (UAT).

---

## Phase 4: Deployment (Week 8)
- [ ] Set up Production Server (VPS/Cloud).
- [ ] Configure Domain & SSL.
- [ ] CI/CD Pipeline Setup.
- [ ] Launch.
