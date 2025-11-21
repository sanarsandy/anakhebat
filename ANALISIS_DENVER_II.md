# Analisis Implementasi Grafik Denver II

## 1. Pengenalan Denver II

Denver II adalah alat skrining perkembangan yang mengevaluasi 4 domain utama:
- **Personal-Social (PS)**: Interaksi sosial, kemandirian
- **Fine Motor-Adaptive (FM)**: Motorik halus, koordinasi tangan-mata
- **Language (L)**: Komunikasi, bahasa
- **Gross Motor (GM)**: Motorik kasar, gerakan tubuh

Grafik Denver II menampilkan:
- **Sumbu X**: Usia anak (bulan)
- **Sumbu Y**: Domain perkembangan (PS, FM, L, GM)
- **Plot**: Milestone yang sudah dicapai (pass) vs belum dicapai (fail) pada setiap usia

## 2. Data yang Tersedia Saat Ini

### 2.1 Struktur Database Milestones
- **Total milestones**: 53 items
- **Kategori saat ini**: 
  - `sensory` (3 items)
  - `motor` (28 items)
  - `perception` (4 items)
  - `cognitive` (18 items)
- **Source**: Semua menggunakan `KPSP` (tidak ada `DENVER`)
- **Age groups**: 9 kelompok usia (3, 6, 9, 12, 18, 24, 36, 48, 60 bulan)

### 2.2 Struktur Database Assessments
- **Status**: `yes`, `no`, `sometimes`
- **Assessment date**: Tanggal penilaian
- **Relationship**: One-to-one dengan milestones

### 2.3 Keterbatasan Data Saat Ini
1. **Kategori tidak sesuai Denver II**:
   - Saat ini: Learning Pyramid (sensory, motor, perception, cognitive)
   - Denver II: PS, FM, L, GM
   - **TIDAK ADA mapping langsung**

2. **Tidak ada data Denver II**:
   - Semua milestones menggunakan source `KPSP`
   - Tidak ada milestones dengan source `DENVER`

3. **Kategori motor terlalu umum**:
   - Denver II memisahkan Fine Motor (FM) dan Gross Motor (GM)
   - Data saat ini hanya punya `motor` (gabungan)

## 3. Analisis Kebutuhan untuk Grafik Denver II

### 3.1 Data yang Diperlukan

#### A. Milestone Data dengan Domain Denver II
**Wajib:**
- Kolom `denver_domain` atau mapping ke domain Denver II:
  - `PS` (Personal-Social)
  - `FM` (Fine Motor-Adaptive)
  - `L` (Language)
  - `GM` (Gross Motor)

**Opsi Implementasi:**
1. **Opsi 1: Tambah kolom `denver_domain`** (Recommended)
   - Tambah kolom baru di tabel `milestones`
   - Mapping existing milestones ke domain Denver II
   - Atau import data Denver II baru

2. **Opsi 2: Mapping berdasarkan kategori existing**
   - `sensory` → bisa ke `PS` atau `FM`
   - `motor` → perlu split ke `FM` dan `GM`
   - `perception` → bisa ke `FM` atau `L`
   - `cognitive` → bisa ke `L` atau `PS`
   - **Masalah**: Mapping tidak akurat, bisa misleading

3. **Opsi 3: Import data Denver II lengkap**
   - Import milestones dengan source `DENVER`
   - Setiap milestone sudah punya domain Denver II yang jelas
   - **Recommended untuk akurasi**

#### B. Assessment Data
**Sudah tersedia:**
- ✅ Status assessment (yes/no/sometimes)
- ✅ Assessment date
- ✅ Relationship dengan milestones

**Perlu ditambahkan:**
- ❌ Mapping ke domain Denver II (jika menggunakan opsi 1 atau 3)

### 3.2 Format Grafik Denver II

Grafik Denver II biasanya menampilkan:
```
Domain: PS | FM | L | GM
Age:    0  3  6  9  12 18 24 36 48 60
        |  |  |  |  |  |  |  |  |  |
        ✓  ✓  ✗  ✓  ✓  |  |  |  |  |  (PS milestones)
        |  ✓  ✓  ✓  ✗  ✓  |  |  |  |  (FM milestones)
        |  |  ✓  ✓  ✓  ✓  ✓  |  |  |  (L milestones)
        ✓  ✓  ✓  ✓  ✓  ✓  ✓  ✓  |  |  (GM milestones)
```

