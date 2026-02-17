import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { SurveyCrud } from "@/types";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  surveyDto: SurveyCrud;
};

export default function Edit({ surveyDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t(($) => $.admin.surveys.edit.header, {
        id: surveyDto.id,
      })}
    >
      <Menu data={surveyDto} />
      <Form
        method="patch"
        data={surveyDto}
        url={Routes.admin_survey_path(surveyDto.id)}
      />
    </AdminLayout>
  );
}
