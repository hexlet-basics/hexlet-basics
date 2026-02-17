import { IconBrandGithub, IconLogin2 } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from "@/components/CrudHorizontalMenu";
import * as Routes from "@/routes.js";
import type { LanguageCrud, LanguageLandingPage } from "@/types/serializers";

type Props = {
  data?: LanguageCrud;
  landingPage?: LanguageLandingPage | null;
};

export function Menu({ data, landingPage }: Props) {
  const { t } = useTranslation();

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_languages_path(),
      label: t(($) => $.helpers.crud.list),
    },
    {
      href: Routes.new_admin_language_path(),
      label: t(($) => $.helpers.crud.add),
    },
  ];

  if (data) {
    const meta = (data as { meta?: { repository_url?: string } }).meta ?? {};
    items.push({
      href: Routes.edit_admin_language_path(data.id),
      label: t(($) => $.helpers.crud.editing),
    });
    if (landingPage) {
      items.push({
        href: Routes.language_path(landingPage.slug),
        external: true,
        label: <IconLogin2 size={15} />,
      });
    }
    items.push({
      href: meta.repository_url ?? "",
      external: true,
      label: <IconBrandGithub size={15} />,
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
