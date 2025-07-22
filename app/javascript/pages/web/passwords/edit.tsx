import {
  Box,
  Button,
  Card,
  Container,
  Stack,
  TextInput,
  Title,
} from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { UserPassword } from '@/types/serializers';

type Props = PropsWithChildren & {
  userPassword: UserPassword;
};

export default function New({ userPassword }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<UserPassword>({
    url: Routes.password_path(),
    method: 'patch',
    container: userPassword,
  });

  return (
    <ApplicationLayout>
      <Container>
        <Stack align="center" gap="md">
          <Title order={1} ta="center" mb="md">
            {t('passwords.edit.title')}
          </Title>
          <Card
            withBorder
            p="xl"
            w={{ base: '100%', sm: '80%', md: '70%', lg: '50%' }}
          >
            <form onSubmit={submit}>
              <TextInput
                {...getInputProps('password')}
                type="password"
                autoComplete="new-password"
                required
              />
              <Box mt="lg" ta="right">
                <Button type="submit" loading={isSubmitting}>
                  {tHelpers('submit.replace')}
                </Button>
              </Box>
            </form>
          </Card>
        </Stack>
      </Container>
    </ApplicationLayout>
  );
}
