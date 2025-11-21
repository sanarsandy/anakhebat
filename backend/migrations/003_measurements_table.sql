-- Create measurements table
CREATE TABLE IF NOT EXISTS measurements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    measurement_date DATE NOT NULL,
    weight DECIMAL(5,2) NOT NULL,
    height DECIMAL(6,2) NOT NULL,
    head_circumference DECIMAL(5,2),
    age_in_days INT NOT NULL,
    age_in_months INT NOT NULL,
    weight_for_age_zscore DECIMAL(5,2),
    height_for_age_zscore DECIMAL(5,2),
    weight_status VARCHAR(50),
    height_status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster queries
CREATE INDEX idx_measurements_child_id ON measurements(child_id);
CREATE INDEX idx_measurements_date ON measurements(measurement_date);
