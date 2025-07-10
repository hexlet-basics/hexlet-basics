import dayjs from 'dayjs';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import type { Grid, SurveyAnswer } from '@/types';

type Props = PropsWithChildren & {
  surveyAnswers: SurveyAnswer[];
  grid: Grid;
};

export default function Index({ grid, surveyAnswers }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<SurveyAnswer, {}>(grid);

  return (
    <AdminLayout header={t('admin.survey_answers.index.header')}>
      <DataTable
        records={surveyAnswers}
        columns={[
          { accessor: 'id' },
          { accessor: 'user_id' },
          { accessor: 'survey_slug' },
          { accessor: 'survey_item_value' },
          {
            accessor: 'created_at',
            sortable: true,
            render: (r) => dayjs(r.created_at).format('LL'),
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
