import { LogIn } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import {
  CrudHorizontalMenu,
  type CrudHorizontalMenuItem,
} from '@/components/CrudHorizontalMenu';
import * as Routes from '@/routes.js';
import type { BlogPostCrud } from '@/types';

type Props = {
  data?: BlogPostCrud;
};

export default function Menu({ data }: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const items: CrudHorizontalMenuItem[] = [
    { href: Routes.admin_blog_posts_path(), label: tHelpers('crud.list') },
    { href: Routes.new_admin_blog_post_path(), label: tHelpers('crud.add') },
  ];

  if (data?.blog_post) {
    items.push({
      href: Routes.edit_admin_blog_post_path(data.blog_post.id),
      label: tHelpers('crud.editing'),
    });
    items.push({
      href: Routes.blog_post_path(data.blog_post.slug!),
      external: true,
      label: <LogIn size={15} />,
    });
  }

  return <CrudHorizontalMenu items={items} />;
}
