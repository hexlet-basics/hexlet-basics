import { Button } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import type { HTTPVerb } from 'use-inertia-form';
import {
  XAutocomplete,
  XCheck,
  XForm,
  XInput,
  XSelect,
  XTextarea,
} from '@/components/forms';

import * as Routes from '@/routes.js';
import type { Language, ReviewCrud } from '@/types';

type Props = {
  data: ReviewCrud;
  url: string;
  courses: Language[];
  method?: HTTPVerb;
};

export default function Form({ courses, data, url, method }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  return (
    <XForm method={method} model="review" data={data} to={url}>
      <XSelect
        field="state"
        items={data.meta.states}
        labelField="value"
        valueField="key"
      />
      <XCheck field="pinned" />
      <XSelect
        field="language_id"
        labelField="slug"
        valueField="id"
        items={courses}
      />
      <XAutocomplete
        field="user_id"
        has="user"
        labelField="email"
        valueField="id"
        source={Routes.search_admin_api_users_path()}
      />
      <XInput field="first_name" autoComplete="name" />
      <XInput field="last_name" autoComplete="name" />
      <XTextarea field="body" rows={8} />
      <Button type="submit">{tHelpers('submit.save')}</Button>
    </XForm>
  );
}
