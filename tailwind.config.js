/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: [
    './routes/**/*.{templ,css,js}',
    './lib/**/*.{templ,css,js}',
  ],
  theme: {
    container: {
      center: true,
      padding: '0.5rem',
    },
    screens: {
      "hidpi": { "raw": "(-webkit-min-device-pixel-ratio: 1.5), (min-resolution: 144dpi)" },
      "sm": "450px",
      "md": "768px",
    },
    extend: {
      typography: ({ theme }) => ({
        warm: {
          css: {
            '--tw-prose-body': theme('colors.warm[800]'),
            '--tw-prose-headings': theme('colors.warm[900]'),
            '--tw-prose-lead': theme('colors.warm[700]'),
            '--tw-prose-links': theme('colors.warm[900]'),
            '--tw-prose-bold': theme('colors.warm[900]'),
            '--tw-prose-counters': theme('colors.warm[600]'),
            '--tw-prose-bullets': theme('colors.warm[600]'),
            '--tw-prose-hr': theme('colors.warm[300]'),
            '--tw-prose-quotes': theme('colors.warm[900]'),
            '--tw-prose-quote-borders': theme('colors.warm[300]'),
            '--tw-prose-captions': theme('colors.warm[800]'),
            '--tw-prose-code': theme('colors.warm[900]'),
            '--tw-prose-pre-code': theme('colors.warm[900]'),
            '--tw-prose-pre-bg': theme('colors.warm[200]'),
            '--tw-prose-th-borders': theme('colors.warm[400]'),
            '--tw-prose-td-borders': theme('colors.warm[300]'),
          }
        }
      }),
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
        "vignette": "radial-gradient(transparent, rgba(0,0,0,0.5))"
      }),
    },
    fontFamily: {
      sans: ['Ysabeau', 'Helvetica Neue', 'Arial', 'sans-serif'],
      serif: ['Cormorant', 'Palatino', "'Times New Roman'", 'serif'],
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
  corePlugins: {
    preflight: true
  }
}

