import { usePage } from '@inertiajs/react';
import { Alert, Anchor, Box, Card } from '@mantine/core';

export function Banner() {
  const { locale } = usePage().props;
  // const { t } = useTranslation();

  return (
    <>
      {locale === 'ru' && (
        <Box
          variant="gradient"
          bg="var(--app-cta-gradient)"
          p="xs"
          my="xl"
          ta="center"
        >
          <Anchor
            target="_blank"
            href="https://special.hexlet.io/vibecode_webinar?utm_source=codebasics&utm_medium=banner&utm_campaign=vibecode_webinar"
          >
            Бесплатный воркшоп по Claude Code →{' '}
            <b>Соберите своё первое AI-приложение</b>
          </Anchor>
        </Box>
      )}
    </>
  );
}
