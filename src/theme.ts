import { Anchor, createTheme } from "@mantine/core";

// Site-wide Mantine theme. `Anchor` is blue by default, which suits links
// inside content but not the header and footer navigation: those take
// `variant="nav"` and keep the surrounding text colour (underline on hover
// stays Anchor's default).
export const theme = createTheme({
  components: {
    Anchor: Anchor.extend({
      styles: (_theme, props) => (props.variant === "nav" ? { root: { color: "inherit" } } : {}),
    }),
  },
});
