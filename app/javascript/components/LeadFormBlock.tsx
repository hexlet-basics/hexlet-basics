import { Box, Button, Select, TextInput } from '@mantine/core';
import { toMerged } from 'es-toolkit';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import { enumToOptions, fromWindow } from '@/lib/utils';
import * as Routes from '@/routes.js';
import type { LeadCrud } from '@/types';
import XssContent from './XssContent';

type Props = {
  leadDto: LeadCrud;
  autoFocus?: boolean;
};

export default function LeadFormBlock({ leadDto, autoFocus = false }: Props) {
  const { t: tAr } = useTranslation('activerecord');
  const { t: tHelpers } = useTranslation('helpers');
  const { t: tViews } = useTranslation();

  const contactMethodEnum = tAr('attributes.user.contact_method/values', {
    returnObjects: true,
  });
  const contactMethodOptions = enumToOptions(contactMethodEnum);

  const preparedLeadDto = toMerged(leadDto, {
    data: {
      ym_client_id: fromWindow('ymClientId'),
      contact_method: 'telegram',
    },
  });

  const {
    getInputProps,
    getSelectProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<LeadCrud>({
    url: Routes.leads_path(),
    method: 'post',
    container: preparedLeadDto,
  });

  return (
    <form onSubmit={submit}>
      {/* Поле ym_client_id передаём скрыто */}
      <input type="hidden" {...getInputProps('ym_client_id')} />

      <Select
        {...getSelectProps(
          'contact_method',
          contactMethodOptions,
          'id',
          'name',
        )}
        required
      />

      <TextInput
        {...getInputProps('contact_value')}
        required
        autoFocus={autoFocus}
      />

      <Box fz="sm" my="lg">
        <XssContent>{tViews('blocks.lead_form_block.description1')}</XssContent>
        <XssContent>{tViews('blocks.lead_form_block.description2')}</XssContent>
      </Box>

      <Button type="submit" fullWidth loading={isSubmitting}>
        {tHelpers('send')}
      </Button>
    </form>
  );
}
