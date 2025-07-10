import { LogIn } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from '@/components/CrudHorizontalMenu';
import * as Routes from '@/routes.js';
import type { LanguageLandingPageCrud } from '@/types/serializers';

type Props = {
  data?: LanguageLandingPageCrud;
};

export function Menu({ data }: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_language_landing_pages_path(),
      label: tHelpers('crud.list'),
    },
    {
      href: Routes.new_admin_language_landing_page_path(),
      label: tHelpers('crud.add'),
    },
  ];

  if (data) {
    items.push({
      href: Routes.edit_admin_language_landing_page_path(
        data.language_landing_page.id,
      ),
      label: tHelpers('crud.editing'),
    });
    items.push({
      href: Routes.language_path(data.language_landing_page.slug!),
      external: true,
      label: <LogIn size={15} />,
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
