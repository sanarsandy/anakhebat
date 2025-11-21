# Analisis & Rencana Implementasi SAAS RBAC & Google OAuth
**Tanggal:** 21 November 2025  
**Status:** Rencana Implementasi Multi-Tenant SAAS dengan Role-Based Access Control

---

## 📊 Executive Summary

Aplikasi Tukem akan dikembangkan menjadi SAAS (Software as a Service) dengan:
1. **Role-Based Access Control (RBAC)** berdasarkan subscription plan (Free vs Paid)
2. **Google OAuth** untuk login alternative
3. **Multi-tenant architecture** dengan plan-based feature limitation

---

## 🎯 Tujuan & Requirement

### 1. User Plans & Tiers
- **Free Plan**: Akses terbatas untuk fitur-fitur dasar
- **Paid Plan**: Akses penuh ke semua fitur
- **Future Plans**: Pro, Enterprise (opsional untuk masa depan)

### 2. Google OAuth
- Login menggunakan akun Google
- Auto-create account jika belum terdaftar
- Link/unlink Google account dengan email existing

### 3. Feature Limitation per Plan
- Batasi fitur berdasarkan subscription plan
- Upgrade/downgrade plan
- Billing management (untuk masa depan)

---

## 📋 ANALISIS STRUKTUR SAAT INI

### Database Schema Saat Ini

**Users Table (Current Schema):**
```sql
- id (UUID, primary key)
- email (unique)
- password_hash (NOT NULL)
- full_name
- role (default 'parent')
- created_at
```

**Note:** `full_name` digunakan, bukan `name`

**Limitation:**
- ❌ Tidak ada field untuk subscription plan
- ❌ Tidak ada field untuk Google OAuth
- ❌ Tidak ada table untuk subscriptions
- ❌ Tidak ada table untuk plan features
- ❌ Tidak ada middleware untuk check plan limits

---

## 🏗️ RENCANA IMPLEMENTASI

### Phase 1: Database Schema Extension

#### 1.1. Update Users Table

**Migration: `009_update_users_for_saas.sql`**

```sql
-- Update users table schema: make password_hash nullable (for Google OAuth users)
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

-- Add columns to users table for SAAS features
ALTER TABLE users ADD COLUMN IF NOT EXISTS 
  google_id VARCHAR(255) UNIQUE,
  google_email VARCHAR(255),
  avatar_url TEXT,
  subscription_plan VARCHAR(50) DEFAULT 'free' CHECK (subscription_plan IN ('free', 'paid', 'pro', 'enterprise')),
  subscription_status VARCHAR(50) DEFAULT 'active' CHECK (subscription_status IN ('active', 'cancelled', 'expired', 'trial')),
  subscription_start_date TIMESTAMP,
  subscription_end_date TIMESTAMP,
  trial_ends_at TIMESTAMP,
  created_via VARCHAR(50) DEFAULT 'email' CHECK (created_via IN ('email', 'google')),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
CREATE INDEX IF NOT EXISTS idx_users_subscription_plan ON users(subscription_plan);
CREATE INDEX IF NOT EXISTS idx_users_subscription_status ON users(subscription_status);
```

**Changes:**
- ✅ `google_id`: ID dari Google OAuth
- ✅ `google_email`: Email dari Google account
- ✅ `avatar_url`: Profile picture dari Google
- ✅ `subscription_plan`: Plan saat ini (free, paid, pro, enterprise)
- ✅ `subscription_status`: Status subscription (active, cancelled, expired, trial)
- ✅ `subscription_start_date`: Tanggal mulai subscription
- ✅ `subscription_end_date`: Tanggal berakhir subscription
- ✅ `trial_ends_at`: Tanggal berakhir trial (untuk paid plan)
- ✅ `created_via`: Metode registrasi (email atau google)

---

#### 1.2. Create Subscription Plans Table

**Migration: `010_create_subscription_plans.sql`**

```sql
-- Subscription plans master data
CREATE TABLE IF NOT EXISTS subscription_plans (
  id VARCHAR(50) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  price_monthly DECIMAL(10, 2) DEFAULT 0,
  price_yearly DECIMAL(10, 2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'IDR',
  trial_days INTEGER DEFAULT 0,
  max_children INTEGER DEFAULT 1,
  max_measurements_per_month INTEGER DEFAULT NULL, -- NULL = unlimited
  max_assessments_per_month INTEGER DEFAULT NULL, -- NULL = unlimited
  features JSONB DEFAULT '{}', -- JSON object dengan feature flags
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default plans
INSERT INTO subscription_plans (id, name, description, price_monthly, price_yearly, trial_days, max_children, max_measurements_per_month, max_assessments_per_month, features) VALUES
('free', 'Free Plan', 'Akses terbatas untuk fitur dasar', 0, 0, 0, 1, 5, 2, '{"pdf_export": false, "unlimited_children": false, "growth_charts": true, "milestone_tracking": true, "denver_ii": false, "immunization_schedule": false, "recommendations": false}'),
('paid', 'Paid Plan', 'Akses penuh ke semua fitur', 99000, 990000, 7, 3, NULL, NULL, '{"pdf_export": true, "unlimited_children": false, "growth_charts": true, "milestone_tracking": true, "denver_ii": true, "immunization_schedule": true, "recommendations": true, "priority_support": false}'),
('pro', 'Pro Plan', 'Akses penuh + fitur premium', 149000, 1490000, 14, NULL, NULL, NULL, '{"pdf_export": true, "unlimited_children": true, "growth_charts": true, "milestone_tracking": true, "denver_ii": true, "immunization_schedule": true, "recommendations": true, "priority_support": true, "data_export": true, "api_access": false}')
ON CONFLICT (id) DO NOTHING;

-- Index
CREATE INDEX IF NOT EXISTS idx_subscription_plans_active ON subscription_plans(is_active);
```

