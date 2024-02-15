/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: ["./views/**/*.html"],
  theme: {
    extend: {
      colors: {
        'light': '#eeeeee',
        'dark': '#333333',
        'text-light': '#ffffff',
        'text-dark': '#000000',
        'gray-light': '#5f5f5f',
        'gray-dark': '#aaaaaa',
      },
      fontFamily: {
        'body': ['Scope One', 'sans-serif'],
        'title': ['Post No Bills Colombo', 'sans-serif'],
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

