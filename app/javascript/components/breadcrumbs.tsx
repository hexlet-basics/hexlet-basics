import type { BreadcrumbItem } from "@/types";
import { Head, Link } from "@inertiajs/react";
import type { HTMLAttributes, PropsWithChildren } from "react";
import { Breadcrumb, type BreadcrumbProps } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import type { BreadcrumbList, ListItem, WithContext } from "schema-dts";

type Props = PropsWithChildren &
  HTMLAttributes<BreadcrumbProps> & {
    items: BreadcrumbItem[];
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

  return (
    <>
      <Head>
        <script type="application/ld+json">
          {JSON.stringify(breadcrumbList)}
        </script>
      </Head>
      <Breadcrumb className={className}>
        <Breadcrumb.Item
          as={Link}
          title={t("languages.show.to_home_title")}
        >
          <Link href="/">
          <i className="bi bi-house" />
          </Link>
        </Breadcrumb.Item>
        {items.map((item, index) => (
          <Breadcrumb.Item
            key={item.url}
            active={items.length === index + 1}
          >
            <Link href={item.url}>
              {item.name}
            </Link>
          </Breadcrumb.Item>
        ))}
      </Breadcrumb>
    </>
  );
}