**Plan Features:**
- `free`: 1 child, 5 measurements/month, 2 assessments/month, fitur dasar
- `paid`: 3 children, unlimited measurements/assessments, semua fitur
- `pro`: Unlimited children, semua fitur + premium features

---

#### 1.3. Create User Subscriptions Table (History)

**Migration: `011_create_user_subscriptions.sql`**

```sql
-- Subscription history untuk tracking changes
CREATE TABLE IF NOT EXISTS user_subscriptions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  plan_id VARCHAR(50) NOT NULL REFERENCES subscription_plans(id),
  status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'cancelled', 'expired', 'trial')),
  started_at TIMESTAMP NOT NULL,
  ended_at TIMESTAMP,
  cancelled_at TIMESTAMP,
  payment_provider VARCHAR(50), -- 'stripe', 'midtrans', 'manual', etc.
  payment_reference VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_status ON user_subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_plan_id ON user_subscriptions(plan_id);
```

**Purpose:**
- Track subscription history
- Audit trail untuk billing
- Support untuk refunds/cancellations

---

#### 1.4. Create Plan Feature Limits Table

**Migration: `012_create_plan_feature_limits.sql`**

```sql
-- Feature limits per plan (untuk flexible configuration)
CREATE TABLE IF NOT EXISTS plan_feature_limits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_id VARCHAR(50) NOT NULL REFERENCES subscription_plans(id) ON DELETE CASCADE,
  feature_key VARCHAR(100) NOT NULL,
  limit_value INTEGER, -- NULL = unlimited
  limit_type VARCHAR(50) NOT NULL, -- 'per_month', 'total', 'per_child'
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(plan_id, feature_key)
);

-- Insert feature limits for free plan
INSERT INTO plan_feature_limits (plan_id, feature_key, limit_value, limit_type, description) VALUES
('free', 'children', 1, 'total', 'Maximum number of children'),
('free', 'measurements', 5, 'per_month', 'Maximum measurements per month'),
('free', 'assessments', 2, 'per_month', 'Maximum assessments per month'),
('free', 'pdf_exports', 0, 'per_month', 'PDF exports per month (0 = disabled)')
ON CONFLICT (plan_id, feature_key) DO NOTHING;

-- Insert feature limits for paid plan
INSERT INTO plan_feature_limits (plan_id, feature_key, limit_value, limit_type, description) VALUES
('paid', 'children', 3, 'total', 'Maximum number of children'),
('paid', 'measurements', NULL, 'per_month', 'Unlimited measurements'),
('paid', 'assessments', NULL, 'per_month', 'Unlimited assessments'),
('paid', 'pdf_exports', NULL, 'per_month', 'Unlimited PDF exports')
ON CONFLICT (plan_id, feature_key) DO NOTHING;

-- Index
CREATE INDEX IF NOT EXISTS idx_plan_feature_limits_plan_id ON plan_feature_limits(plan_id);
```

**Purpose:**
- Flexible configuration untuk plan limits
- Easy to add/modify limits tanpa code changes
- Support untuk different limit types (per month, total, per child)

---

### Phase 2: Backend Implementation

#### 2.1. Update Models

**File: `backend/models/user.go`**

```go
type User struct {
    ID                  string    `json:"id" db:"id"`
    Email               string    `json:"email" db:"email"`
    PasswordHash        string    `json:"-" db:"password_hash"`
    Name                string    `json:"name" db:"name"`
    GoogleID            *string   `json:"google_id,omitempty" db:"google_id"`
    GoogleEmail         *string   `json:"google_email,omitempty" db:"google_email"`
    AvatarURL           *string   `json:"avatar_url,omitempty" db:"avatar_url"`
    SubscriptionPlan    string    `json:"subscription_plan" db:"subscription_plan"`
    SubscriptionStatus  string    `json:"subscription_status" db:"subscription_status"`
    SubscriptionStartDate *time.Time `json:"subscription_start_date,omitempty" db:"subscription_start_date"`
    SubscriptionEndDate *time.Time `json:"subscription_end_date,omitempty" db:"subscription_end_date"`
    TrialEndsAt         *time.Time `json:"trial_ends_at,omitempty" db:"trial_ends_at"`
    CreatedVia          string    `json:"created_via" db:"created_via"`
    CreatedAt           time.Time `json:"created_at" db:"created_at"`
    UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type SubscriptionPlan struct {
    ID                     string                 `json:"id" db:"id"`
    Name                   string                 `json:"name" db:"name"`
    Description            string                 `json:"description" db:"description"`
    PriceMonthly           float64                `json:"price_monthly" db:"price_monthly"`
    PriceYearly            float64                `json:"price_yearly" db:"price_yearly"`
    Currency               string                 `json:"currency" db:"currency"`
    TrialDays              int                    `json:"trial_days" db:"trial_days"`
    MaxChildren            *int                   `json:"max_children,omitempty" db:"max_children"`
    MaxMeasurementsPerMonth *int                  `json:"max_measurements_per_month,omitempty" db:"max_measurements_per_month"`
    MaxAssessmentsPerMonth *int                   `json:"max_assessments_per_month,omitempty" db:"max_assessments_per_month"`
    Features               map[string]interface{} `json:"features" db:"features"`
    IsActive               bool                   `json:"is_active" db:"is_active"`
    CreatedAt              time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt              time.Time              `json:"updated_at" db:"updated_at"`
}
```

