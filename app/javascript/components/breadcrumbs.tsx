import { Head } from '@inertiajs/react';
import { Breadcrumbs, Center, Text } from '@mantine/core';
import { Home } from 'lucide-react';
import type { HTMLAttributes, PropsWithChildren, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import type { BreadcrumbList, ListItem, WithContext } from 'schema-dts';
import type { BreadcrumbItem } from '@/types';
import AppAnchor from './AppAnchor';

type Props = PropsWithChildren & {
  items: BreadcrumbItem[];
};

type BreadcrumbItemWithActive = {
  title: ReactNode;
  href: string;
  active?: boolean;
};

export function XBreadcrumb({ items = [] }: Props) {
  const { t } = useTranslation();

  const itemListElement: ListItem[] = items.map((item, index) => {
    return {
      position: index + 1,
      '@type': 'ListItem',
      item: {
        '@id': item.url,
        name: item.name,
      },
    };
  });

  const breadcrumbList: WithContext<BreadcrumbList> = {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    itemListElement,
  };

  const breadcrumbs: BreadcrumbItemWithActive[] = [
    {
      title: (
        <Center c="dimmed">
          <Home size={15} />
        </Center>
      ),
      href: '/',
    },
    ...items.map((item, index) => ({
      title: item.name,
      href: item.url,
      active: items.length === index + 1,
    })),
  ];

  return (
    <>
      <Head>
        <script type="application/ld+json">
          {JSON.stringify(breadcrumbList)}
        </script>
      </Head>
      <Breadcrumbs
        styles={{
          breadcrumb: { whiteSpace: 'normal' },
        }}
      >
        {breadcrumbs.map((item) =>
          item.active ? (
            <Text key={item.href} c="dimmed" size="sm">
              {item.title}
            </Text>
          ) : (
            <AppAnchor key={item.href} href={item.href} size="sm">
              {item.title}
            </AppAnchor>
          ),
        )}
      </Breadcrumbs>
    </>
  );
}
