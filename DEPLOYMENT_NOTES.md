# Deployment Notes - WhatsApp OTP Authentication

## Migration Required

Sebelum deploy ke production, pastikan untuk menjalankan migration database berikut:

### Migration File
- `backend/migrations/009_add_phone_otp_auth.sql`

### Cara Menjalankan Migration

#### Option 1: Menggunakan script migration
```bash
cd backend
./scripts/run-migrations.sh
```

#### Option 2: Manual via psql
```bash
psql -U tukem_user -d tukem_db -f backend/migrations/009_add_phone_otp_auth.sql
```

#### Option 3: Via Docker
```bash
docker exec -i tukem-db psql -U tukem_user -d tukem_db < backend/migrations/009_add_phone_otp_auth.sql
```

### Yang Dilakukan Migration

1. **Menambahkan kolom ke tabel `users`:**
   - `phone_number` (VARCHAR(20), unique, nullable)
   - `phone_verified` (BOOLEAN, default false)
   - `phone_verified_at` (TIMESTAMP, nullable)
   - Membuat `email` menjadi optional (DROP NOT NULL)

2. **Membuat tabel baru:**
   - `otp_codes` - Menyimpan kode OTP untuk autentikasi
   - `otp_rate_limits` - Rate limiting untuk request OTP

3. **Membuat index:**
   - Index untuk `phone_number` di tabel `users`
   - Index untuk pencarian OTP yang efisien

## Environment Variables

Pastikan environment variables berikut sudah di-set di production:

```bash
# WhatsApp Gateway Configuration
WHATSAPP_GATEWAY_URL=https://anakhebat.web.id/services/wa-gateway/api/send-file
WHATSAPP_API_KEY=YOUR_SECURE_API_KEY
WHATSAPP_SENDER_NUMBER=6281234567890

# Google OAuth (jika digunakan)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=https://yourdomain.com/api/auth/google/callback
FRONTEND_URL=https://yourdomain.com

# JWT Secret
JWT_SECRET=your_secure_jwt_secret
```

## Perubahan Fitur

### 1. Registration
- **Sebelum:** Register menggunakan email + password
- **Sekarang:** Register menggunakan nomor WhatsApp + nama lengkap
- Tidak lagi memerlukan password

### 2. Login
- **WhatsApp OTP Login:**
  - User harus sudah terdaftar terlebih dahulu
  - Request OTP akan gagal jika user belum terdaftar
  - Error message: "User belum terdaftar. Silakan daftar terlebih dahulu."

- **Google OAuth Login:**
  - Tetap berfungsi seperti sebelumnya
  - User yang login via Google akan otomatis dibuat jika belum ada

### 3. Database Changes
- Tabel `users` sekarang mendukung:
  - Phone number authentication
  - Email menjadi optional (untuk backward compatibility dengan Google OAuth)
  - Multiple auth providers: 'phone', 'google', 'phone_google'

## Testing Checklist

Setelah deployment, test hal berikut:

1. ✅ Register dengan nomor WhatsApp baru
2. ✅ Login dengan nomor yang sudah terdaftar (OTP)
3. ✅ Login dengan nomor yang belum terdaftar (harus error)
4. ✅ Resend OTP
5. ✅ Google OAuth login (jika digunakan)
6. ✅ Rate limiting OTP (max 3 requests per 15 menit)

## Rollback Plan

Jika perlu rollback:

1. **Database:** Tidak ada rollback migration yang disediakan. Jika perlu, buat manual:
   ```sql
   -- Hapus tabel baru
   DROP TABLE IF EXISTS otp_codes;
   DROP TABLE IF EXISTS otp_rate_limits;
   
   -- Hapus kolom baru (HATI-HATI: pastikan tidak ada data penting)
   ALTER TABLE users DROP COLUMN IF EXISTS phone_number;
   ALTER TABLE users DROP COLUMN IF EXISTS phone_verified;
   ALTER TABLE users DROP COLUMN IF EXISTS phone_verified_at;
   
   -- Restore email NOT NULL (jika diperlukan)
   ALTER TABLE users ALTER COLUMN email SET NOT NULL;
   ```

2. **Code:** Revert commit atau checkout ke commit sebelumnya

## Notes

- Migration ini **tidak destructive** - tidak menghapus data existing
- User yang sudah ada tetap bisa login dengan email/password (jika menggunakan auth_provider 'email')
- Google OAuth tetap berfungsi normal
- Phone number disimpan dalam format internasional (+62...)

