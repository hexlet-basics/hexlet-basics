import { useTranslation } from "react-i18next";
import AdminLayout from "@/layouts/AdminLayout";
import * as Routes from "@/routes.js";
import type { BannerUpdate } from "@/types/serializers";
import Form from "./shared/form";
import { Menu } from "./shared/menu";

type Props = {
  bannerDto: BannerUpdate;
};

export default function Edit({ bannerDto }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t(($) => $.admin.banners.edit.header)}>
      <Menu data={bannerDto} />
      <Form
        method="patch"
        data={bannerDto}
        url={Routes.admin_banner_path(bannerDto.id)}
      />
    </AdminLayout>
  );
}
