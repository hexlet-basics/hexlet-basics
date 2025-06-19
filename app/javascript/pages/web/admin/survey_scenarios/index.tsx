import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import * as Routes from '@/routes.js';
import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import AppAnchor from '@/components/AppAnchor';
import type { Grid, SurveyScenario } from '@/types';
import { Menu } from './shared/menu';
import useDataTableProps from '@/hooks/useDataTableProps';
import { Edit } from 'lucide-react';
import dayjs from 'dayjs';

type Props = PropsWithChildren & {
  surveyScenarios: SurveyScenario[];
  grid: Grid;
};

export default function Index({ grid, surveyScenarios }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<SurveyScenario, {}>(grid);
  console.log(surveyScenarios)

  const renderActions = (item: SurveyScenario) => (
    <AppAnchor href={Routes.edit_admin_survey_scenario_path(item.id)}>
      <Edit size={14} />
    </AppAnchor>
  );

  return (
    <AdminLayout header={t('admin.survey_scenarios.index.header')}>
      <Menu />
      <DataTable
        records={surveyScenarios}
        columns={[
          { accessor: 'id' },
          { accessor: 'name' },
          { accessor: 'survey_item_value', title: 'Parent Item' },
          {
            accessor: 'created_at',
            sortable: true,
            render: (r) => dayjs(r.created_at).format('LL'),
          },
          { accessor: 'actions', title: 'actions', render: renderActions },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}

