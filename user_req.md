# USER REQUIREMENT DOCUMENT (URD)
**Nama Proyek:** Aplikasi Pemantauan Tumbuh Kembang Anak (Web-Based SaaS)  
**Versi:** 1.0  
**Status:** Draft  
**Tanggal:** 20 November 2025  

---

## 1. Pendahuluan

### 1.1 Tujuan Proyek
Membangun aplikasi berbasis web (*web-based*) yang komprehensif untuk memantau pertumbuhan fisik (antropometri) dan perkembangan mental/sensorik anak usia 0-5 tahun. Aplikasi ini bertujuan untuk mengisi celah pasar yang didominasi aplikasi *mobile* dengan menyediakan visualisasi data yang mendalam di layar besar (desktop) serta laporan siap cetak untuk keperluan medis.

### 1.2 Lingkup Masalah (Scope)
1.  **Pencatatan Fisik:** Menggunakan standar **WHO Child Growth Standards** dengan metode kalkulasi Z-Score (LMS).
2.  **Pencatatan Perkembangan:** Menggunakan landasan teori **Piramida Belajar (Williams & Shellenberger)** dan **CDC Milestones**.
3.  **Intervensi Dini:** Memberikan saran stimulasi otomatis berbasis video/artikel jika ditemukan keterlambatan.
4.  **Pelaporan:** Menghasilkan laporan PDF profesional untuk rujukan ke Dokter Spesialis Anak (DSA).

### 1.3 Target Pengguna
* **Orang Tua (Primary User):** Memantau anak, input data rutin.
* **Dokter/Terapis (Secondary User - Future Dev):** Menerima laporan hasil pantauan.
* **Administrator Sistem:** Mengelola data master (standar WHO, database pertanyaan).

---

## 2. Kebutuhan Fungsional (Functional Requirements)

Bagian ini menjelaskan fitur-fitur yang **harus** ada dalam sistem. Prioritas: **P1 (High)**, **P2 (Medium)**, **P3 (Low)**.

### 2.1 Modul Manajemen Akun & Profil Anak
| ID | Deskripsi Kebutuhan | Prioritas |
| :--- | :--- | :--- |
| **FR-01** | Sistem harus memungkinkan user mendaftar dan login (Email/Password & Google OAuth). | P1 |
| **FR-02** | User dapat menambahkan data anak (Nama, Tgl Lahir, Gender, BB Lahir, TB Lahir). | P1 |
| **FR-03** | Sistem harus memiliki fitur **Koreksi Usia** untuk anak prematur (Input HPL vs Tgl Lahir Asli). | P1 |
| **FR-04** | User dapat mengelola lebih dari satu profil anak dalam satu akun (Multi-profile). | P2 |

### 2.2 Modul Pertumbuhan Fisik (Anthropometry Engine)
| ID | Deskripsi Kebutuhan | Prioritas |
| :--- | :--- | :--- |
| **FR-05** | User dapat menginput data pengukuran berkala: Tanggal Ukur, Berat Badan (kg), Tinggi Badan (cm), Lingkar Kepala (cm). | P1 |
| **FR-06** | **[Logic]** Sistem harus menghitung usia presisi (Tahun, Bulan, Hari) saat data diinput. | P1 |
| **FR-07** | **[Logic]** Sistem wajib menghitung **Z-Score** otomatis menggunakan rumus LMS WHO untuk 3 indikator: BB/U (Berat/Umur), TB/U (Tinggi/Umur), BB/TB (Berat/Tinggi). | P1 |
| **FR-08** | Sistem harus menampilkan interpretasi status gizi secara otomatis (Misal: Z-Score < -3 = "Gizi Buruk"). | P1 |
| **FR-09** | **[UI]** Sistem menampilkan **Grafik Pertumbuhan** interaktif (Chart.js) yang membandingkan data anak dengan kurva standar WHO (Garis SD -3 s.d +3). | P1 |

### 2.3 Modul Perkembangan & Piramida Belajar
| ID | Deskripsi Kebutuhan | Prioritas |
| :--- | :--- | :--- |
| **FR-10** | Sistem menyediakan *checklist* perkembangan yang disesuaikan dengan usia anak (Berdasarkan CDC/KPSP). | P1 |
| **FR-11** | Checklist harus dikategorikan berdasarkan **Level Piramida Belajar**: Sensorik (Dasar), Motorik, Persepsi, Kognitif (Puncak). | P2 |
| **FR-12** | User dapat memilih status capaian: "Ya / Bisa", "Tidak / Belum", "Kadang-kadang". | P1 |
| **FR-13** | **[Logic]** Sistem memberikan peringatan (*Warning*) jika user fokus pada level Kognitif padahal level Sensorik belum tuntas. | P2 |
| **FR-14** | Sistem mendeteksi **Red Flags** (Tanda Bahaya) jika ada milestone kritis yang terlewat (Misal: 18 bulan belum ada kontak mata). | P1 |

