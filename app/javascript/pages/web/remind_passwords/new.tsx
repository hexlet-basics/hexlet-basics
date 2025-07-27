import { Box, Button, Card, Center, Container, TextInput } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/Elements/AppAnchor';
import { useAppForm } from '@/hooks/useAppForm';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { PasswordReminderForm } from '@/types';

type Props = PropsWithChildren & {
  passwordReminder: PasswordReminderForm;
};

export default function New({ passwordReminder }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<PasswordReminderForm>({
    url: Routes.remind_password_path(),
    method: 'post',
    container: passwordReminder,
  });

  return (
    <ApplicationLayout center header={t('remind_passwords.new.title')}>
      <Container>
        <Center>
          <Card
            withBorder
            shadow="sm"
            p="xl"
            w={{ base: '100%', sm: '80%', md: '70%', lg: '50%' }}
          >
            <form onSubmit={submit}>
              <TextInput
                {...getInputProps('email')}
                required
                autoFocus
                autoComplete="email"
              />
              <Box my="lg" ta="right">
                {t('users.new.have_account')}{' '}
                <AppAnchor fw="bold" href={Routes.new_session_path()}>
                  {t('users.new.sign_in')}
                </AppAnchor>
              </Box>
              <Button fullWidth type="submit" loading={isSubmitting}>
                {tHelpers('submit.remind_password_form.create')}
              </Button>
            </form>
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
