/** @type {import('tailwindcss').Config} */
module.exports = {
  tailwindConfig: './styles/tailwind.config.js',
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    fontSize: {
      'sm': '0.8rem',
      'base': '1rem',
      'xl': '1.5rem',
      '2xl': `2rem`,
      '3xl':  '3rem',
    },
    colors: {
      'black-0': '#000',
      'white-0': '#fff',
      'accent-1': '#E3E2FF',
      'accent-2': '#77B6FF',
      'accent-3':  '#D7E9FF',
      'surface-1': '#E3E2FF',
      'surface-2': '#EDF5FF',
      'emphasis': '#132941',
      'caution': '#FF4D4D',
    },
    extend: {
      width: {
        '1/2': '50%',
        '3/4': '75%',
        '9/10': '90%',
      },
    },
  },
  plugins: [
  ],
};