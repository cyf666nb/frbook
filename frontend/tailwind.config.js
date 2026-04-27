/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'primary': '#1A1A1A',
        'secondary': '#666666',
        'accent-orange': '#E85D2B',
        'accent-cream': '#F5F0E8',
        'accent-blue': '#2B9CD8',
        'accent-green': '#4ADE80',
      },
      fontFamily: {
        'serif': ['Playfair Display', 'Times New Roman', 'Georgia', 'serif'],
        'sans': ['Inter', 'Helvetica Neue', 'Arial', 'sans-serif'],
      },
      boxShadow: {
        'card': '0 20px 40px rgba(0, 0, 0, 0.1)',
        'card-hover': '0 30px 60px rgba(0, 0, 0, 0.15)',
      },
    },
  },
  plugins: [],
}
