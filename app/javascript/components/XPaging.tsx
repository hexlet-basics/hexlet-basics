import { router, usePage } from '@inertiajs/react';
import { Center, Pagination } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { fromWindow } from '@/lib/utils.ts';
import type { Pagy, SharedProps } from '@/types';

type Props = PropsWithChildren & {
  pagy: Pagy;
  only?: string[];
};

export default function XPaging({ pagy, only }: Props) {
  const { url } = usePage<SharedProps>();
  const origin = fromWindow('location')?.origin;
  const only_props = only && only.length > 0 ? [...only, 'pagy'] : undefined;

  const handleChange = (page: number) => {
    if (!origin) return;

    const u = new URL(url, origin);
    u.searchParams.set('page', String(page));

    router.visit(u.toString(), {
      replace: true,
      preserveState: true,
      only: only_props,
    });
  };

  return (
    <Center>
      <Pagination
        my="xl"
        total={pagy.last}
        value={pagy.page}
        onChange={handleChange}
        boundaries={1}
        siblings={1}
      />
    </Center>
  );
}
