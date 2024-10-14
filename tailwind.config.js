/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/views/*.templ", "./web/views/user.templ", "./web/**/*.templ", "./web/views/*.templ" ],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["light", "dark", "black", "cupcake"],
  },
  plugins: [
    require('daisyui'),
  ],
}

