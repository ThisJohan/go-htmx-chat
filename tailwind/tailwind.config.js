/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../views/**/*.{templ,html}"],
  theme: {
    extend: {
      colors:{
        'primary': '#222831',
        'primary-foreground': '#fff'
      }
    },
  },
  plugins: [],
  darkMode: 'selector',
}