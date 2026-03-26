import { Link } from "@inertiajs/react";
import { ActionIcon, Group } from "@mantine/core";
import { IconBrandGithub, IconFileSearch, IconLink } from "@tabler/icons-react";
import dayjs from "dayjs";
import { DataTable } from "mantine-datatable";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import useDataTableProps from "@/hooks/useDataTableProps";
import AdminLayout from "@/layouts/AdminLayout";
import { propsForExternalLink } from "@/lib/utils";
import * as Routes from "@/routes.js";
import type { Grid, LanguageLesson } from "@/types";

type Props = PropsWithChildren & {
  lessons: LanguageLesson[];
  grid: Grid;
};

export default function Index({ grid, lessons }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<LanguageLesson, {}>(grid);

  const renderActions = (item: LanguageLesson) => {
    return (
      <Group gap="xs">
        {/* <AppAnchor */}
        {/*   me="xs" */}
        {/*   href={Routes.edit_admin_language_category_path(item.id)} */}
        {/* > */}
        {/*   <ActionIcon variant="default" size="xs"> */}
        {/*     <Edit /> */}
        {/*   </ActionIcon> */}
        {/* </AppAnchor> */}
        <a
          href={Routes.language_lesson_path(item.language!.slug!, item.slug!)}
          {...propsForExternalLink()}
        >
          <ActionIcon variant="default" size="xs">
            <IconLink />
          </ActionIcon>
        </a>
        <a href={item.source_code_url!} {...propsForExternalLink()}>
          <ActionIcon variant="default" size="xs">
            <IconBrandGithub />
          </ActionIcon>
        </a>
        <ActionIcon
          component={Link}
          href={Routes.review_admin_language_lesson_path(item.id)}
          method="post"
          variant="subtle"
        >
          <IconFileSearch size={14} />
        </ActionIcon>
      </Group>
    );
  };

  return (
    <AdminLayout header={t(($) => $.admin.language_lessons.index.header)}>
      <DataTable
        records={lessons}
        columns={[
          { accessor: "id" },
          { accessor: "language.slug" },
          { accessor: "slug" },
          {
            accessor: "created_at",
            render: (r) => dayjs(r.created_at).format("LL"),
          },
          {
            accessor: "actions",
            title: "actions",
            render: renderActions,
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
