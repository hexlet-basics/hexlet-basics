import "@mantine/core";

declare module "@mantine/core" {
  export interface TitleProps {
    /** Делает размер заголовка адаптивным (clamp на основе size/order) */
    responsive?: boolean;
  }

  export interface TextProps {
    responsive?: string;
  }
}