---

#### 2.2. Google OAuth Handler

**File: `backend/handlers/oauth.go`** (NEW)

```go
package handlers

import (
    "encoding/json"
    "io"
    "net/http"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/jmoiron/sqlx"
    "github.com/labstack/echo/v4"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "tukem-backend/models"
    "tukem-backend/utils"
)

var (
    googleOauthConfig *oauth2.Config
)

func InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
    googleOauthConfig = &oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        RedirectURL:  redirectURL,
        Scopes:       []string{"openid", "profile", "email"},
        Endpoint:     google.Endpoint,
    }
}

// GoogleLogin initiates OAuth flow
func GoogleLogin(c echo.Context) error {
    if googleOauthConfig == nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "OAuth not configured"})
    }
    
    url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles OAuth callback
func GoogleCallback(db *sqlx.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        code := c.QueryParam("code")
        if code == "" {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "Code not provided"})
        }
        
        // Exchange code for token
        token, err := googleOauthConfig.Exchange(c.Request().Context(), code)
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to exchange token"})
        }
        
        // Get user info from Google
        client := googleOauthConfig.Client(c.Request().Context(), token)
        resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user info"})
        }
        defer resp.Body.Close()
        
        data, _ := io.ReadAll(resp.Body)
        var googleUser struct {
            ID      string `json:"id"`
            Email   string `json:"email"`
            Name    string `json:"name"`
            Picture string `json:"picture"`
        }
        json.Unmarshal(data, &googleUser)
        
        // Check if user exists by email or google_id
        var user models.User
        err = db.Get(&user, "SELECT * FROM users WHERE email = $1 OR google_id = $2", googleUser.Email, googleUser.ID)
        
        if err != nil {
            // Create new user
            userID := uuid.New().String()
            _, err = db.Exec(`
                INSERT INTO users (id, email, name, google_id, google_email, avatar_url, subscription_plan, created_via)
                VALUES ($1, $2, $3, $4, $5, $6, 'free', 'google')
            `, userID, googleUser.Email, googleUser.Name, googleUser.ID, googleUser.Email, googleUser.Picture)
            
            if err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
            }
            
            user.ID = userID
            user.Email = googleUser.Email
            user.FullName = googleUser.Name
            user.Role = "parent"
            user.GoogleID = &googleUser.ID
            user.GoogleEmail = &googleUser.Email
            user.AvatarURL = &googleUser.Picture
        } else {
            // Update existing user with Google info if needed
            if user.GoogleID == nil {
                db.Exec("UPDATE users SET google_id = $1, google_email = $2, avatar_url = $3 WHERE id = $4",
                    googleUser.ID, googleUser.Email, googleUser.Picture, user.ID)
            }
        }
        
        // Generate JWT token
        jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "user_id": user.ID,
            "email":   user.Email,
            "exp":     time.Now().Add(time.Hour * 72).Unix(),
        })
        
        tokenString, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
        
        // Redirect to frontend with token
        frontendURL := os.Getenv("FRONTEND_URL") + "/auth/callback?token=" + tokenString
        return c.Redirect(http.StatusTemporaryRedirect, frontendURL)
    }
}
```

---

#### 2.3. Plan Checker Utility

**File: `backend/utils/plan_checker.go`** (NEW)

