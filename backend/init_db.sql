-- ============================================
-- TUKEM DATABASE INITIALIZATION SCRIPT
-- This file runs automatically when PostgreSQL container starts for the first time
-- ============================================

-- Create Extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- 1. USERS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    full_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'parent',
    google_id VARCHAR(255) UNIQUE,
    auth_provider VARCHAR(50) DEFAULT 'email',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- ============================================
-- 2. CHILDREN TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS children (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    dob DATE NOT NULL,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
    birth_weight FLOAT,
    birth_height FLOAT,
    is_premature BOOLEAN DEFAULT FALSE,
    gestational_age INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_children_parent_id ON children(parent_id);

-- ============================================
-- 3. MEASUREMENTS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS measurements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    child_id UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    measurement_date DATE NOT NULL,
    weight DECIMAL(5,2) NOT NULL,
    height DECIMAL(6,2) NOT NULL,
    head_circumference DECIMAL(5,2),
    age_in_days INT NOT NULL,
    age_in_months INT NOT NULL,
    storage_days INT,
    storage_months INT,
    weight_for_age_zscore DECIMAL(5,2),
    height_for_age_zscore DECIMAL(5,2),
    weight_for_height_zscore NUMERIC(5,2),
    head_circumference_zscore NUMERIC(5,2),
    weight_status VARCHAR(50),
    height_status VARCHAR(50),
    nutritional_status VARCHAR(100),
    weight_for_height_status VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_measurements_child_id ON measurements(child_id);
CREATE INDEX IF NOT EXISTS idx_measurements_date ON measurements(measurement_date);

-- ============================================
-- 4. MILESTONES TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS milestones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    age_months INT NOT NULL,
    min_age_range INT,
    max_age_range INT,
    category VARCHAR(50) NOT NULL,
    question TEXT NOT NULL,
    question_en TEXT,
    source VARCHAR(50) DEFAULT 'KPSP',
    is_red_flag BOOLEAN DEFAULT FALSE,
    pyramid_level INT NOT NULL,
    denver_domain VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_milestones_age ON milestones(age_months);
CREATE INDEX IF NOT EXISTS idx_milestones_range ON milestones(min_age_range, max_age_range);
CREATE INDEX IF NOT EXISTS idx_milestones_category ON milestones(category);
CREATE INDEX IF NOT EXISTS idx_milestones_denver_domain ON milestones(denver_domain) WHERE denver_domain IS NOT NULL;

-- ============================================
-- 5. ASSESSMENTS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    child_id UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    milestone_id UUID NOT NULL REFERENCES milestones(id) ON DELETE CASCADE,
    assessment_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(child_id, milestone_id)
);

CREATE INDEX IF NOT EXISTS idx_assessments_child ON assessments(child_id);
CREATE INDEX IF NOT EXISTS idx_assessments_milestone ON assessments(milestone_id);

-- ============================================
-- 6. WHO STANDARDS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS who_standards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    indicator VARCHAR(50) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    age_months INT,
    height_cm DECIMAL(5,2),
    l_value DECIMAL(10,6) NOT NULL,
    m_value DECIMAL(10,6) NOT NULL,
    s_value DECIMAL(10,6) NOT NULL,
    sd3neg DECIMAL(10,4),
    sd2neg DECIMAL(10,4),
    sd1neg DECIMAL(10,4),
    sd0 DECIMAL(10,4),
    sd1 DECIMAL(10,4),
    sd2 DECIMAL(10,4),
    sd3 DECIMAL(10,4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(indicator, gender, age_months, height_cm)
);

CREATE INDEX IF NOT EXISTS idx_who_indicator_gender ON who_standards(indicator, gender);
CREATE INDEX IF NOT EXISTS idx_who_age ON who_standards(age_months) WHERE age_months IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_who_height ON who_standards(height_cm) WHERE height_cm IS NOT NULL;

-- ============================================
-- 7. STIMULATION CONTENT TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS stimulation_content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    milestone_id UUID REFERENCES milestones(id) ON DELETE SET NULL,
    category VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content_type VARCHAR(20) NOT NULL DEFAULT 'article',
    url TEXT NOT NULL,
    thumbnail_url TEXT,
    age_min_months INT,
    age_max_months INT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stimulation_content_milestone ON stimulation_content(milestone_id);
CREATE INDEX IF NOT EXISTS idx_stimulation_content_category ON stimulation_content(category);
CREATE INDEX IF NOT EXISTS idx_stimulation_content_age_range ON stimulation_content(age_min_months, age_max_months);
CREATE INDEX IF NOT EXISTS idx_stimulation_content_active ON stimulation_content(is_active);

-- ============================================
-- 8. IMMUNIZATION SCHEDULE TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS immunization_schedule (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    name_id VARCHAR(255),
    description TEXT,
    age_min_days INT,
    age_optimal_days INT,
    age_max_days INT,
    age_min_months INT,
    age_optimal_months INT,
    age_max_months INT,
    dose_number INT NOT NULL,
    total_doses INT,
    interval_from_previous_days INT,
    interval_from_previous_months INT,
    category VARCHAR(50) DEFAULT 'wajib',
    priority VARCHAR(20) DEFAULT 'medium',
    is_required BOOLEAN DEFAULT TRUE,
    notes TEXT,
    source VARCHAR(100) DEFAULT 'IDAI',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_immunization_schedule_age ON immunization_schedule(age_min_days, age_max_days);
CREATE INDEX IF NOT EXISTS idx_immunization_schedule_category ON immunization_schedule(category);
CREATE INDEX IF NOT EXISTS idx_immunization_schedule_active ON immunization_schedule(is_active);

-- ============================================
-- 9. CHILD IMMUNIZATIONS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS child_immunizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    child_id UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    immunization_schedule_id UUID NOT NULL REFERENCES immunization_schedule(id),
    given_date DATE NOT NULL,
    given_at_age_days INT,
    given_at_age_months INT,
    location VARCHAR(255),
    healthcare_facility VARCHAR(255),
    doctor_name VARCHAR(255),
    vaccine_batch_number VARCHAR(100),
    notes TEXT,
    is_on_schedule BOOLEAN,
    is_catch_up BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_child_immunizations_unique ON child_immunizations(child_id, immunization_schedule_id);
CREATE INDEX IF NOT EXISTS idx_child_immunizations_child ON child_immunizations(child_id);
CREATE INDEX IF NOT EXISTS idx_child_immunizations_date ON child_immunizations(given_date);
CREATE INDEX IF NOT EXISTS idx_child_immunizations_schedule ON child_immunizations(immunization_schedule_id);

-- ============================================
-- DONE! Tables created automatically on first start
-- ============================================
