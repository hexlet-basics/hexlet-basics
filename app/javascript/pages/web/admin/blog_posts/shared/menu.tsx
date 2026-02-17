import { ActionIcon } from "@mantine/core";
import { IconLogin2, IconSchool } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from "@/components/CrudHorizontalMenu";
import * as Routes from "@/routes.js";
import type { BlogPostCrud } from "@/types";

type Props = {
  data?: BlogPostCrud;
};

export default function Menu({ data }: Props) {
  const { t } = useTranslation();

  const items: CrudHorizontalMenuItem[] = [
    {
      href: Routes.admin_blog_posts_path(),
      label: t(($) => $.helpers.crud.list),
    },
    {
      href: Routes.new_admin_blog_post_path(),
      label: t(($) => $.helpers.crud.add),
    },
  ];

  if (data) {
    items.push({
      href: Routes.edit_admin_blog_post_path(data.id),
      label: t(($) => $.helpers.crud.editing),
    });
    items.push({
      href: Routes.related_courses_admin_blog_post_path(data.id),
      method: "post",
      label: <IconSchool size={15} />,
    });
    items.push({
      href: Routes.blog_post_path(data.slug!),
      external: true,
      label: <IconLogin2 size={15} />,
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
