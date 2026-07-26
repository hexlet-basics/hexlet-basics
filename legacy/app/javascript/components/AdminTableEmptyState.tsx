import { EmptyState } from "@mantine/core";
import { IconDatabaseOff } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";

export default function AdminTableEmptyState() {
  const { t } = useTranslation();

  return (
    <EmptyState icon={<IconDatabaseOff />} title={t(($) => $.common.empty)} />
  );
}
