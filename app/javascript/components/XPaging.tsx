import { Link, router, usePage } from '@inertiajs/react';
import { Center, Pagination } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { fromWindow } from '@/lib/utils.ts';
import type { Pagy, SharedProps } from '@/types';

type Props = PropsWithChildren & {
  pagy: Pagy;
};

const defaultInertiaLinkOptions = {
  replace: true,
  preserveState: true,
};

export default function XPaging({ pagy }: Props) {
  const { url } = usePage<SharedProps>();
  const origin = fromWindow('location')?.origin;

  const paginationUrl = (page: number) => {
    const u = new URL(url, origin);
    u.searchParams.set('page', String(page));
    return u.toString();
  };

  const linkComponent = (page: number) => {
    return {
      component: Link,
      href: paginationUrl(page),
      ...defaultInertiaLinkOptions,
    };
  };

  const controlUrl = (control: string) => {
    switch (control) {
      case 'next':
        return pagy.page < pagy.last
          ? {
              component: Link,
              href: paginationUrl(pagy.page + 1),
              ...defaultInertiaLinkOptions,
            }
          : {};
      case 'previous':
        return pagy.page > 1
          ? {
              component: Link,
              href: paginationUrl(pagy.page - 1),
              ...defaultInertiaLinkOptions,
            }
          : {};
      default:
        return {};
    }
  };

  return (
    <Center>
      <Pagination
        my="xl"
        total={pagy.last}
        value={pagy.page}
        getItemProps={(page) => linkComponent(page)}
        getControlProps={(control) => controlUrl(control)}
        boundaries={1}
        siblings={1}
      />
    </Center>
  );
}