```go
package utils

import (
    "database/sql"
    "fmt"
    "time"
    
    "github.com/jmoiron/sqlx"
)

// CheckFeatureAccess checks if user has access to a feature
func CheckFeatureAccess(db *sqlx.DB, userID string, featureKey string) (bool, error) {
    // Get user plan
    var planID string
    var subscriptionStatus string
    var subscriptionEndDate *time.Time
    
    err := db.QueryRow(`
        SELECT subscription_plan, subscription_status, subscription_end_date
        FROM users WHERE id = $1
    `, userID).Scan(&planID, &subscriptionStatus, &subscriptionEndDate)
    
    if err != nil {
        return false, err
    }
    
    // Check if subscription is active
    if subscriptionStatus != "active" {
        if subscriptionEndDate != nil && time.Now().After(*subscriptionEndDate) {
            return false, fmt.Errorf("subscription expired")
        }
        return false, fmt.Errorf("subscription not active")
    }
    
    // Get feature from plan
    var hasFeature bool
    err = db.QueryRow(`
        SELECT (features->$1)::boolean
        FROM subscription_plans
        WHERE id = $2
    `, featureKey, planID).Scan(&hasFeature)
    
    if err == sql.ErrNoRows {
        return false, nil
    }
    if err != nil {
        return false, err
    }
    
    return hasFeature, nil
}

// CheckLimit checks if user has reached a limit
func CheckLimit(db *sqlx.DB, userID string, limitKey string) (bool, int, error) {
    // Get user plan and limit
    var planID string
    var limitValue *int
    var limitType string
    
    err := db.QueryRow(`
        SELECT u.subscription_plan, pf.limit_value, pf.limit_type
        FROM users u
        JOIN plan_feature_limits pf ON pf.plan_id = u.subscription_plan
        WHERE u.id = $1 AND pf.feature_key = $2
    `, userID, limitKey).Scan(&planID, &limitValue, &limitType)
    
    if err == sql.ErrNoRows {
        // No limit defined = unlimited
        return true, 0, nil
    }
    if err != nil {
        return false, 0, err
    }
    
    if limitValue == nil {
        // Unlimited
        return true, 0, nil
    }
    
    // Get current usage based on limit_type
    var currentUsage int
    switch limitType {
    case "per_month":
        err = getMonthlyUsage(db, userID, limitKey, &currentUsage)
    case "total":
        err = getTotalUsage(db, userID, limitKey, &currentUsage)
    case "per_child":
        // Need child_id, handle separately
        return true, *limitValue, nil
    }
    
    if err != nil {
        return false, 0, err
    }
    
    return currentUsage < *limitValue, *limitValue - currentUsage, nil
}

func getMonthlyUsage(db *sqlx.DB, userID string, limitKey string, usage *int) error {
    now := time.Now()
    startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    
    switch limitKey {
    case "measurements":
        return db.Get(usage, `
            SELECT COUNT(*) FROM measurements m
            JOIN children c ON c.id = m.child_id
            WHERE c.parent_id = $1 AND m.created_at >= $2
        `, userID, startOfMonth)
    case "assessments":
        return db.Get(usage, `
            SELECT COUNT(DISTINCT assessment_date) FROM assessments a
            JOIN children c ON c.id = a.child_id
            WHERE c.parent_id = $1 AND a.created_at >= $2
        `, userID, startOfMonth)
    case "pdf_exports":
        // Track in separate table or log
        return db.Get(usage, `SELECT 0`) // Placeholder
    }
    
    return fmt.Errorf("unknown limit key: %s", limitKey)
}

func getTotalUsage(db *sqlx.DB, userID string, limitKey string, usage *int) error {
    switch limitKey {
    case "children":
        return db.Get(usage, `SELECT COUNT(*) FROM children WHERE parent_id = $1`, userID)
    }
    
    return fmt.Errorf("unknown limit key: %s", limitKey)
}
```

---

#### 2.4. Plan Checker Middleware

**File: `backend/middleware/plan_checker.go`** (NEW)

```go
package middleware

import (
    "net/http"
    "tukem-backend/db"
    "tukem-backend/utils"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

// RequireFeature checks if user has access to a feature
func RequireFeature(featureKey string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := c.Get("user").(*jwt.Token)
            claims := user.Claims.(*jwt.MapClaims)
            userID := (*claims)["user_id"].(string)
            
            hasAccess, err := utils.CheckFeatureAccess(db.DB, userID, featureKey)
            if err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check access"})
            }
            
            if !hasAccess {
                return c.JSON(http.StatusForbidden, map[string]interface{}{
                    "error": "Feature not available in your plan",
                    "feature": featureKey,
                    "code": "FEATURE_LIMITED",
                })
            }
            
            return next(c)
        }
    }
}

// CheckLimit checks if user has reached a limit before processing
func CheckLimit(limitKey string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := c.Get("user").(*jwt.Token)
            claims := user.Claims.(*jwt.MapClaims)
            userID := (*claims)["user_id"].(string)
            
            hasQuota, remaining, err := utils.CheckLimit(db.DB, userID, limitKey)
            if err != nil {
                return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check limit"})
            }
            
            if !hasQuota {
                return c.JSON(http.StatusForbidden, map[string]interface{}{
                    "error": "Limit reached",
                    "limit": limitKey,
                    "remaining": remaining,
                    "code": "LIMIT_REACHED",
                })
            }
            
            // Store remaining quota in context
            c.Set("remaining_quota", remaining)
            
            return next(c)
        }
    }
}
```

---

#### 2.5. Update Routes

**File: `backend/main.go`** (UPDATED)

```go
// Google OAuth routes
api.GET("/auth/google", handlers.GoogleLogin)
api.GET("/auth/google/callback", handlers.GoogleCallback(db.DB))

// Protected routes with plan checks
api.GET("/children/:id/export-pdf", handlers.ExportChildReport, 
    middleware.JWT([]byte(os.Getenv("JWT_SECRET"))),
    middleware.RequireFeature("pdf_export"))

api.POST("/children", handlers.CreateChild,
    middleware.JWT([]byte(os.Getenv("JWT_SECRET"))),
    middleware.CheckLimit("children"))

api.POST("/children/:id/measurements", handlers.CreateMeasurement,
    middleware.JWT([]byte(os.Getenv("JWT_SECRET"))),
    middleware.CheckLimit("measurements"))

// Plan management routes
api.GET("/plans", handlers.GetPlans)
api.GET("/user/plan", handlers.GetUserPlan, middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))
api.POST("/user/upgrade", handlers.UpgradePlan, middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))
```

