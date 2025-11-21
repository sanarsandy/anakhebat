# Test Cases untuk useCorrectedAge Composable

## Test Scenarios

### 1. Anak Normal (Non-Premature)
**Input:**
- dob: "2024-01-01"
- isPremature: false
- gestationalAgeWeeks: null
- measurementDate: "2024-06-01" (5 bulan)

**Expected:**
- chronologicalMonths: 5
- correctedMonths: null
- useCorrected: false
- chronologicalDisplay: "5 bulan"
- correctedDisplay: null

---

### 2. Bayi Prematur < 24 Bulan - Menggunakan Corrected Age
**Input:**
- dob: "2024-01-01" (prematur 34 minggu)
- isPremature: true
- gestationalAgeWeeks: 34
- measurementDate: "2024-06-01" (5 bulan chronological)

**Calculation:**
- Weeks premature = 40 - 34 = 6 minggu = 42 hari
- Corrected days = 152 days (chronological) - 42 days = 110 days
- Corrected months ≈ 3.6 bulan ≈ 3 bulan

**Expected:**
- chronologicalMonths: 5
- correctedMonths: ~3
- useCorrected: true (karena < 24 bulan)
- chronologicalDisplay: "5 bulan"
- correctedDisplay: "3 bulan"

---

### 3. Bayi Prematur >= 24 Bulan - Tidak Menggunakan Corrected Age
**Input:**
- dob: "2022-01-01" (prematur 32 minggu)
- isPremature: true
- gestationalAgeWeeks: 32
- measurementDate: "2024-06-01" (29 bulan chronological)

**Expected:**
- chronologicalMonths: 29
- correctedMonths: 29 (sama dengan chronological)
- useCorrected: false (karena >= 24 bulan)
- chronologicalDisplay: "2 tahun 5 bulan"
- correctedDisplay: null

---

### 4. Bayi Prematur Tepat 24 Bulan
**Input:**
- dob: "2022-06-01" (prematur 35 minggu)
- isPremature: true
- gestationalAgeWeeks: 35
- measurementDate: "2024-06-01" (tepat 24 bulan)

**Expected:**
- chronologicalMonths: 24
- correctedMonths: 24
- useCorrected: false (boundary case - 730 days)
- chronologicalDisplay: "2 tahun 0 bulan"
- correctedDisplay: null

---

### 5. Edge Case: Gestational Age >= 40 weeks
**Input:**
- dob: "2024-01-01"
- isPremature: true
- gestationalAgeWeeks: 41 (seharusnya tidak prematur, tapi untuk test)
- measurementDate: "2024-06-01"

**Expected:**
- Weeks premature = 40 - 41 = -1 → 0 (safety check)
- useCorrected: false
- chronologicalMonths: 5

---

### 6. Edge Case: Negative Corrected Age
**Input:**
- dob: "2024-06-01" (sangat prematur 25 minggu)
- isPremature: true
- gestationalAgeWeeks: 25
- measurementDate: "2024-06-15" (hanya 14 hari setelah lahir)

**Calculation:**
- Weeks premature = 40 - 25 = 15 minggu = 105 hari
- Corrected days = 14 - 105 = -91 days → 0 (safety check)

**Expected:**
- chronologicalMonths: 0
- correctedMonths: 0 (tidak boleh negatif)
- useCorrected: true
- chronologicalDisplay: "0 bulan"
- correctedDisplay: "0 bulan"

---

### 7. Format Age Display - Tahun dan Bulan
**Input:**
- dob: "2022-01-01"
- isPremature: false
- measurementDate: "2024-06-01"

**Expected:**
- chronologicalDisplay: "2 tahun 5 bulan"

---

### 8. Format Age Display - Hanya Bulan
**Input:**
- dob: "2024-01-01"
- isPremature: false
- measurementDate: "2024-05-01"

**Expected:**
- chronologicalDisplay: "4 bulan"

---

## Integration Tests

### Dashboard Component
- ✅ Age info computed property ter-update saat selectedChild berubah
- ✅ Badge "Usia Koreksi" muncul untuk prematur < 24 bulan
- ✅ Badge tidak muncul untuk non-prematur atau >= 24 bulan
- ✅ Info prematur muncul dengan gestational age

### Child Selector Component
- ✅ Setiap child menampilkan age dengan indikasi corrected jika applicable
- ✅ Dropdown tetap readable dengan info tambahan

### Children List Page
- ✅ Setiap child card menampilkan age info
- ✅ Badge corrected age hanya untuk yang applicable
- ✅ Info prematur tersedia

### Growth Page
- ✅ Latest measurement menampilkan indikasi corrected age jika digunakan
- ✅ Conditional display berdasarkan child dan measurement date

