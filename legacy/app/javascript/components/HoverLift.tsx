import { Box, type BoxProps } from "@mantine/core";
import { motion } from "motion/react";
import type React from "react";

type HoverLiftProps = BoxProps &
  React.ComponentProps<typeof motion.div> & {
    lift?: number;
    hoverScale?: number;
  };

export function HoverLift({
  lift = 6,
  hoverScale = 1.02,
  style,
  whileHover,
  transition,
  ...props
}: HoverLiftProps) {
  return (
    <Box
      component={motion.div}
      {...props}
      style={
        {
          borderRadius: 16,
          willChange: "transform",
          ...(style as React.CSSProperties),
        } as React.CSSProperties
      }
      whileHover={{
        y: -lift,
        scale: hoverScale,
        boxShadow: "0 18px 45px rgba(0,0,0,0.14)",
        ...(whileHover as object),
      }}
      transition={{
        type: "spring",
        stiffness: 500,
        damping: 35,
        ...transition,
      }}
    />
  );
}
