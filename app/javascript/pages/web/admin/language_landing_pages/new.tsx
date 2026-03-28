import { Grid } from "@mantine/core";
import { useTranslation } from "react-i18next";

import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type {
  Language,
  LanguageLandingPage,
  LanguageLandingPageCrud,
} from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  landingPageDto: LanguageLandingPageCrud;
  landingPages: LanguageLandingPage[];
  languages: Language[];
};

export default function New({
  landingPageDto,
  landingPages,
  languages,
}: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.language_landing_pages.new.header)}>
      <Menu />
      <Grid>
        <Grid.Col span={7}>
          <Form
            data={landingPageDto}
            landingPages={landingPages}
            languages={languages}
            url={Routes.admin_language_landing_pages_path()}
          />
        </Grid.Col>
      </Grid>
    </AdminLayout>
  );
}
