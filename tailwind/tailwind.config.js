/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "../templates/**/*.{gohtml,html}" // For local computer
    // "/templates/**/*.{gohtml,html}" // For mount it in Docker volume
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

