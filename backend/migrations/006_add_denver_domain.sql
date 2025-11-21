-- Migration: Add Denver II domain column to milestones table
-- This allows milestones to be categorized by Denver II domains: PS, FM, L, GM

ALTER TABLE milestones 
ADD COLUMN IF NOT EXISTS denver_domain VARCHAR(10); -- PS, FM, L, GM, or NULL

-- Add index for faster queries by Denver domain
CREATE INDEX IF NOT EXISTS idx_milestones_denver_domain ON milestones(denver_domain) WHERE denver_domain IS NOT NULL;

-- Add comment
COMMENT ON COLUMN milestones.denver_domain IS 'Denver II domain: PS (Personal-Social), FM (Fine Motor-Adaptive), L (Language), GM (Gross Motor)';