**Visualisasi:**
- **Bar chart** atau **line chart** per domain
- **X-axis**: Usia (bulan)
- **Y-axis**: Persentase milestone yang dicapai per domain
- **Color coding**: 
  - Hijau: Semua milestone tercapai
  - Kuning: Beberapa milestone belum tercapai
  - Merah: Banyak milestone belum tercapai (red flag)

## 4. Rekomendasi Implementasi

### 4.1 Opsi Terbaik: Hybrid Approach

**Langkah 1: Tambah Kolom `denver_domain`**
```sql
ALTER TABLE milestones 
ADD COLUMN denver_domain VARCHAR(10); -- PS, FM, L, GM, atau NULL
```

**Langkah 2: Import Data Denver II**
- Import milestones dengan source `DENVER`
- Setiap milestone sudah punya `denver_domain` yang jelas
- Minimal 4-6 milestones per domain per age group

**Langkah 3: Mapping Existing Data (Optional)**
- Mapping milestones KPSP yang relevan ke domain Denver II
- Hanya untuk milestone yang jelas mapping-nya
- Jangan paksa mapping jika tidak jelas

### 4.2 Struktur Data yang Diperlukan

#### Milestones dengan Denver Domain
```json
{
  "age_months": 6,
  "denver_domain": "PS",
  "question": "Smiles spontaneously at people",
  "source": "DENVER",
  "is_red_flag": false
}
```

#### Query untuk Grafik
```sql
SELECT 
  m.denver_domain,
  m.age_months,
  COUNT(*) as total_milestones,
  COUNT(CASE WHEN a.status = 'yes' THEN 1 END) as passed,
  COUNT(CASE WHEN a.status = 'no' THEN 1 END) as failed
FROM milestones m
LEFT JOIN assessments a ON m.id = a.milestone_id AND a.child_id = $1
WHERE m.denver_domain IS NOT NULL
GROUP BY m.denver_domain, m.age_months
ORDER BY m.denver_domain, m.age_months
```

## 5. Kesimpulan

### 5.1 Data Saat Ini: **TIDAK CUKUP**

**Alasan:**
1. ❌ Tidak ada kolom `denver_domain`
2. ❌ Tidak ada data dengan source `DENVER`
3. ❌ Kategori existing tidak sesuai dengan domain Denver II
4. ❌ Mapping manual tidak akurat

### 5.2 Yang Perlu Ditambahkan

**Wajib:**
1. ✅ Tambah kolom `denver_domain` di tabel `milestones`
2. ✅ Import data Denver II (minimal 20-30 milestones per domain)
3. ✅ Update migration untuk kolom baru
4. ✅ Update seed function untuk data Denver II

**Opsional (untuk akurasi lebih baik):**
1. Mapping milestones KPSP yang relevan ke domain Denver II
2. Tambah kolom `denver_item_number` untuk referensi standar Denver II

### 5.3 Estimasi Data yang Diperlukan

**Minimum untuk grafik yang meaningful:**
- **PS domain**: 20-25 milestones (usia 0-60 bulan)
- **FM domain**: 20-25 milestones
- **L domain**: 20-25 milestones
- **GM domain**: 20-25 milestones
- **Total**: 80-100 milestones dengan `denver_domain`

**Ideal:**
- 4-6 milestones per domain per age group (3, 6, 9, 12, 18, 24, 36, 48, 60 bulan)
- Total: ~180-240 milestones

## 6. Next Steps

1. **Buat migration** untuk tambah kolom `denver_domain`
2. **Siapkan data Denver II** dalam format JSON
3. **Update seed function** untuk import data Denver II
4. **Buat API endpoint** untuk fetch data Denver II
5. **Buat komponen grafik** Denver II di frontend
6. **Implementasi visualisasi** dengan Chart.js

## 7. Referensi

- Denver II Manual: Standard milestones untuk setiap domain
- WHO Developmental Milestones: Referensi tambahan
- CDC Milestones: Cross-reference untuk validasi

