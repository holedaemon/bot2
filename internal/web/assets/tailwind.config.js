/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../templates/*.qtpl",
    "./node_modules/flowbite/**/*.js"
  ],
  darkMode: "class",
  theme: {
    fontFamily: {
      sans: ["Comic Sans MS", "sans-serif"]
    },
    extend: {},
  },
  plugins: [
    require("flowbite/plugin")
  ],
}

