import { IconLogin2 } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from "@/components/CrudHorizontalMenu";
import * as Routes from "@/routes.js";
import type ReviewCrud from "@/types/serializers/ReviewCrud";

type Props = {
  data?: ReviewCrud;
};

export function Menu({ data }: Props) {
  const { t } = useTranslation();

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_reviews_path(),
      label: t(($) => $.helpers.crud.list),
    },
    {
      href: Routes.new_admin_review_path(),
      label: t(($) => $.helpers.crud.add),
    },
  ];

  if (data) {
    items.push({
      href: Routes.edit_admin_review_path(data.id),
      label: t(($) => $.helpers.crud.editing),
    });
    items.push({
      href: Routes.reviews_path(),
      external: true,
      label: <IconLogin2 size={15} />,
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
