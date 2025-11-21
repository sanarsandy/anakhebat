# Data Specifications & Seeding Plan
**Project Name:** Tukem (Tumbuh Kembang)
**Purpose:** Define the JSON structure for static medical data (WHO Standards & Milestones) to be seeded into the database.

---

## 1. WHO Growth Standards (Anthropometry)
**Source:** WHO Child Growth Standards (LMS Method).
**Format:** JSON files per indicator/gender.

### 1.1 File Structure
```
/data/who/
  ├── wfa_boys_0_5.json   (Weight-for-age)
  ├── wfa_girls_0_5.json
  ├── hfa_boys_0_5.json   (Height-for-age)
  ├── hfa_girls_0_5.json
  ├── wfh_boys_0_5.json   (Weight-for-height)
  └── wfh_girls_0_5.json
```

### 1.2 JSON Schema (LMS Table)
Each entry represents a specific age (month) or height (cm).

```json
[
  {
    "Month": 0,
    "L": 1,
    "M": 3.3464,
    "S": 0.13526,
    "SD3neg": 2.1,
    "SD2neg": 2.5,
    "SD1neg": 2.9,
    "SD0": 3.3,
    "SD1": 3.9,
    "SD2": 4.4,
    "SD3": 5.0
  },
  {
    "Month": 1,
    "L": 0.88,
    "M": 4.2,
    "S": 0.12,
    ...
  }
]
```
*Note: `L`, `M`, `S` are the variables needed for the Z-Score formula: `Z = ((value/M)^L - 1) / (L * S)`.*

---

## 2. Developmental Milestones (Checklists)
**Source:** CDC Milestones & KPSP (Kuesioner Pra Skrining Perkembangan).
**Format:** JSON file grouped by age.

### 2.1 File Structure
```
/data/milestones/
  ├── milestones_master.json
```

### 2.2 JSON Schema
```json
[
  {
    "age_months": 2,
    "description": "Milestone Usia 2 Bulan",
    "items": [
      {
        "id": 201,
        "category": "social",
        "question": "Does baby smile at people?",
        "is_red_flag": false
      },
      {
        "id": 202,
        "category": "motor",
        "question": "Can baby hold head up?",
        "is_red_flag": true
      }
    ]
  },
  {
    "age_months": 4,
    ...
  }
]
```

### 2.3 Categories
- `social`: Social & Emotional
- `language`: Language/Communication
- `cognitive`: Cognitive (Learning, Thinking, Problem-solving)
- `motor`: Movement/Physical Development

---

## 3. Implementation Strategy
1.  **Acquisition:** Download raw .txt/.csv data from WHO website.
2.  **Conversion:** Write a script (Python/JS) to convert raw data to the JSON format above.
3.  **Seeding:** Create a database seeder script (Go/SQL) to read these JSONs and insert into `who_standards` and `milestones` tables.
