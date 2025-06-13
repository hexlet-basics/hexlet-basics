import { useTranslation } from "react-i18next";
import * as Routes from "@/routes.js";
import type ReviewCrud from "@/types/serializers/ReviewCrud";
import { LogIn } from "lucide-react";
import { CrudHorizontalMenu, CrudHorizontalMenuItem } from "@/components/CrudHorizontalMenu";

type Props = {
  data?: ReviewCrud;
};

export function Menu({ data }: Props) {
  const { t: tHelpers } = useTranslation("helpers");

  const items: CrudHorizontalMenuItem[] = [
    { href: Routes.admin_reviews_path(), label: tHelpers("crud.list") },
    { href: Routes.new_admin_review_path(), label: tHelpers("crud.add") },
  ];

  if (data) {
    items.push(
      { href: Routes.edit_admin_review_path(data.review.id), label: tHelpers("crud.editing") },
    );
    items.push(
      { href: Routes.reviews_path(), external: true, label: <LogIn size={15} /> },
    );
  }

  return <CrudHorizontalMenu items={items} />;
}
