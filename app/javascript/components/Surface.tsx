import type {
  CardProps,
  ElementProps,
  MantineColor,
  MantineGradient,
} from "@mantine/core";
import { Card, useMantineTheme } from "@mantine/core";
import { useMemo } from "react";

export type SurfaceVariant =
  | "filled"
  | "light"
  | "outline"
  | "transparent"
  | "white"
  | "default"
  | "gradient"
  | "muted";

export interface SurfaceProps
  extends Omit<CardProps, "variant">,
    ElementProps<"div"> {
  variant?: SurfaceVariant;
  color?: MantineColor;
  gradient?: MantineGradient;
  autoContrast?: boolean;
}

export const Surface = ({
  variant = "light",
  color,
  gradient,
  autoContrast,
  style,
  withBorder = false,
  ...others
}: SurfaceProps) => {
  const theme = useMantineTheme();

  const combinedStyle = useMemo(() => {
    const isMuted = variant === "muted";

    const colors = isMuted
      ? null
      : theme.variantColorResolver({
          color: color || theme.primaryColor,
          theme,
          gradient,
          variant: variant || "light",
          autoContrast,
        });

    const variantStyles = isMuted
      ? {
          backgroundColor:
            "light-dark(var(--mantine-color-gray-1), var(--mantine-color-dark-6))",
          color: undefined,
        }
      : {
          backgroundColor: color || variant ? colors?.background : undefined,
          color: autoContrast ? colors?.color : undefined,
          border: withBorder ? undefined : colors?.border,
        };

    return {
      ...variantStyles,
      ...style,
    };
  }, [variant, color, gradient, autoContrast, style, theme, withBorder]);

  return <Card {...others} withBorder={withBorder} style={combinedStyle} />;
};

Surface.displayName = "@hexlet/Surface";
