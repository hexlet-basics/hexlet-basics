import { Box } from "@mantine/core";
import type { BoxProps } from "@mantine/core";
import { motion } from "motion/react";
import type React from "react";
import type { ReactElement } from "react";

type HoverLiftProps = BoxProps &
  React.ComponentProps<typeof motion.div> & {
    lift?: number;
    hoverScale?: number;
  };

export function HoverLift({
  lift = 6,
  hoverScale = 1.02,
  whileHover,
  transition,
  ...props
}: HoverLiftProps): ReactElement {
  return (
    <Box
      component={motion.div}
      {...props}
      bdrs={16}
      whileHover={{
        boxShadow: "0 18px 45px rgba(0,0,0,0.14)",
        scale: hoverScale,
        y: -lift,
        ...(typeof whileHover === "object" && !Array.isArray(whileHover) ? whileHover : {}),
      }}
      transition={{
        damping: 35,
        stiffness: 500,
        type: "spring",
        ...transition,
      }}
    />
  );
}
