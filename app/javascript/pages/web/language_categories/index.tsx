import type { PropsWithChildren } from "react";
import { Container } from "react-bootstrap";

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { LanguageCategory } from "@/types/serializers";

import * as Routes from "@/routes.js";
import { Link } from "@inertiajs/react";

type Props = PropsWithChildren & {
  categories: LanguageCategory[];
};

export default function Index({ categories }: Props) {
  const { t } = useTranslation();
  const header = t("language_categories.index.header");

  const items = [
    {
      name: header,
      url: Routes.language_categories_path(),
    },
  ];

  return (
    <ApplicationLayout items={items} header={header}>
      <Container>
        <div className="row row-cols-1 row-cols-lg-2">
          {categories.map((category) => (
            <div className="col mb-3" key={category.id}>
              <Link className="display-5" href={Routes.language_category_path(category.slug!)}>
                {category.name}
              </Link>
            </div>
          ))}
        </div>
      </Container>
    </ApplicationLayout>
  );
}
