import { useTranslation } from "react-i18next";
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from "@/components/CrudHorizontalMenu";
import * as Routes from "@/routes.js";
import type BannerUpdate from "@/types/serializers/BannerUpdate";

type Props = {
  data?: BannerUpdate;
};

export function Menu({ data }: Props) {
  const { t } = useTranslation();

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_banners_path(),
      label: t(($) => $.helpers.crud.list),
    },
    {
      href: Routes.new_admin_banner_path(),
      label: t(($) => $.helpers.crud.add),
    },
  ];

  if (data) {
    items.push({
      href: Routes.edit_admin_banner_path(data.id),
      label: t(($) => $.helpers.crud.editing),
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
