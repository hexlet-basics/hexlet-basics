import { Center, Pagination } from '@mantine/core';
import type { PropsWithChildren } from 'react';
import type { Pagy } from '@/types';

type Props = PropsWithChildren & {
  pagy: Pagy;
};

export default function XPaging({ pagy }: Props) {
  return (
    <Center>
      <Pagination
        my="xl"
        total={pagy.last}
        value={pagy.page}
        boundaries={1}
        siblings={1}
      />
    </Center>
  );
}
