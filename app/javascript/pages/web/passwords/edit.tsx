import { Link } from '@inertiajs/react';
import { Button, Card, Container, Stack, Title } from '@mantine/core';
import type { PropsWithChildren } from 'react';

import { useTranslation } from 'react-i18next';
import { Submit } from 'use-inertia-form';

import { XForm, XInput } from '@/components/forms';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { UserPassword } from '@/types/serializers';

type Props = PropsWithChildren & {
  userPassword: UserPassword;
};

export default function New({ userPassword }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

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
            <XForm
              method="patch"
              model="user_password_form"
              data={{ user_password_form: userPassword }}
              to={Routes.password_path()}
            >
              <XInput
                field="password"
                type="password"
                autoComplete="new-password"
              />
              <Button type="submit">{tHelpers('submit.replace')}</Button>
            </XForm>
          </Card>
        </Stack>
      </Container>
    </ApplicationLayout>
  );
}
