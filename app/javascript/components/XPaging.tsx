import { router } from '@inertiajs/react';
import { Center, Pagination } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import type { Pagy } from '@/types';

type Props = PropsWithChildren & {
  pagy: Pagy;
  path: (page: number) => string;
  only?: string[];
};

export default function XPaging({ pagy, path, only }: Props) {
  const handleChange = (page: number) => {
    router.visit(path(page), {
      replace: true,
      preserveState: true,
      only,
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
