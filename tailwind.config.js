/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './**/*.{templ,css,js}',
  ],
  theme: {
    container: {
      center: true,
      padding: '0.5rem',
    },
    screens: {
      "sm": "640px",
      "md": "768px",
    },
    extend: {
      colors: {
        warm: {
          100: "oklch(98.23% 0.0152 67.74)",
          200: "oklch(96% 0.0152 67.74)",
          300: "oklch(92% 0.0152 67.74)",
          400: "oklch(80% 0.0152 67.74)",
          500: "oklch(70% 0.0212 67.74)",
          600: "oklch(60% 0.0212 67.74)",
          700: "oklch(50% 0.0234 67.74)",
          800: "oklch(30% 0.0234 67.74)",
          900: "oklch(20% 0.0234 67.74)",
        },
      },
      backgroundImage: ({ theme }) => ({
        "noise": "url(\"/assets/images/noise.png\")",
        "vignette": "radial-gradient(transparent, rgba(0,0,0,0.3))"
      }),
    },
    fontFamily: {
      sans: ['Ysabeau', 'Helvetica Neue', 'Arial', 'sans-serif'],
      serif: ['Cormorant', 'Palatino', "'Times New Roman'", 'serif'],
    },
  },
  plugins: [

  ],
  corePlugins: {
    preflight: true
  }
}

