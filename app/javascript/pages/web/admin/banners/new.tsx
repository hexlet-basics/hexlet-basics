import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { BannerCreate } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  bannerDto: BannerCreate;
};

export default function New({ bannerDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.banners.new.header)}>
      <Menu />
      <Form data={bannerDto} url={Routes.admin_banners_path()} />
    </AdminLayout>
  );
}