---

### Phase 3: Frontend Implementation

#### 3.1. Update Auth Store

**File: `frontend/stores/auth.ts`** (UPDATED)

```typescript
export const useAuthStore = defineStore('auth', () => {
    // ... existing code ...
    
    const subscriptionPlan = ref<string | null>(null)
    const subscriptionStatus = ref<string | null>(null)
    const planFeatures = ref<Record<string, boolean>>({})
    
    function setUserPlan(plan: string, status: string, features: Record<string, boolean>) {
        subscriptionPlan.value = plan
        subscriptionStatus.value = status
        planFeatures.value = features
    }
    
    function hasFeature(featureKey: string): boolean {
        return planFeatures.value[featureKey] === true
    }
    
    // Google OAuth login
    async function loginWithGoogle() {
        window.location.href = `${apiBase}/api/auth/google`
    }
    
    // Handle OAuth callback
    function handleOAuthCallback(token: string) {
        setToken(token)
        // Fetch user data with plan info
        fetchUserProfile()
    }
    
    async function fetchUserProfile() {
        // Fetch user data including plan info
        // Update subscriptionPlan, subscriptionStatus, planFeatures
    }
    
    return {
        // ... existing returns ...
        subscriptionPlan,
        subscriptionStatus,
        planFeatures,
        hasFeature,
        loginWithGoogle,
        handleOAuthCallback,
        fetchUserProfile,
    }
})
```

---

#### 3.2. OAuth Callback Page

**File: `frontend/pages/auth/callback.vue`** (NEW)

```vue
<template>
  <div class="min-h-screen flex items-center justify-center">
    <div v-if="loading" class="text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
      <p class="mt-4 text-gray-600">Logging in...</p>
    </div>
    <div v-else-if="error" class="text-center text-red-600">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const error = ref(null)

onMounted(async () => {
  const token = route.query.token as string
  
  if (!token) {
    error.value = 'No token provided'
    loading.value = false
    return
  }
  
  try {
    authStore.handleOAuthCallback(token)
    await authStore.fetchUserProfile()
    router.push('/dashboard')
  } catch (e) {
    error.value = 'Failed to authenticate'
    loading.value = false
  }
})
</script>
```

---

#### 3.3. Plan Management Component

**File: `frontend/components/PlanUpgrade.vue`** (NEW)

```vue
<template>
  <div v-if="showUpgrade" class="bg-amber-50 border border-amber-200 rounded-lg p-4 mb-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="font-semibold text-amber-900">{{ message }}</h3>
        <p class="text-sm text-amber-700 mt-1">{{ description }}</p>
      </div>
      <NuxtLink to="/settings/billing" class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700">
        Upgrade Plan
      </NuxtLink>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  feature: {
    type: String,
    required: true
  },
  message: {
    type: String,
    default: 'Feature not available in free plan'
  },
  description: {
    type: String,
    default: 'Upgrade to paid plan to unlock this feature'
  }
})

const authStore = useAuthStore()
const showUpgrade = computed(() => {
  return !authStore.hasFeature(props.feature)
})
</script>
```

---

#### 3.4. Settings/Billing Page

**File: `frontend/pages/settings/billing.vue`** (NEW)

```vue
<template>
  <div class="p-6 max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold mb-6">Billing & Subscription</h1>
    
    <!-- Current Plan -->
    <div class="bg-white rounded-xl shadow-sm p-6 mb-6">
      <h2 class="text-xl font-bold mb-4">Current Plan</h2>
      <div class="flex items-center justify-between">
        <div>
          <p class="text-2xl font-bold">{{ currentPlan.name }}</p>
          <p class="text-gray-600">{{ currentPlan.description }}</p>
        </div>
        <span class="px-4 py-2 bg-green-100 text-green-700 rounded-lg">
          {{ subscriptionStatus }}
        </span>
      </div>
    </div>
    
    <!-- Available Plans -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="plan in plans" :key="plan.id" 
           :class="['rounded-xl p-6 border-2', plan.id === currentPlan.id ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200']">
        <h3 class="text-xl font-bold mb-2">{{ plan.name }}</h3>
        <p class="text-3xl font-bold mb-4">
          {{ formatPrice(plan.price_monthly) }}
          <span class="text-sm font-normal text-gray-600">/bulan</span>
        </p>
        <ul class="space-y-2 mb-6">
          <li v-for="feature in plan.features" :key="feature">
            ✓ {{ feature }}
          </li>
        </ul>
        <button v-if="plan.id !== currentPlan.id" 
                @click="upgradePlan(plan.id)"
                class="w-full px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700">
          {{ plan.id === 'free' ? 'Downgrade' : 'Upgrade' }}
        </button>
        <button v-else disabled
                class="w-full px-4 py-2 bg-gray-300 text-gray-600 rounded-lg cursor-not-allowed">
          Current Plan
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
// Implementation for plan management
</script>
```

