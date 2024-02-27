/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: ["./views/**/*.html"],
  theme: {
    extend: {
      colors: {
        'light': '#eeeeee',
        'dark': '#333333',
        'gray-light': '#5f5f5f',
        'gray-dark': '#aaaaaa',
      },
      fontFamily: {
        body: ['scope-one', 'sans-serif'],
        title: ['post-no-bills-colombo', 'sans-serif'],
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}

