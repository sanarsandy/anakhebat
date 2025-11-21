-- Migration: Create WHO Standards Table
-- This table stores WHO Child Growth Standards LMS parameters
-- for calculating Z-scores and growth assessment

CREATE TABLE IF NOT EXISTS who_standards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Indicator type: wfa (weight-for-age), hfa (height-for-age), 
    -- wfh (weight-for-height), hcfa (head circumference-for-age)
    indicator VARCHAR(50) NOT NULL,
    
    -- Gender: male or female
    gender VARCHAR(10) NOT NULL,
    
    -- Age in months (0-60) for age-based indicators
    -- NULL for weight-for-height which uses height instead
    age_months INT,
    
    -- Height in cm for weight-for-height indicator
    -- NULL for age-based indicators
    height_cm DECIMAL(5,2),
    
    -- LMS parameters for Z-score calculation
    -- Formula: Z = ((value/M)^L - 1) / (L * S)
    l_value DECIMAL(10,6) NOT NULL,  -- Box-Cox transformation
    m_value DECIMAL(10,6) NOT NULL,  -- Median
    s_value DECIMAL(10,6) NOT NULL,  -- Coefficient of variation
    
    -- Pre-calculated standard deviation values for reference
    sd3neg DECIMAL(10,4),  -- -3 SD
    sd2neg DECIMAL(10,4),  -- -2 SD
    sd1neg DECIMAL(10,4),  -- -1 SD
    sd0 DECIMAL(10,4),     -- Median (same as M)
    sd1 DECIMAL(10,4),     -- +1 SD
    sd2 DECIMAL(10,4),     -- +2 SD
    sd3 DECIMAL(10,4),     -- +3 SD
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure unique combinations
    UNIQUE(indicator, gender, age_months, height_cm)
);

-- Indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_who_indicator_gender ON who_standards(indicator, gender);
CREATE INDEX IF NOT EXISTS idx_who_age ON who_standards(age_months) WHERE age_months IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_who_height ON who_standards(height_cm) WHERE height_cm IS NOT NULL;

-- Add Z-score and status columns to measurements table
ALTER TABLE measurements 
ADD COLUMN IF NOT EXISTS weight_for_height_zscore NUMERIC(5,2),
ADD COLUMN IF NOT EXISTS head_circumference_zscore NUMERIC(5,2),
ADD COLUMN IF NOT EXISTS nutritional_status VARCHAR(100),
ADD COLUMN IF NOT EXISTS height_status VARCHAR(100),
ADD COLUMN IF NOT EXISTS weight_for_height_status VARCHAR(100);

-- Add comment for documentation
COMMENT ON TABLE who_standards IS 'WHO Child Growth Standards LMS parameters for Z-score calculation (0-60 months)';
COMMENT ON COLUMN who_standards.l_value IS 'Box-Cox transformation parameter (L)';
COMMENT ON COLUMN who_standards.m_value IS 'Median value (M)';
COMMENT ON COLUMN who_standards.s_value IS 'Coefficient of variation (S)';
