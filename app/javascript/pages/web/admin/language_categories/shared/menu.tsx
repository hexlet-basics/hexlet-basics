import { LogIn } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from '@/components/CrudHorizontalMenu';
import * as Routes from '@/routes.js';
import type { LanguageCategoryCrud } from '@/types';

type Props = {
  data?: LanguageCategoryCrud;
};

export function Menu({ data }: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_language_categories_path(),
      label: tHelpers('crud.list'),
    },
    {
      href: Routes.new_admin_language_category_path(),
      label: tHelpers('crud.add'),
    },
  ];

  if (data) {
    items.push({
      href: Routes.edit_admin_language_category_path(data.language_category.id),
      label: tHelpers('crud.editing'),
    });
    items.push({
      href: Routes.language_category_path(data.language_category.slug!),
      external: true,
      label: <LogIn size={15} />,
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
