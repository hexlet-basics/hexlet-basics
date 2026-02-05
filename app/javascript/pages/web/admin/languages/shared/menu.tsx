import { LogIn } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from '@/components/CrudHorizontalMenu';
import * as Routes from '@/routes.js';
import type { LanguageCrud, LanguageLandingPage } from '@/types/serializers';

type Props = {
  data?: LanguageCrud;
  landingPage?: LanguageLandingPage | null;
};

export function Menu({ data, landingPage }: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_languages_path(),
      label: tHelpers(($) => $.crud.list),
    },
    {
      href: Routes.new_admin_language_path(),
      label: tHelpers(($) => $.crud.add),
    },
  ];

  if (data) {
    const meta = (data as { meta?: { repository_url?: string } }).meta ?? {};
    items.push({
      href: Routes.edit_admin_language_path(data.data.id),
      label: tHelpers(($) => $.crud.editing),
    });
    if (landingPage) {
      items.push({
        href: Routes.language_path(landingPage.slug),
        external: true,
        label: <LogIn size={15} />,
      });
    }
    items.push({
      href: meta.repository_url ?? '',
      external: true,
      label: <i className="bi bi-github" />,
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
