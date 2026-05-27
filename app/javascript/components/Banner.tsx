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
            href="https://ru.hexlet.io/programs/vibecoding-claudecode?utm_source=codebasics&utm_medium=banner&utm_campaign=vibecoding_course"
          >
            Курс по вайбкодингу <b>со скидкой 17%</b> • только до 7 июня
          </Anchor>
        </Box>
      )}
    </>
  );
}