---

### Phase 4: Feature Gating

#### 4.1. Update Components

**Contoh: PDF Export Button**

```vue
<template>
  <div>
    <PlanUpgrade v-if="!canExportPDF" 
                 feature="pdf_export"
                 message="PDF Export hanya tersedia di Paid Plan"
                 description="Upgrade ke Paid Plan untuk mengexport laporan PDF" />
    
    <button v-if="canExportPDF" @click="downloadPDF">
      Export PDF
    </button>
  </div>
</template>

<script setup>
const authStore = useAuthStore()
const canExportPDF = computed(() => authStore.hasFeature('pdf_export'))
</script>
```

---

## 📊 FEATURE MAPPING PER PLAN

| Feature | Free | Paid | Pro |
|---------|------|------|-----|
| **Children** | 1 | 3 | Unlimited |
| **Measurements** | 5/month | Unlimited | Unlimited |
| **Assessments** | 2/month | Unlimited | Unlimited |
| **Growth Charts** | ✅ | ✅ | ✅ |
| **Milestone Tracking** | ✅ | ✅ | ✅ |
| **PDF Export** | ❌ | ✅ | ✅ |
| **Denver II** | ❌ | ✅ | ✅ |
| **Immunization Schedule** | ❌ | ✅ | ✅ |
| **Recommendations** | ❌ | ✅ | ✅ |
| **Data Export (CSV)** | ❌ | ❌ | ✅ |
| **Priority Support** | ❌ | ❌ | ✅ |

---

## 🚀 IMPLEMENTATION TIMELINE

### Week 1: Database & Backend Foundation
- [ ] Day 1-2: Database migrations (009-012)
- [ ] Day 3-4: Update models dan utilities
- [ ] Day 5: Plan checker middleware

### Week 2: Google OAuth & Plan Management
- [ ] Day 1-2: Google OAuth implementation
- [ ] Day 3-4: Plan management API
- [ ] Day 5: Testing OAuth flow

### Week 3: Frontend Integration
- [ ] Day 1-2: Update auth store & OAuth callback
- [ ] Day 3-4: Billing/settings page
- [ ] Day 5: Feature gating di components

### Week 4: Testing & Polish
- [ ] Day 1-2: Integration testing
- [ ] Day 3-4: UI/UX polish
- [ ] Day 5: Documentation & deployment prep

---

## 🔒 SECURITY CONSIDERATIONS

1. **OAuth Security**
   - Validate state parameter
   - Verify token signature
   - Check token expiration

2. **Plan Validation**
   - Always validate on backend
   - Never trust frontend-only checks
   - Rate limiting for API calls

3. **Subscription Status**
   - Check expiration dates
   - Handle trial periods
   - Grace period for expired subscriptions

---

## 📝 ENVIRONMENT VARIABLES

```bash
# Google OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback

# Frontend URL for OAuth callback
FRONTEND_URL=http://localhost:3000

# Existing
JWT_SECRET=your-secret
DATABASE_URL=postgresql://...
```

---

## 🎯 NEXT STEPS

1. **Setup Google OAuth Credentials**
   - Create project di Google Cloud Console
   - Enable Google+ API
   - Create OAuth 2.0 credentials
   - Configure redirect URIs

2. **Database Migrations**
   - Apply migrations 009-012
   - Verify data integrity
   - Test plan seeding

3. **Backend Implementation**
   - Implement OAuth handler
   - Implement plan checker
   - Add middleware to routes

4. **Frontend Implementation**
   - Update auth store
   - Create OAuth callback page
   - Create billing/settings page
   - Add feature gating

---

---

## 🔧 KONFIGURASI DINAMIS PLANS & FEATURES

### ✅ Ya, Fitur Plans Bisa Dinamis!

Desain yang sudah dibuat **mendukung konfigurasi dinamis** tanpa perlu mengubah code. Berikut caranya:

#### 1. Features Configuration (JSONB)

Features disimpan di field `JSONB` di table `subscription_plans`, sehingga bisa diubah langsung di database:

```sql
-- Update features untuk plan tertentu
UPDATE subscription_plans 
SET features = jsonb_set(
    features, 
    '{pdf_export}', 
    'true'::jsonb
)
WHERE id = 'free';

-- Atau update multiple features sekaligus
UPDATE subscription_plans 
SET features = '{
    "pdf_export": true,
    "growth_charts": true,
    "milestone_tracking": true,
    "denver_ii": false,
    "immunization_schedule": false,
    "recommendations": false
}'::jsonb
WHERE id = 'free';
```

**Keuntungan:**
- ✅ Tidak perlu deploy code baru
- ✅ Bisa diubah kapan saja
- ✅ Support untuk feature flags baru
- ✅ Flexible untuk A/B testing

---

#### 2. Feature Limits Configuration

Limits disimpan di table `plan_feature_limits` yang bisa diubah dinamis:

