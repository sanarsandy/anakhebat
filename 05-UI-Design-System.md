# UI Design System & Style Guide
**Project Name:** Tukem (Tumbuh Kembang)
**Theme:** Modern, Trustworthy, Playful yet Professional.

---

## 1. Color Palette
We use a **Soft Pastel & Vibrant** mix to convey "Child-Friendly" but "Medical Grade" trust.

### Primary Colors (Brand)
- **Primary Blue:** `#4F46E5` (Indigo 600) - *Trust, Technology, Professionalism*
- **Primary Soft:** `#E0E7FF` (Indigo 100) - *Backgrounds, Highlights*

### Secondary Colors (Status)
- **Success (Normal/Safe):** `#10B981` (Emerald 500) - *Good Growth, Milestones Met*
- **Warning (Risk):** `#F59E0B` (Amber 500) - *Borderline, Needs Attention*
- **Danger (Stunting/Delay):** `#EF4444` (Red 500) - *Critical, Red Flags*

### Neutral Colors
- **Text Dark:** `#1F2937` (Gray 800) - *Readability*
- **Text Muted:** `#6B7280` (Gray 500) - *Subtitles*
- **Background:** `#F9FAFB` (Gray 50) - *App Background*
- **Surface:** `#FFFFFF` (White) - *Cards, Modals*

---

## 2. Typography
**Font Family:** `Inter` or `Outfit` (Google Fonts).
- **Headings:** `Outfit`, Bold (700).
- **Body:** `Inter`, Regular (400) & Medium (500).

### Scale
- **H1 (Page Title):** 24px / 32px (Mobile), 30px / 36px (Desktop)
- **H2 (Section Title):** 20px / 28px
- **H3 (Card Title):** 18px / 24px
- **Body:** 16px / 24px
- **Small:** 14px / 20px

---

## 3. Components

### 3.1 Cards (Container)
- **Style:** White background, Rounded-xl (12px), Soft Shadow (`shadow-sm` or `shadow-md`).
- **Usage:** Enclose every major section (Growth Chart, Milestone List).

### 3.2 Buttons
- **Primary:** `bg-indigo-600 text-white hover:bg-indigo-700 rounded-lg px-4 py-2 font-medium transition-all`.
- **Secondary:** `bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 rounded-lg px-4 py-2`.
- **Danger:** `bg-red-50 text-red-600 hover:bg-red-100 rounded-lg px-4 py-2`.

### 3.3 Charts (Growth Curve)
- **Library:** Chart.js
- **Background:** White
- **Grid Lines:** Dashed, Gray 200
- **Datasets:**
    - *Child's Data:* Solid Blue Line, Circle Points.
    - *WHO SD 0 (Median):* Green Dotted Line.
    - *WHO SD -2/+2:* Orange Dotted Line.
    - *WHO SD -3/+3:* Red Dotted Line.

### 3.4 Milestone Checklist
- **Item:** Flex row, Text on left, Toggle/Radio on right.
- **Checked (Yes):** Green border, Green background tint.
- **Unchecked (No):** Gray border.

---

## 4. Layout Structure (Dashboard)
- **Sidebar (Desktop) / Bottom Nav (Mobile):** Navigation (Home, Growth, Development, Settings).
- **Header:** User Profile, Child Selector (Dropdown).
- **Main Content:** Padding 4 (16px) Mobile, Padding 8 (32px) Desktop. Max-width `7xl`.
