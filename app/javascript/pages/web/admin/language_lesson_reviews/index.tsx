import { ActionIcon, Button, Modal, Select } from '@mantine/core';
import dayjs from 'dayjs';
import { Link, Search } from 'lucide-react';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/Elements/AppAnchor';
import MarkdownViewer from '@/components/MarkdownViewer';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Grid, LanguageLessonReview } from '@/types';

type Props = PropsWithChildren & {
  reviews: LanguageLessonReview[];
  language_slugs: Array<string>;
  grid: Grid & {
    fields: {
      language_slug_eq: string;
    };
  };
};

function DataBox({
  review,
  col,
}: {
  review: LanguageLessonReview;
  col: string;
}) {
  const [visible, setVisible] = useState(false);
  return (
    <>
      <Button variant="subtle" onClick={() => setVisible(true)}>
        data
      </Button>
      <Modal opened={visible} onClose={() => setVisible(false)} size="50vw">
        <MarkdownViewer>{review.summary}</MarkdownViewer>
      </Modal>
    </>
  );
}

export default function Index({ grid, reviews, language_slugs }: Props) {
  const { t } = useTranslation();
  const { gridProps, filters } = useDataTableProps<
    LanguageLessonReview,
    typeof grid.fields
  >(grid);

  const renderActions = (item: LanguageLessonReview) => {
    return (
      <>
        {/* <AppAnchor */}
        {/*   me="xs" */}
        {/*   href={Routes.edit_admin_language_category_path(item.id)} */}
        {/* > */}
        {/*   <ActionIcon variant="default" size="xs"> */}
        {/*     <Edit /> */}
        {/*   </ActionIcon> */}
        {/* </AppAnchor> */}
        <AppAnchor
          external
          href={Routes.language_lesson_path(item.language_slug!, item.slug!)}
        >
          <ActionIcon variant="default" size="xs">
            <Link />
          </ActionIcon>
        </AppAnchor>
        {/* <AppAnchor external href={item.source_code_url}> */}
        {/*   <ActionIcon variant="default" size="xs"> */}
        {/*     <Github /> */}
        {/*   </ActionIcon> */}
        {/* </AppAnchor> */}
        {/* <Link */}
        {/*   onClick={confirmDeleting} */}
        {/*   className="btn btn-link link-body-emphasis p-0 m-0" */}
        {/*   method="delete" */}
        {/*   href={Routes.admin_language_category_path(data.id)} */}
        {/* > */}
        {/*   {<i className="bi bi-file-x" />} */}
        {/* </Link> */}
      </>
    );
  };

  const filterLanguageSlug = (
    <Select
      data={language_slugs}
      value={filters.values.language_slug_eq}
      onChange={filters.getOnChange('language_slug_eq')}
      leftSection={<Search size={16} />}
      comboboxProps={{ withinPortal: false }}
      clearable
      searchable
    />
  );

  return (
    <AdminLayout header={t('admin.language_lesson_reviews.index.header')}>
      <DataTable
        records={reviews}
        columns={[
          { accessor: 'id' },
          { accessor: 'language_slug', filter: filterLanguageSlug },
          { accessor: 'slug' },
          {
            accessor: 'review',
            title: 'review',
            render: (rec) => <DataBox review={rec} col="review" />,
          },
          {
            accessor: 'created_at',
            render: (r) => dayjs(r.created_at).format('LL'),
          },
          {
            accessor: 'actions',
            title: 'actions',
            render: renderActions,
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
