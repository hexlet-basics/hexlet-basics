import { Button } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import type { HTTPVerb } from 'use-inertia-form';
import { XFile, XForm, XInput, XSelect, XTextarea } from '@/components/forms';
import type { BlogPostCrud } from '@/types/serializers';

type Props = {
  data: BlogPostCrud;
  url: string;
  method?: HTTPVerb;
};

const locales = [
  { name: 'Russian', code: 'ru' },
  { name: 'English', code: 'en' },
];

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  return (
    <XForm method={method} model="blog_post" data={data} to={url}>
      <XSelect
        field="state"
        items={data.meta.states}
        labelField="value"
        valueField="key"
      />
      <XInput field="name" />
      <XFile metaName="cover_thumb_variant_url" field="cover" />
      <XInput field="slug" />
      <XTextarea field="description" rows={5} />
      <XTextarea field="body" rows={12} />
      <Button type="submit" mt="xl">
        {tHelpers('submit.save')}
      </Button>
    </XForm>
  );
}
