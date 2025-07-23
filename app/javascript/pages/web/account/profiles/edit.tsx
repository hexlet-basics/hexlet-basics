import { Box, Button, Card, Center, Container, TextInput } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/AppAnchor';
import { useAppForm } from '@/hooks/useAppForm';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { UserProfileForm } from '@/types/serializers';

type Props = PropsWithChildren & {
  form: UserProfileForm;
};

export default function Edit({ form }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<UserProfileForm>({
    url: Routes.account_profile_path(),
    method: 'patch',
    container: form,
  });

  return (
    <ApplicationLayout center header={t('account.profiles.edit.title')}>
      <Container mt="xl">
        <Center>
          <Card
            shadow="sm"
            withBorder
            p="xl"
            w={{ base: '100%', xs: '70%', sm: '50%' }}
          >
            <form onSubmit={submit}>
              <TextInput {...getInputProps('first_name')} autoComplete="name" />
              <TextInput {...getInputProps('last_name')} autoComplete="name" />
              <Button type="submit" fullWidth mt="xl" loading={isSubmitting}>
                {tHelpers('submit.save')}
              </Button>
            </form>

            <Box>
              <AppAnchor
                withConfirmation
                mt="xl"
                href={Routes.account_profile_path()}
                method="delete"
                c="red"
              >
                {t('account.profiles.edit.delete')}
              </AppAnchor>
            </Box>
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
