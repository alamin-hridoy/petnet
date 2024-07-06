module.exports = {
  theme: {
    // PROXTERA colors
    //
    //
    colors: {
      transparent: "transparent",
      white: "#FFF",
      black: "#000",

      grey: {
        "dark-5": "#5E647D",
        "dark-10": "#959BB2",
        "dark-15": "#495272",
        default: "#D5DAED",
      },
      blue: {
        "dark-5": "#2c5cd6",
        "dark-10": "#271f4e",
        "dark-15": "#15226d",
        "dark-20": "#112661",
        "dark-25": "#15226D",
        "dark-30": "#101b5d",
        default: "#517ef0",
        "light-5": "#4e79e9",
        "light-10": "#bae8fb",
      },
      red: {
        default: "#C53030",
        "light-5": "#fff6f6",
        "light-10": "#ffebeb",
      },
      green: {
        default: "#86F4BE",
      },
      purple: {
        default: "#40337d",
      }
    },

    //
    fontFamily: {
      sans: [
        '"Rubik"',

        // TODO: What are the fallback fonts to use?
        "system-ui",
        "-apple-system",
        "BlinkMacSystemFont",
        '"Segoe UI"',
        "Roboto",
        '"Helvetica Neue"',
        "Arial",
      ],
    },
    extend: {
      screens: {
        'xs': {max: '768px'},
        'max-lg': {max: '1024px'},
      },
      maxWidth: {
        '420': '26.25rem', // 420px
      },
      height: {
        "14": "3.438rem", // 55px
      },
      lineHeight: {
        "14": "3.438rem", // 55px
      },
      minHeight: {
        "800": "50rem", // 800px
      },
      fontSize: {
        "60": "3.75rem", // 60px
        "80": "5rem", // 60px
        "96": "6rem", // 96px
      },
      zIndex: {
        '-1': '-1',
      }
    }
  },
  variants: {
    display: ['responsive', 'hover', 'focus','group-hover'],
  },
  plugins: [],
  future: {
    removeDeprecatedGapUtilities: true,
  },
};
