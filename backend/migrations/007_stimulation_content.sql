-- Create stimulation_content table for intervention recommendations
CREATE TABLE IF NOT EXISTS stimulation_content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    milestone_id UUID REFERENCES milestones(id) ON DELETE SET NULL, -- Optional: link to specific milestone
    category VARCHAR(50) NOT NULL, -- 'sensory', 'motor', 'perception', 'cognitive'
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content_type VARCHAR(20) NOT NULL DEFAULT 'article', -- 'video' or 'article'
    url TEXT NOT NULL, -- URL to content (video URL or article link)
    thumbnail_url TEXT, -- Optional thumbnail image URL
    age_min_months INT, -- Minimum age for this content (optional)
    age_max_months INT, -- Maximum age for this content (optional)
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster queries
CREATE INDEX IF NOT EXISTS idx_stimulation_content_milestone ON stimulation_content(milestone_id);
CREATE INDEX IF NOT EXISTS idx_stimulation_content_category ON stimulation_content(category);
CREATE INDEX IF NOT EXISTS idx_stimulation_content_age_range ON stimulation_content(age_min_months, age_max_months);
CREATE INDEX IF NOT EXISTS idx_stimulation_content_active ON stimulation_content(is_active);

-- Add comment
COMMENT ON TABLE stimulation_content IS 'Stimulation content (videos/articles) for intervention recommendations based on milestone status';