```sql
-- Update limit untuk measurements di free plan
UPDATE plan_feature_limits 
SET limit_value = 10  -- Ubah dari 5 menjadi 10
WHERE plan_id = 'free' AND feature_key = 'measurements';

-- Tambah limit baru
INSERT INTO plan_feature_limits (plan_id, feature_key, limit_value, limit_type, description)
VALUES ('free', 'pdf_exports', 1, 'per_month', '1 PDF export per month');

-- Hapus limit (unlimited)
UPDATE plan_feature_limits 
SET limit_value = NULL 
WHERE plan_id = 'paid' AND feature_key = 'measurements';
```

**Keuntungan:**
- ✅ Bisa adjust limits tanpa code changes
- ✅ Support untuk promosi/limited time offers
- ✅ Easy to add new limit types

---

#### 3. Admin API untuk Manage Plans (RECOMMENDED)

Untuk kemudahan management, disarankan membuat **Admin API** untuk manage plans:

**File: `backend/handlers/admin_plans.go`** (NEW)

```go
package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "tukem-backend/db"
    "tukem-backend/models"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/jmoiron/sqlx"
    "github.com/labstack/echo/v4"
)

// GetPlans - Get all subscription plans (public endpoint)
func GetPlans(c echo.Context) error {
    var plans []models.SubscriptionPlan
    err := db.DB.Select(&plans, `
        SELECT * FROM subscription_plans 
        WHERE is_active = true 
        ORDER BY price_monthly ASC
    `)
    
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch plans"})
    }
    
    return c.JSON(http.StatusOK, plans)
}

// AdminGetPlans - Get all plans with details (admin only)
func AdminGetPlans(c echo.Context) error {
    // Check if user is admin (add role check)
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*jwt.MapClaims)
    role := (*claims)["role"].(string)
    
    if role != "admin" {
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
    }
    
    var plans []models.SubscriptionPlan
    err := db.DB.Select(&plans, `SELECT * FROM subscription_plans ORDER BY price_monthly ASC`)
    
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch plans"})
    }
    
    return c.JSON(http.StatusOK, plans)
}

// AdminUpdatePlanFeatures - Update plan features dynamically
func AdminUpdatePlanFeatures(c echo.Context) error {
    // Check admin role
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*jwt.MapClaims)
    role := (*claims)["role"].(string)
    
    if role != "admin" {
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
    }
    
    planID := c.Param("plan_id")
    var req struct {
        Features map[string]interface{} `json:"features"`
    }
    
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
    }
    
    // Convert to JSONB
    featuresJSON, _ := json.Marshal(req.Features)
    
    _, err := db.DB.Exec(`
        UPDATE subscription_plans 
        SET features = $1::jsonb, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
    `, string(featuresJSON), planID)
    
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update plan"})
    }
    
    return c.JSON(http.StatusOK, map[string]string{"message": "Plan updated successfully"})
}

// AdminUpdatePlanLimit - Update plan limit dynamically
func AdminUpdatePlanLimit(c echo.Context) error {
    // Check admin role
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*jwt.MapClaims)
    role := (*claims)["role"].(string)
    
    if role != "admin" {
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
    }
    
    planID := c.Param("plan_id")
    featureKey := c.Param("feature_key")
    
    var req struct {
        LimitValue *int   `json:"limit_value"` // NULL = unlimited
        LimitType  string `json:"limit_type"`
        Description string `json:"description"`
    }
    
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
    }
    
    // Upsert limit
    _, err := db.DB.Exec(`
        INSERT INTO plan_feature_limits (plan_id, feature_key, limit_value, limit_type, description)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (plan_id, feature_key) 
        DO UPDATE SET 
            limit_value = EXCLUDED.limit_value,
            limit_type = EXCLUDED.limit_type,
            description = EXCLUDED.description,
            updated_at = CURRENT_TIMESTAMP
    `, planID, featureKey, req.LimitValue, req.LimitType, req.Description)
    
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update limit"})
    }
    
    return c.JSON(http.StatusOK, map[string]string{"message": "Limit updated successfully"})
}

// AdminCreatePlan - Create new plan dynamically
func AdminCreatePlan(c echo.Context) error {
    // Check admin role
    user := c.Get("user").(*jwt.Token)
    claims := user.Claims.(*jwt.MapClaims)
    role := (*claims)["role"].(string)
    
    if role != "admin" {
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
    }
    
    var plan models.SubscriptionPlan
    if err := c.Bind(&plan); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
    }
    
    featuresJSON, _ := json.Marshal(plan.Features)
    
    _, err := db.DB.Exec(`
        INSERT INTO subscription_plans (
            id, name, description, price_monthly, price_yearly, currency,
            trial_days, max_children, max_measurements_per_month, 
            max_assessments_per_month, features, is_active
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb, $12)
    `, plan.ID, plan.Name, plan.Description, plan.PriceMonthly, plan.PriceYearly,
        plan.Currency, plan.TrialDays, plan.MaxChildren, plan.MaxMeasurementsPerMonth,
        plan.MaxAssessmentsPerMonth, string(featuresJSON), plan.IsActive)
    
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create plan"})
    }
    
    return c.JSON(http.StatusCreated, plan)
}
```

