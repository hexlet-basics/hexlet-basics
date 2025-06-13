import type { BreadcrumbItem } from "@/types";
import { Head, Link } from "@inertiajs/react";
import type { HTMLAttributes, PropsWithChildren, ReactNode } from "react";
import { Breadcrumbs, Anchor, Text, Center } from '@mantine/core';
import { useTranslation } from "react-i18next";
import type { BreadcrumbList, ListItem, WithContext } from "schema-dts";
import { Home } from "lucide-react";

type Props = PropsWithChildren &
  HTMLAttributes<HTMLDivElement> & {
    items: BreadcrumbItem[];
  };

type BreadcrumbItemWithActive = {
  title: ReactNode;
  href: string;
  active?: boolean;
};

export function XBreadcrumb({ items = [], className }: Props) {
  const { t } = useTranslation();

  const itemListElement: ListItem[] = items.map((item, index) => {
    return {
      position: index + 1,
      "@type": "ListItem",
      item: {
        "@id": item.url,
        name: item.name,
      },
    };
  });

  const breadcrumbList: WithContext<BreadcrumbList> = {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    itemListElement,
  };

  const breadcrumbs: BreadcrumbItemWithActive[] = [
    { title: <Center c="dimmed"><Home size={15} /></Center>, href: "/" },
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
      <Breadcrumbs className={className}>
        {breadcrumbs.map((item) => (
          item.active ? (
            <Text key={item.href} c="dimmed" size="sm">
              {item.title}
            </Text>
          ) : (
            <Anchor
              key={item.href}
              component={Link}
              href={item.href}
              size="sm"
            >
              {item.title}
            </Anchor>
          )
        ))}
      </Breadcrumbs>
    </>
  );
}
