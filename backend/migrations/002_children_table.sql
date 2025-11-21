-- Create children table
CREATE TABLE IF NOT EXISTS children (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    dob DATE NOT NULL,
    gender VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female')),
    birth_weight DECIMAL(5,2) NOT NULL,
    birth_height DECIMAL(5,2) NOT NULL,
    is_premature BOOLEAN DEFAULT FALSE,
    gestational_age INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster queries
CREATE INDEX idx_children_parent_id ON children(parent_id);
