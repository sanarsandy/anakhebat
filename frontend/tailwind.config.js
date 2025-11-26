/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./components/**/*.{js,vue,ts}",
        "./layouts/**/*.vue",
        "./pages/**/*.vue",
        "./plugins/**/*.{js,ts}",
        "./app.vue",
        "./error.vue",
    ],
    theme: {
        extend: {
            colors: {
                // Jurnal Si Kecil Brand Colors - Minimalis
                'jurnal': {
                    // Jurnal Teal (Tosca) - Primary Brand Color
                    'teal': {
                        50: '#F0FDFA',   // Very light background
                        100: '#E6FCF7',  // Light background
                        200: '#B2F5EA',  // Subtle highlight
                        300: '#7DD3C0',  // Light accent
                        400: '#4FD1B5',  // Medium
                        500: '#2DD4BF',  // Primary Teal (from logo)
                        600: '#14B8A6',  // Primary Dark (main brand)
                        700: '#0D9488',  // Darker
                        800: '#0F766E',  // Very dark
                        900: '#134E4A',  // Darkest
                    },
                    // Warm Coral (Koral) - Secondary/Accent
                    'coral': {
                        50: '#FFF5F5',
                        100: '#FFE5E5',
                        200: '#FFCCCC',
                        300: '#FF9999',
                        400: '#FF6B6B',  // Warm Coral (from logo)
                        500: '#FF5252',  // Coral Medium
                        600: '#E63946',  // Coral Dark
                        700: '#C92A2A',
                        800: '#A61E1E',
                        900: '#7F1D1D',
                    },
                    // Sunny Gold (Kuning Emas) - Accent
                    'gold': {
                        50: '#FFFBEB',
                        100: '#FEF3C7',
                        200: '#FDE68A',
                        300: '#FCD34D',
                        400: '#FBBF24',  // Sunny Gold (from logo)
                        500: '#F59E0B',  // Gold Medium
                        600: '#D97706',
                        700: '#B45309',
                        800: '#92400E',
                        900: '#78350F',
                    },
                    // Soft Charcoal - Text Color (minimalis)
                    'charcoal': {
                        50: '#F9FAFB',   // Off-White background
                        100: '#F3F4F6',  // Light gray background
                        200: '#E5E7EB',  // Border
                        300: '#D1D5DB',  // Light border
                        400: '#9CA3AF',  // Muted text
                        500: '#6B7280',  // Medium text
                        600: '#4B5563',  // Body text
                        700: '#374151',  // Soft Charcoal (main text)
                        800: '#1F2937',  // Heading
                        900: '#111827',  // Dark heading
                    },
                },
            },
            fontFamily: {
                'sans': ['Poppins', 'system-ui', 'sans-serif'],
            },
            borderRadius: {
                'soft': '12px',
                'soft-lg': '16px',
                'pill': '9999px',
            },
        },
    },
    plugins: [],
}
