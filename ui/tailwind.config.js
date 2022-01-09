module.exports = {
  content: ["./src/**/*.{html,js}"],
  purge: [ "./src/**/*.{js,jsx,ts,tsx}", "./public/index.html"],
  darkMode: false,
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    styled: true,
    themes: [{
      'mastro' : {
        'primary': '#3ab2d9',
        'primary-focus': '#2288aa',
        'primary-content': '#ffffff',
        'secondary': '#ecf000',
        'secondary-focus': '#bdb600',
        'secondary-content': '#ffffff',
        'accent': '#37cdbe',
        'accent-focus': '#2aa79b',
        'accent-content': '#ffffff',
        'neutral': '#3d4451',
        'neutral-focus': '#2a2e37',
        'neutral-content': '#ffffff',
        'base-100': '#ffffff',
        'base-200': '#f9fafb',
        'base-300': '#d1d5db',
        'base-content': '#1f2937',
        'info': '#2094f3',
        'success': '#009485',
        'warning': '#ff9900',
        'error': '#ff5724',
      }
    }],
    base: true,
    utils: true,
    logs: true,
    rtl: false,
  },
}
