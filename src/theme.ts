import { Anchor, createTheme } from "@mantine/core";

// Site-wide Mantine theme. `Anchor` is blue by default, which suits links
// inside content but not the header and footer navigation: those take
// `variant="nav"` and keep the surrounding text colour (underline on hover
// stays Anchor's default).
//
// Headings shrink with the viewport as Bootstrap 5's RFS does: below 1200px
// (`xl`) the size is `rem + vw`, from there up it is Mantine's own size, so the
// desktop looks as before. The coefficients are Bootstrap's own for h1–h4
// scaled to Mantine's sizes (34, 26, 22, 18px); h5 and h6 are too small to shrink.
export const theme = createTheme({
  headings: {
    sizes: {
      h1: { fontSize: "min(calc(1.169rem + 1.275vw), 2.125rem)" },
      h2: { fontSize: "min(calc(1.077rem + 0.731vw), 1.625rem)" },
      h3: { fontSize: "min(calc(1.021rem + 0.471vw), 1.375rem)" },
      h4: { fontSize: "min(calc(0.956rem + 0.225vw), 1.125rem)" },
    },
  },
  components: {
    Anchor: Anchor.extend({
      styles: (_theme, props) => (props.variant === "nav" ? { root: { color: "inherit" } } : {}),
    }),
  },
});
