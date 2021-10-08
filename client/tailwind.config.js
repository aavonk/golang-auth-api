module.exports = {
  purge: ['./src/**/*.{js,jsx,ts,tsx}', './public/index.html'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      // textColor: theme => theme('colors')
      textColor: {
        'primary-dark': '#1a1f36',
        'primary-reg': '#3c4257',
        error: '#cd3d64',
      },
      minWidth: {
        1085: '1085px',
        1080: '1080px',
      },
      width: {
        1085: '1085px',
        1080: '1080px',
      },
    },
  },

  variants: {
    extend: {},
  },
  plugins: [],
};
