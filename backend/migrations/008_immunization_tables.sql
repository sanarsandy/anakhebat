-- Create immunization_schedule table (Master Data)
CREATE TABLE IF NOT EXISTS immunization_schedule (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL, -- Nama imunisasi, e.g., "DPT-HB-Hib"
    name_id VARCHAR(255), -- Nama dalam Bahasa Indonesia
    description TEXT, -- Deskripsi singkat
    
    -- Timing Information (in days)
    age_min_days INT, -- Usia minimum dalam hari (0 = saat lahir)
    age_optimal_days INT, -- Usia optimal dalam hari
    age_max_days INT, -- Usia maksimum dalam hari (untuk catch-up)
    
    -- Timing Information (for display in months)
    age_min_months INT, -- Usia minimum dalam bulan (untuk display)
    age_optimal_months INT, -- Usia optimal dalam bulan (untuk display)
    age_max_months INT, -- Usia maksimum dalam bulan (untuk display)
    
    -- Dose Information
    dose_number INT NOT NULL, -- Dosis ke berapa (1, 2, 3, booster)
    total_doses INT, -- Total dosis yang diperlukan (misal: DPT perlu 3 dosis primer)
    
    -- Interval Information (untuk multi-dose)
    interval_from_previous_days INT, -- Interval minimum dari dosis sebelumnya (hari)
    interval_from_previous_months INT, -- Interval minimum (bulan)
    
    -- Category & Priority
    category VARCHAR(50) DEFAULT 'wajib', -- 'wajib', 'tambahan', 'catch-up'
    priority VARCHAR(20) DEFAULT 'medium', -- 'high', 'medium', 'low'
    is_required BOOLEAN DEFAULT TRUE, -- Apakah imunisasi wajib
    
    -- Additional Info
    notes TEXT, -- Catatan khusus, kontraindikasi, dll
    source VARCHAR(100) DEFAULT 'IDAI', -- Sumber rekomendasi
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create child_immunizations table (User Data)
CREATE TABLE IF NOT EXISTS child_immunizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    immunization_schedule_id UUID NOT NULL REFERENCES immunization_schedule(id),
    
    -- Pemberian Imunisasi
    given_date DATE NOT NULL, -- Tanggal pemberian
    given_at_age_days INT, -- Usia saat diberikan (dalam hari)
    given_at_age_months INT, -- Usia saat diberikan (dalam bulan, untuk display)
    
    -- Detail Pemberian
    location VARCHAR(255), -- Lokasi/tempat pemberian
    healthcare_facility VARCHAR(255), -- Fasilitas kesehatan
    doctor_name VARCHAR(255), -- Nama dokter/pemberi
    vaccine_batch_number VARCHAR(100), -- Nomor batch vaksin
    notes TEXT, -- Catatan tambahan
    
    -- Status
    is_on_schedule BOOLEAN, -- Apakah diberikan sesuai jadwal
    is_catch_up BOOLEAN DEFAULT FALSE, -- Apakah ini catch-up
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Prevent duplicate: one child can't have the same dose twice
CREATE UNIQUE INDEX IF NOT EXISTS idx_child_immunizations_unique 
ON child_immunizations(child_id, immunization_schedule_id);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_immunization_schedule_age ON immunization_schedule(age_min_days, age_max_days);
CREATE INDEX IF NOT EXISTS idx_immunization_schedule_category ON immunization_schedule(category);
CREATE INDEX IF NOT EXISTS idx_immunization_schedule_active ON immunization_schedule(is_active);
CREATE INDEX IF NOT EXISTS idx_child_immunizations_child ON child_immunizations(child_id);
CREATE INDEX IF NOT EXISTS idx_child_immunizations_date ON child_immunizations(given_date);
CREATE INDEX IF NOT EXISTS idx_child_immunizations_schedule ON child_immunizations(immunization_schedule_id);

-- Add comment
COMMENT ON TABLE immunization_schedule IS 'Master data jadwal imunisasi berdasarkan rekomendasi IDAI';
COMMENT ON TABLE child_immunizations IS 'Riwayat imunisasi yang sudah diberikan untuk setiap anak';

