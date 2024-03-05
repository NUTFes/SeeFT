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
      'black-100': '#000',
      'white-100': '#fff',
      'accent-primary': '#E3E2FF',
      'accent-secondary': '#77B6FF',
      'accent-tertiary':  '#D7E9FF',
      'surfase-primary': '#E3E2FF',
      'surfase-sesondary': '#EDF5FF',
      'emphasis': '#132941',
      'caution': '#FF4D4D',
    },
    extend: {
      fontFamily: {
        'sans': ['"Proxima Nova"', ...defaultTheme.fontFamily.sans],
      },
      width: {
        '1/2': '50%',
        '3/4': '75%',
        '9/10': '90%',
      },
    },
  },
  plugins: [
    function ({ addUtilities }) {
      const newUtilities = {
        '.text-shadow': {
          textShadow: '0px 2px 3px darkgrey',
        },
        '.text-shadow-md': {
          textShadow: '0px 3px 3px darkgrey',
        },
        '.text-shadow-logo': {
          textShadow: '4px 2px 0px #333',
        },
        '.text-shadow-none': {
          textShadow: 'none',
        },
      };

      addUtilities(newUtilities);
    },
  ],
};