import type { BreadcrumbItem } from "@/types";
import type { HTMLAttributes, PropsWithChildren } from "react";
import { Breadcrumb, type BreadcrumbProps } from "react-bootstrap";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren &
  HTMLAttributes<BreadcrumbProps> & {
    items: BreadcrumbItem[];
  };

export function XBreadcrumb({ items = [], className }: Props) {
  const { t } = useTranslation();
  return (
    <Breadcrumb className={className}>
      <Breadcrumb.Item href="/" title={t("languages.show.to_home_title")}>
        <i className="bi bi-house" />
      </Breadcrumb.Item>
      {items.map((item, index) => (
        <Breadcrumb.Item
          key={item.url}
          href={item.url}
          active={items.length === index + 1}
        >
          {item.name}
        </Breadcrumb.Item>
      ))}
    </Breadcrumb>
  );
}
