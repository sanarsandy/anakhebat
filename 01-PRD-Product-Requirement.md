# Product Requirement Document (PRD)
**Project Name:** Aplikasi Pemantauan Tumbuh Kembang Anak (Tukem)
**Version:** 1.1
**Status:** Draft
**Date:** 20 November 2025

---

## 1. Executive Summary
Aplikasi SaaS berbasis web (*web-based*) yang komprehensif untuk memantau pertumbuhan fisik (antropometri) dan perkembangan mental/sensorik anak usia 0-5 tahun (0-60 bulan). Aplikasi ini bertujuan untuk mengisi celah pasar yang didominasi aplikasi *mobile* dengan menyediakan visualisasi data yang mendalam di layar besar (desktop) serta laporan siap cetak untuk keperluan medis.

## 2. User Personas
### Primary User: Orang Tua (The Anxious Parent)
- **Goal:** Memantau kesehatan dan perkembangan anak secara rutin dan akurat.
- **Pain Point:** Cemas dengan informasi yang simpang siur, butuh kepastian data medis (bukan mitos), ingin laporan yang bisa dibawa ke dokter.
- **Behavior:** Rajin mencatat, suka visualisasi data, butuh panduan stimulasi jika ada masalah.

### Secondary User: Dokter/Terapis (Future Dev)
- **Goal:** Menerima laporan hasil pantauan historis yang valid.
- **Pain Point:** Data anamnesa dari orang tua sering tidak lengkap atau tidak terstruktur.

### Admin System
- **Goal:** Mengelola data master (Standar WHO, Milestones) dan user.

## 3. Functional Requirements
Prioritas: **P1 (High)**, **P2 (Medium)**, **P3 (Low)**

### 3.1 Modul Manajemen Akun & Profil Anak
| ID | Requirement | Priority |
| :--- | :--- | :--- |
| **FR-01** | Register & Login (Email/Password & Google OAuth). | P1 |
| **FR-02** | Manajemen Data Anak (Nama, Tgl Lahir, Gender, BB/TB Lahir). | P1 |
| **FR-03** | **Usia Koreksi:** Fitur otomatis untuk anak prematur (Input HPL vs Tgl Lahir Asli). | P1 |
| **FR-04** | **Multi-profile:** Satu akun orang tua bisa mengelola banyak data anak. | P2 |

### 3.2 Modul Pertumbuhan Fisik (Anthropometry Engine)
| ID | Requirement | Priority |
| :--- | :--- | :--- |
| **FR-05** | Input Pengukuran: Tanggal, Berat (kg), Tinggi (cm), Lingkar Kepala (cm). | P1 |
| **FR-06** | **Age Calculation:** Hitung usia presisi (Tahun, Bulan, Hari). | P1 |
| **FR-07** | **Z-Score Calculation:** Otomatis hitung Z-Score (LMS WHO) untuk BB/U, TB/U, BB/TB. | P1 |
| **FR-08** | **Interpretation:** Status gizi otomatis (e.g., "Gizi Buruk" jika Z < -3). | P1 |
| **FR-09** | **Growth Chart:** Grafik interaktif (Chart.js) dengan garis standar SD -3 s.d +3. | P1 |

### 3.3 Modul Perkembangan (Development Tracker)
| ID | Requirement | Priority |
| :--- | :--- | :--- |
| **FR-10** | Checklist Milestone dinamis sesuai usia (CDC/KPSP). | P1 |
| **FR-11** | Kategori berdasarkan **Piramida Belajar**: Sensorik, Motorik, Persepsi, Kognitif. | P2 |
| **FR-12** | Status Capaian: "Ya", "Tidak", "Kadang-kadang". | P1 |
| **FR-13** | **Logic Warning:** Peringatan jika loncat level (e.g., Kognitif ok tapi Sensorik belum). | P2 |
| **FR-14** | **Red Flags:** Deteksi tanda bahaya kritis (e.g., 18 bulan belum kontak mata). | P1 |

### 3.4 Modul Intervensi & Rekomendasi
| ID | Requirement | Priority |
| :--- | :--- | :--- |
| **FR-15** | Rekomendasi Stimulasi (Video/Artikel) jika status "Tidak/Belum". | P2 |
| **FR-16** | Jadwal Imunisasi (IDAI). | P3 |

### 3.5 Modul Pelaporan (Output)
| ID | Requirement | Priority |
| :--- | :--- | :--- |
| **FR-17** | **PDF Export:** Laporan profesional untuk Dokter (Grafik + Milestone Log). | P1 |
| **FR-18** | Dashboard Summary: Ringkasan status kesehatan anak secara real-time. | P1 |

## 4. Non-Functional Requirements
- **NFR-01 Performance:** Dashboard load < 2s. Z-Score calc < 200ms.
- **NFR-02 Security:** Password Hash (Bcrypt), HTTPS, JWT Auth.
- **NFR-03 Reliability:** Data backup & consistency.
- **NFR-04 UI/UX:** Responsive Design (Mobile/Desktop), Color-coded indicators (Green/Red).

## 5. Success Metrics (MVP)
1. User dapat register dan input data anak.
2. Grafik pertumbuhan terbentuk sesuai standar WHO.
3. Sistem mendeteksi status gizi (Normal/Malnutrisi).
4. Export PDF berfungsi dengan baik.