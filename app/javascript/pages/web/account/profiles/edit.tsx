import { Link } from '@inertiajs/react';
import { Anchor, Box, Button, Card, Center, Container } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/AppAnchor';
import { XForm, XInput } from '@/components/forms';
import useConfirmation from '@/hooks/useConfirmation';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { UserProfileForm } from '@/types/serializers';

type Props = PropsWithChildren & {
  form: UserProfileForm;
};

export default function Edit(props: Props) {
  const { form } = props;
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');
  const { t: tAr } = useTranslation('activerecord');

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
            <XForm
              model="user"
              method="patch"
              data={{ user: form }}
              to={Routes.account_profile_path()}
            >
              <XInput field="first_name" autoComplete="name" />
              <XInput field="last_name" autoComplete="name" />
              <Button type="submit" fullWidth mt="xl">
                {tHelpers('submit.save')}
              </Button>
            </XForm>

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