### 2.4 Modul Intervensi & Rekomendasi (SaaS Value)
| ID | Deskripsi Kebutuhan | Prioritas |
| :--- | :--- | :--- |
| **FR-15** | Jika status capaian "Tidak/Belum", sistem otomatis merekomendasikan konten stimulasi spesifik (Video/Artikel). | P2 |
| **FR-16** | Sistem memberikan jadwal imunisasi rekomendasi IDAI berdasarkan tanggal lahir. | P3 |

### 2.5 Modul Pelaporan (Output)
| ID | Deskripsi Kebutuhan | Prioritas |
| :--- | :--- | :--- |
| **FR-17** | Sistem dapat men-generate laporan PDF ("Doctor's Report") yang berisi ringkasan grafik dan list milestone yang gagal. | P1 |
| **FR-18** | Dashboard utama harus menampilkan ringkasan status terakhir ("Anak Anda sehat, namun perlu stimulasi motorik kasar"). | P1 |

---

## 3. Kebutuhan Non-Fungsional (Non-Functional Requirements)

Bagian ini menjelaskan standar kualitas sistem.

### 3.1 Performa & Scalability
* **NFR-01:** Halaman dashboard utama harus dimuat dalam waktu < 2 detik (menggunakan Nuxt.js SSR/SPA).
* **NFR-02:** API Response time untuk kalkulasi Z-Score harus < 200ms (menggunakan Go/PHP Optimized).
* **NFR-03:** Sistem mampu menangani *concurrent users* (pengguna bersamaan) tanpa degradasi performa.

### 3.2 Keamanan (Security)
* **NFR-04:** Password user harus di-hash menggunakan algoritma kuat (Bcrypt/Argon2).
* **NFR-05:** Seluruh komunikasi data wajib menggunakan protokol HTTPS/SSL.
* **NFR-06:** API harus dilindungi dengan mekanisme Authentication (JWT - JSON Web Token).
* **NFR-07:** Perlindungan terhadap serangan umum web (SQL Injection, XSS, CSRF).

### 3.3 Usability & UI/UX
* **NFR-08:** Desain antarmuka harus *Responsive* (nyaman dibuka di Laptop, Tablet, maupun HP).
* **NFR-09:** Visualisasi grafik harus mudah dibaca oleh orang awam (menggunakan kode warna: Hijau=Aman, Merah=Bahaya).

---

## 4. Arsitektur & Batasan Sistem (System Constraints)

* **Frontend Framework:** Nuxt.js (Vue 3)
* **Backend Language:** Go (Golang) atau PHP (Laravel/Hyperf) - *To be decided by User*.
* **Database:** MariaDB / MySQL / PostgreSQL.
* **Charting Library:** Chart.js atau ApexCharts.
* **PDF Engine:** TCPDF / DomPDF / Puppeteer.
* **Deployment:** Docker Container di Linux Server.

---

## 5. Model Data (Entitas Utama)

Gambaran kasar entitas yang akan dikelola sistem:

1.  **Master Data WHO:** Menyimpan tabel *L, M, S* standar baku global.
2.  **Master Milestones:** Menyimpan bank pertanyaan checklist & kategori piramida belajar.
3.  **Users:** Data akun orang tua.
4.  **Children:** Profil anak.
5.  **Measurements:** Log data fisik (BB, TB, LK) per tanggal.
6.  **Assessments:** Log jawaban checklist perkembangan per tanggal.

---

## 6. Kriteria Keberhasilan (Acceptance Criteria)

Proyek dianggap berhasil tahap MVP (*Minimum Viable Product*) jika:
1.  User bisa login dan input data anak.
2.  Grafik pertumbuhan terbentuk sesuai standar WHO (bisa dibandingkan dengan kurva manual WHO).
3.  Sistem bisa mengeluarkan peringatan jika input data berada di area "Gizi Buruk" atau "Stunting".
4.  Fungsi Export PDF berjalan dan menghasilkan dokumen yang rapi.