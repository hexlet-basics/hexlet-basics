import { Container, Stack, Text, Title } from '@mantine/core';

import { useTranslation } from 'react-i18next';

import ApplicationLayout from '@/pages/layouts/ApplicationLayout';

type Props = {
  code: string;
  header: string;
  description: string;
};

export default function Show({ code, header, description }: Props) {
  const { t } = useTranslation();

  return (
    <ApplicationLayout>
      <Container ta="center" my="xl" py="xl">
        <Stack align="center" gap="md">
          <Text fw={700}>{code}</Text>
          <Title order={1}>{header}</Title>
          <Text size="lg">{description}</Text>
        </Stack>
      </Container>
    </ApplicationLayout>
  );
}
