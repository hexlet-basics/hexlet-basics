import { Box, Button, TextInput } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import * as Routes from '@/routes.js';
import type { UserSignUpForm } from '@/types';
import AppAnchor from './Elements/AppAnchor';
import XssContent from './XssContent';

type Props = PropsWithChildren & {
  userDto: UserSignUpForm;
  autoFocus?: boolean;
};

export default function SignUpFormBlock({ userDto, autoFocus = false }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<UserSignUpForm>({
    url: Routes.users_path(),
    method: 'post',
    container: userDto,
  });

  return (
    <form onSubmit={submit}>
      <TextInput
        {...getInputProps('first_name')}
        autoComplete="name"
        autoFocus={autoFocus}
      />
      <TextInput {...getInputProps('email')} required autoComplete="email" />
      <TextInput
        {...getInputProps('password')}
        required
        type="password"
        autoComplete="current-password"
      />
      <Box my="lg" ta="right">
        {t('users.new.have_account')}{' '}
        <AppAnchor fw="bold" href={Routes.new_session_path()}>
          {t('users.new.sign_in')}
        </AppAnchor>
      </Box>
      <Button type="submit" fullWidth loading={isSubmitting}>
        {tHelpers('submit.user_sign_up_form.create')}
      </Button>

      <XssContent fz="sm" mt="xs">
        {t('users.new.confirmation_html', {
          url: Routes.page_path('tos'),
        })}
      </XssContent>
    </form>
  );
}
