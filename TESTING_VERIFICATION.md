# Testing & Verification - Corrected Age UI Indication

## Konsistensi Frontend vs Backend

### Perhitungan Corrected Age

#### Backend (Go):
```go
// weeksPremature = 40 - gestationalAgeWeeks
daysToSubtract = weeksPremature * 7
correctedDays = chronoDays - daysToSubtract
correctedMonths = int(float64(correctedDays) / 30.44)
```

#### Frontend (TypeScript):
```typescript
// weeksPremature = 40 - gestationalAgeWeeks
daysToSubtract = weeksPremature * 7
correctedDays = daysDiff - daysToSubtract
correctedMonths = Math.max(0, Math.floor(correctedDays / 30.44))
```

✅ **Konsisten**: Logika sama, menggunakan 30.44 days per month

---

## Test Cases Manual

### Test Case 1: Anak Normal
**Setup:**
- DOB: 2024-01-01
- isPremature: false
- Current Date: 2024-06-01 (5 bulan)

**Expected UI:**
- ✅ Dashboard: Menampilkan "5 bulan" tanpa badge corrected age
- ✅ Child Selector: Menampilkan "5 bulan"
- ✅ Children List: Menampilkan "5 bulan"

---

### Test Case 2: Bayi Prematur 34 minggu, 5 bulan chronological
**Setup:**
- DOB: 2024-01-01
- isPremature: true
- gestationalAge: 34 minggu
- Current Date: 2024-06-01

**Calculation:**
- Chronological: 152 days = ~5 bulan
- Weeks premature: 40 - 34 = 6 minggu = 42 days
- Corrected: 152 - 42 = 110 days = ~3.6 bulan ≈ 3 bulan
- useCorrected: true (karena 152 < 730)

**Expected UI:**
- ✅ Dashboard: 
  - Menampilkan "5 bulan"
  - Badge amber "Usia Koreksi: 3 bulan"
  - Text kecil: "Prematur (34 minggu) - Menggunakan usia koreksi"
- ✅ Child Selector:
  - Menampilkan "5 bulan (Usia Koreksi: 3 bulan)"
  - Badge "Koreksi: 3 bulan"
- ✅ Children List:
  - Menampilkan "5 bulan"
  - Badge amber "Usia Koreksi: 3 bulan"
  - Info "Prematur (34 minggu)"

---

### Test Case 3: Bayi Prematur 32 minggu, 25 bulan chronological
**Setup:**
- DOB: 2022-06-01
- isPremature: true
- gestationalAge: 32 minggu
- Current Date: 2024-07-01 (25 bulan = 761 days)

**Calculation:**
- Chronological: 761 days = 25 bulan
- Days >= 730, so useCorrected = false

**Expected UI:**
- ✅ Dashboard:
  - Menampilkan "2 tahun 1 bulan"
  - ❌ TIDAK ada badge corrected age
  - Info prematur tetap muncul tapi tanpa "Menggunakan usia koreksi"
- ✅ Child Selector:
  - Menampilkan "2 tahun 1 bulan"
  - ❌ TIDAK ada badge corrected
- ✅ Children List:
  - Menampilkan "2 tahun 1 bulan"
  - ❌ TIDAK ada badge corrected age
  - Info "Prematur (32 minggu)" tetap muncul

---

### Test Case 4: Growth Page - Latest Measurement
**Setup:**
- Child: Prematur 34 minggu, 5 bulan
- Latest Measurement: 2024-06-01

**Expected UI:**
- ✅ Growth Page:
  - Menampilkan age_display dari measurement (sudah processed by backend)
  - Badge "Usia Koreksi digunakan" muncul

---

## Edge Cases Testing

### Edge Case 1: Gestational Age >= 40
**Setup:**
- isPremature: true (seharusnya false, tapi untuk test)
- gestationalAge: 41 minggu

**Expected:**
- ✅ weeksPremature = 40 - 41 = -1 → 0 (safety check)
- ✅ useCorrected = false
- ✅ Menampilkan chronological age only

---

### Edge Case 2: Negative Corrected Age
**Setup:**
- DOB: 2024-06-01
- isPremature: true
- gestationalAge: 25 minggu (sangat prematur)
- Measurement Date: 2024-06-15 (hanya 14 hari setelah lahir)

**Calculation:**
- Chronological: 14 days
- Weeks premature: 40 - 25 = 15 minggu = 105 days
- Corrected: 14 - 105 = -91 → 0 (safety check)

**Expected:**
- ✅ correctedMonths = 0 (tidak boleh negatif)
- ✅ Menampilkan "0 bulan" untuk corrected age

---

### Edge Case 3: Boundary - Tepat 24 Bulan (730 days)
**Setup:**
- DOB: 2022-06-01
- isPremature: true
- gestationalAge: 35 minggu
- Measurement Date: 2024-06-01 (tepat 730 days)

**Expected:**
- ✅ daysDiff = 730
- ✅ useCorrected = false (boundary case: >= 730)
- ✅ Menampilkan chronological age only

---

### Edge Case 4: Boundary - 729 days (< 24 bulan)
**Setup:**
- DOB: 2022-06-01
- isPremature: true
- gestationalAge: 35 minggu
- Measurement Date: 2024-05-31 (729 days)

**Expected:**
- ✅ daysDiff = 729
- ✅ useCorrected = true (karena < 730)
- ✅ Menampilkan corrected age

---

### Edge Case 5: Null/Undefined Gestational Age
**Setup:**
- isPremature: true
- gestationalAge: null atau undefined

**Expected:**
- ✅ useCorrected = false
- ✅ Menampilkan chronological age only

---

## Integration Testing Checklist

### Dashboard Component
- [ ] Age info computed reactive saat selectedChild berubah
- [ ] Badge corrected age muncul untuk prematur < 24 bulan
- [ ] Badge tidak muncul untuk non-prematur
- [ ] Badge tidak muncul untuk prematur >= 24 bulan
- [ ] Info prematur dengan gestational age muncul
- [ ] Layout tidak broken dengan badge tambahan

### Child Selector Component
- [ ] Dropdown tetap readable dengan info tambahan
- [ ] Badge tidak terlalu besar
- [ ] Setiap child menampilkan info yang benar

### Children List Page
- [ ] Card tidak terlalu besar dengan info tambahan
- [ ] Badge hanya muncul untuk yang applicable
- [ ] Info prematur tersedia

### Growth Page
- [ ] Latest measurement badge hanya muncul jika applicable
- [ ] Tidak muncul untuk non-prematur
- [ ] Tidak muncul untuk prematur >= 24 bulan

---

## Browser Testing

### Responsive Design
- [ ] Mobile: Badge masih readable
- [ ] Tablet: Layout tetap baik
- [ ] Desktop: Semua info terlihat jelas

### Browser Compatibility
- [ ] Chrome: Berfungsi
- [ ] Firefox: Berfungsi
- [ ] Safari: Berfungsi
- [ ] Edge: Berfungsi

---

## Performance Testing

- [ ] Computed property tidak menyebabkan re-render berlebihan
- [ ] Calculate corrected age tidak lambat
- [ ] Tidak ada lag saat switching child

---

## Accessibility

- [ ] Badge memiliki kontras warna yang cukup
- [ ] Text readable
- [ ] Info penting tidak hilang di mobile

