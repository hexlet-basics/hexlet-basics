import { Link, usePage } from '@inertiajs/react';
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

const paginationUrl = (page: number, url: string) => {
  const origin = fromWindow('location')?.origin;
  if (!origin) {
    return;
  }

  const u = new URL(url, origin);
  u.searchParams.set('page', String(page));
  return u.toString();
};

const linkComponent = (page: number, url: string) => {
  const linkHref = paginationUrl(page, url);

  return linkHref
    ? {
        component: Link,
        href: linkHref,
        ...defaultInertiaLinkOptions,
      }
    : {};
};

export default function XPaging({ pagy }: Props) {
  const { url } = usePage<SharedProps>();

  return (
    <Center>
      <Pagination
        my="xl"
        total={pagy.last}
        value={pagy.page}
        getItemProps={(page) => linkComponent(page, url)}
        getControlProps={(control) => {
          const page = pagy[control];
          return page ? linkComponent(page, url) : {};
        }}
        boundaries={1}
        siblings={1}
      />
    </Center>
  );
}
