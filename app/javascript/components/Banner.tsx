import { usePage } from "@inertiajs/react";
import { Alert, Anchor, Box, Card } from "@mantine/core";
import { propsForExternalLink } from "@/lib/utils";

export function Banner() {
  const { locale } = usePage().props;
  // const { t } = useTranslation();

  return (
    <>
      {locale === "ru" && (
        <Box
          variant="gradient"
          bg="var(--app-cta-gradient)"
          p="xs"
          my="xl"
          ta="center"
        >
          <Anchor
            href="https://ru.hexlet.io/courses_for_beginners?utm_source=code-basics&utm_medium=referral&utm_campaign=newyear2025"
            {...propsForExternalLink()}
          >
            Курсы с наставником и опытом работы во время учебы со{" "}
            <b>скидкой -30%</b> → Подробнее
          </Anchor>
        </Box>
      )}
    </>
  );
}