**Admin Routes:**
```go
// Admin routes (protected with admin role check)
admin := api.Group("/admin")
admin.Use(customMiddleware.JWTMiddleware())
admin.Use(customMiddleware.RequireAdmin()) // New middleware

admin.GET("/plans", handlers.AdminGetPlans)
admin.PUT("/plans/:plan_id/features", handlers.AdminUpdatePlanFeatures)
admin.PUT("/plans/:plan_id/limits/:feature_key", handlers.AdminUpdatePlanLimit)
admin.POST("/plans", handlers.AdminCreatePlan)
```

---

#### 4. Admin Middleware

**File: `backend/middleware/admin.go`** (NEW)

```go
package middleware

import (
    "net/http"
    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
)

// RequireAdmin checks if user has admin role
func RequireAdmin() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            user := c.Get("user").(*jwt.Token)
            claims := user.Claims.(*jwt.MapClaims)
            role := (*claims)["role"].(string)
            
            if role != "admin" {
                return c.JSON(http.StatusForbidden, map[string]string{
                    "error": "Admin access required",
                })
            }
            
            return next(c)
        }
    }
}
```

---

#### 5. Frontend Admin Panel (Optional)

Untuk kemudahan management, bisa dibuat admin panel di frontend:

**File: `frontend/pages/admin/plans.vue`** (NEW)

```vue
<template>
  <div class="p-6 max-w-6xl mx-auto">
    <h1 class="text-3xl font-bold mb-6">Manage Subscription Plans</h1>
    
    <!-- Plans List -->
    <div class="space-y-4 mb-8">
      <div v-for="plan in plans" :key="plan.id" class="bg-white rounded-xl shadow-sm p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-2xl font-bold">{{ plan.name }}</h2>
          <button @click="editPlan(plan)" class="px-4 py-2 bg-indigo-600 text-white rounded-lg">
            Edit
          </button>
        </div>
        
        <!-- Features Editor -->
        <div class="grid grid-cols-2 gap-4">
          <div v-for="(value, key) in plan.features" :key="key" class="flex items-center gap-2">
            <input 
              type="checkbox" 
              :checked="value"
              @change="updateFeature(plan.id, key, $event.target.checked)"
              class="w-4 h-4"
            />
            <label class="text-sm">{{ key }}</label>
          </div>
        </div>
        
        <!-- Limits Editor -->
        <div class="mt-4">
          <h3 class="font-semibold mb-2">Limits</h3>
          <div class="space-y-2">
            <div v-for="limit in plan.limits" :key="limit.feature_key" class="flex items-center gap-2">
              <input 
                type="number" 
                :value="limit.limit_value"
                @change="updateLimit(plan.id, limit.feature_key, $event.target.value)"
                class="w-24 px-2 py-1 border rounded"
                placeholder="Unlimited"
              />
              <span class="text-sm">{{ limit.feature_key }} ({{ limit.limit_type }})</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// Implementation untuk manage plans
</script>
```

---

## 📝 CONTOH PENGGUNAAN

### Scenario 1: Menambah Fitur ke Free Plan

```sql
-- Tambahkan PDF Export ke Free Plan (1x per bulan)
UPDATE subscription_plans 
SET features = jsonb_set(features, '{pdf_export}', 'true'::jsonb)
WHERE id = 'free';

-- Tambahkan limit 1 PDF export per bulan
INSERT INTO plan_feature_limits (plan_id, feature_key, limit_value, limit_type)
VALUES ('free', 'pdf_exports', 1, 'per_month')
ON CONFLICT (plan_id, feature_key) DO UPDATE SET limit_value = 1;
```

### Scenario 2: Promosi - Naikkan Limit Free Plan

```sql
-- Promosi: Free plan bisa 10 measurements per bulan (biasanya 5)
UPDATE plan_feature_limits 
SET limit_value = 10
WHERE plan_id = 'free' AND feature_key = 'measurements';
```

### Scenario 3: Tambah Plan Baru

```sql
-- Tambah "Starter Plan" dengan fitur di antara Free dan Paid
INSERT INTO subscription_plans (id, name, description, price_monthly, features)
VALUES (
    'starter',
    'Starter Plan',
    'Plan untuk keluarga kecil',
    49000,
    '{
        "pdf_export": true,
        "growth_charts": true,
        "milestone_tracking": true,
        "denver_ii": false,
        "immunization_schedule": true,
        "recommendations": false
    }'::jsonb
);
```

---

## ✅ KESIMPULAN

**Ya, semua fitur plans bisa dikonfigurasi secara dinamis!**

**Cara konfigurasi:**
1. ✅ **Langsung di Database** - Update SQL langsung
2. ✅ **Via Admin API** - API endpoints untuk manage plans (recommended)
3. ✅ **Via Admin Panel** - UI untuk manage plans (optional, untuk kemudahan)

**Keuntungan:**
- ✅ Tidak perlu deploy code baru untuk ubah plans
- ✅ Bisa adjust pricing, features, limits kapan saja
- ✅ Support untuk promosi dan limited offers
- ✅ Easy to add new plans atau features
- ✅ Flexible untuk A/B testing

**Rekomendasi:**
- Implement Admin API untuk manage plans
- Buat admin panel untuk kemudahan management
- Document semua feature keys yang tersedia
- Buat audit log untuk track perubahan plans

---

**Updated:** 21 November 2025

