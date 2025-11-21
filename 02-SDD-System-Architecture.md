# System Design Document (SDD)
**Project Name:** Aplikasi Pemantauan Tumbuh Kembang Anak (Tukem)
**Version:** 1.0
**Date:** 20 November 2025

---

## 1. System Architecture
The system follows a standard **Client-Server Architecture** with a RESTful API.

```mermaid
graph TD
    Client[Web Client (Nuxt.js)] -->|HTTPS/JSON| LoadBalancer[Nginx / Load Balancer]
    LoadBalancer --> API[Backend API (Go/Golang)]
    API --> DB[(PostgreSQL Database)]
    API --> Cache[(Redis Cache - Optional)]
    API --> PDF[PDF Engine (GoPDF/Puppeteer)]
```

## 2. Technology Stack
### 2.1 Frontend (Web)
- **Framework:** Nuxt.js 3 (Vue 3)
- **Language:** TypeScript
- **Styling:** TailwindCSS (Recommended for speed) or Vanilla CSS
- **Charts:** Chart.js (via vue-chartjs)
- **State Management:** Pinia

### 2.2 Backend (API)
- **Language:** Go (Golang) - *Chosen for high performance & concurrency*
- **Framework:** Echo or Gin
- **Auth:** JWT (JSON Web Tokens)

### 2.3 Database
- **RDBMS:** PostgreSQL 15+
- **Why PostgreSQL?** Better support for complex queries and JSONB fields (useful for flexible milestone data).

### 2.4 Infrastructure
- **Containerization:** Docker & Docker Compose
- **OS:** Linux (Ubuntu/Alpine)

## 3. Database Schema Design (ERD Draft)

### 3.1 Users Table
| Column | Type | Description |
| :--- | :--- | :--- |
| id | UUID | Primary Key |
| email | VARCHAR | Unique, Indexed |
| password_hash | VARCHAR | Bcrypt hash |
| full_name | VARCHAR | |
| role | ENUM | 'parent', 'admin' |
| created_at | TIMESTAMP | |

### 3.2 Children Table
| Column | Type | Description |
| :--- | :--- | :--- |
| id | UUID | Primary Key |
| parent_id | UUID | FK to Users |
| name | VARCHAR | |
| dob | DATE | Date of Birth |
| gender | ENUM | 'male', 'female' |
| birth_weight | FLOAT | kg |
| birth_height | FLOAT | cm |
| is_premature | BOOLEAN | |
| gestational_age | INT | Weeks (if premature) |

### 3.3 Measurements Table (Growth)
| Column | Type | Description |
| :--- | :--- | :--- |
| id | UUID | Primary Key |
| child_id | UUID | FK to Children |
| date_measured | DATE | |
| weight | FLOAT | kg |
| height | FLOAT | cm |
| head_circ | FLOAT | cm |
| age_in_months | INT | Calculated at input |
| z_score_weight_age | FLOAT | Calculated |
| z_score_height_age | FLOAT | Calculated |
| z_score_weight_height | FLOAT | Calculated |

### 3.4 Milestones Table (Master Data)
| Column | Type | Description |
| :--- | :--- | :--- |
| id | INT | Primary Key |
| age_month | INT | Target age (e.g. 2, 4, 6) |
| category | ENUM | 'motorik', 'sensorik', 'bahasa', 'sosial' |
| question | TEXT | The checklist item |
| is_red_flag | BOOLEAN | Critical milestone? |

### 3.5 Assessments Table (Development)
| Column | Type | Description |
| :--- | :--- | :--- |
| id | UUID | Primary Key |
| child_id | UUID | FK to Children |
| milestone_id | INT | FK to Milestones |
| status | ENUM | 'yes', 'no', 'sometimes' |
| assessed_at | TIMESTAMP | |

## 4. API Endpoints Overview

### Auth
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/auth/me`

### Children
- `POST /api/children` (Create profile)
- `GET /api/children` (List all profiles)
- `GET /api/children/:id` (Detail)

### Growth
- `POST /api/children/:id/measurements`
- `GET /api/children/:id/measurements` (History for chart)

### Development
- `GET /api/milestones?age=X`
- `POST /api/children/:id/assessments` (Submit checklist)
