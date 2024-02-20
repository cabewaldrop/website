module.exports = {
  content: ["./views/*.html"],
  theme: {
    extend: {
      typography: {
        DEFAULT: {
          css: {
            img: {
              margin: "auto",
            },
          },
        },
      },
    },
  },
  plugins: [require("@tailwindcss/typography")],
};
