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
            href="https://special.hexlet.io/vibecode_webinar?utm_source=codebasics&utm_medium=banner&utm_campaign=vibecode_webinar"
            {...propsForExternalLink()}
          >
            Бесплатный воркшоп по Claude Code →{" "}
            <b>Соберите своё первое AI-приложение</b>
          </Anchor>
        </Box>
      )}
    </>
  );
}
