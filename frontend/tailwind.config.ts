import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./src/**/*.{js,ts,jsx,tsx,mdx}'],
  theme: {
    extend: {
      colors: {
        primary: '#ff0036',
        'primary-pressed': '#e60030',
        ink: '#000000',
        body: '#33332e',
        mute: '#62625b',
        ash: '#91918c',
        stone: '#c8c8c1',
        hairline: '#dadad3',
        canvas: '#ffffff',
        surface: '#f6f6f3',
        'surface-soft': '#fbfbf9',
        'secondary-bg': '#e5e5e0',
      },
      borderRadius: {
        card: '16px',
        lg: '32px',
        full: '9999px',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      maxWidth: {
        container: '1280px',
      },
    },
  },
  plugins: [],
};

export default config;
