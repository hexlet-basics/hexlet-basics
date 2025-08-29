import dayjs from 'dayjs';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import type { Grid, Survey, SurveyAnswer } from '@/types';
import { arrayToSelectData } from '@/lib/utils';
import { Search } from 'lucide-react';
import { Select } from '@mantine/core';

type Props = PropsWithChildren & {
  surveyAnswers: SurveyAnswer[];
  surveys: Survey[];
  grid: Grid & {
    fields: {
      survey_slug_eq: string;
    };
  };
};

export default function Index({ grid, surveyAnswers, surveys }: Props) {
  const { t } = useTranslation();
  const { gridProps, filters } = useDataTableProps<
    SurveyAnswer,
    typeof grid.fields
  >(grid);

  const surveyFilterSlug = (
    <Select
      data={arrayToSelectData(surveys, 'slug', "question")}
      value={filters.values.survey_slug_eq}
      onChange={filters.getOnChange('survey_slug_eq')}
      leftSection={<Search size={16} />}
      comboboxProps={{ withinPortal: false }}
      clearable
      searchable
    />
  );

  return (
    <AdminLayout header={t('admin.survey_answers.index.header')}>
      <DataTable
        records={surveyAnswers}
        columns={[
          { accessor: 'id' },
          { accessor: 'user_id' },
          { accessor: 'survey_slug', filter: surveyFilterSlug },
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
