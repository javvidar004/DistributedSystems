import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        // Custom color palette
        primary: '#3B82F6', // Blue
        secondary: '#A16207', // Brown/Gold
        background: '#F9FAF6', // Off-white for backgrounds
        'dark-background': '#1F2937', // Dark gray for nav/sidebars
      },
    },
  },
  plugins: [],
}
export default config