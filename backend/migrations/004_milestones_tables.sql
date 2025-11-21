-- Create milestones table (Master Data)
CREATE TABLE IF NOT EXISTS milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    age_months INT NOT NULL, -- Target age: 3, 6, 9, 12, etc.
    min_age_range INT, -- Start of age range (e.g., 4 for 4-6 months)
    max_age_range INT, -- End of age range (e.g., 6 for 4-6 months)
    category VARCHAR(50) NOT NULL, -- 'sensory', 'motor', 'perception', 'cognitive'
    question TEXT NOT NULL, -- The milestone checklist item
    question_en TEXT, -- English version (optional)
    source VARCHAR(50) DEFAULT 'KPSP', -- 'CDC', 'KPSP', 'DENVER'
    is_red_flag BOOLEAN DEFAULT FALSE, -- Critical milestone
    pyramid_level INT NOT NULL, -- 1=Sensory, 2=Motor, 3=Perception, 4=Cognitive
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create assessments table (User Data)
CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    milestone_id UUID NOT NULL REFERENCES milestones(id) ON DELETE CASCADE,
    assessment_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'yes', 'no', 'sometimes'
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(child_id, milestone_id) -- One assessment per milestone per child
);

-- Create indexes
CREATE INDEX idx_milestones_age ON milestones(age_months);
CREATE INDEX idx_milestones_range ON milestones(min_age_range, max_age_range);
CREATE INDEX idx_milestones_category ON milestones(category);
CREATE INDEX idx_assessments_child ON assessments(child_id);
CREATE INDEX idx_assessments_milestone ON assessments(milestone_id);
