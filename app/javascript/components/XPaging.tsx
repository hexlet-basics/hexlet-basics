import { router, usePage } from '@inertiajs/react';
import { Center, Pagination } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { fromWindow } from '@/lib/utils.ts';
import type { Pagy, SharedProps } from '@/types';

type Props = PropsWithChildren & {
  pagy: Pagy;
};

export default function XPaging({ pagy }: Props) {
  const { url } = usePage<SharedProps>();
  const origin = fromWindow('location')?.origin;

  const handleChange = (page: number) => {
    if (!origin) return;

    const u = new URL(url, origin);
    u.searchParams.set('page', String(page));

    router.visit(u.toString(), {
      replace: true,
      preserveState: true,
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
