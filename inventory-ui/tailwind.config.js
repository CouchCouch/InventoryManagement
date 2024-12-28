/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,jsx}"],
  theme: {
    extend: {
      'colors': 
      { 'jasper': { DEFAULT: '#bf4e30', 100: '#26100a', 200: '#4d2013', 300: '#732f1d', 400: '#993f26', 500: '#bf4e30', 600: '#d46e52', 700: '#df927d', 800: '#e9b7a9', 900: '#f4dbd4' }, 'ash_gray': { DEFAULT: '#c6ccb2', 100: '#2a2e1f', 200: '#555c3d', 300: '#7f8a5c', 400: '#a4ad85', 500: '#c6ccb2', 600: '#d1d6c2', 700: '#dde0d1', 800: '#e8ebe0', 900: '#f4f5f0' }, 'dark_green': { DEFAULT: '#093824', 100: '#020b07', 200: '#04170f', 300: '#062216', 400: '#072e1d', 500: '#093824', 600: '#168555', 700: '#22d286', 800: '#66e6af', 900: '#b3f3d7' }, 'lavender': { DEFAULT: '#e5eafa', 100: '#0f1f51', 200: '#1f3da1', 300: '#4467db', 400: '#95a9eb', 500: '#e5eafa', 600: '#eaeefb', 700: '#f0f2fc', 800: '#f5f7fd', 900: '#fafbfe' }, 'aquamarine': { DEFAULT: '#78fecf', 100: '#004a30', 200: '#019461', 300: '#01de91', 400: '#2cfeb4', 500: '#78fecf', 600: '#91fed8', 700: '#adffe2', 800: '#c8ffec', 900: '#e4fff5' } }
    },
  },
  plugins: [],
}

