import { Button, Code, Modal } from '@mantine/core';
import dayjs from 'dayjs';
import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import useDataTableProps from '@/hooks/useDataTableProps';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Grid, Lead } from '@/types';

type Props = PropsWithChildren & {
  leads: Lead[];
  grid: Grid;
};

function DataBox({ lead, col }: { lead: Lead; col: string }) {
  const [visible, setVisible] = useState(false);
  return (
    <>
      <Button variant="subtle" onClick={() => setVisible(true)}>
        data
      </Button>
      <Modal opened={visible} onClose={() => setVisible(false)} size="50vw">
        <Code block>{JSON.stringify(lead[col as keyof Lead], null, 2)}</Code>
      </Modal>
    </>
  );
}

export default function Index({ grid, leads }: Props) {
  const { t } = useTranslation();
  const { gridProps } = useDataTableProps<Lead, {}>(grid);

  return (
    <AdminLayout header={t('admin.leads.index.header')}>
      <DataTable
        records={leads}
        columns={[
          { accessor: 'id' },
          { accessor: 'email' },
          { accessor: 'user_id' },
          { accessor: 'full_name' },
          { accessor: 'phone' },
          { accessor: 'telegram' },
          { accessor: 'whatsapp' },
          {
            accessor: 'survey_answers_data',
            title: 'answers',
            render: (rec) => <DataBox lead={rec} col="survey_answers_data" />,
          },
          {
            accessor: 'courses_data',
            title: 'courses',
            render: (rec) => <DataBox lead={rec} col="course_data" />,
          },
          {
            accessor: 'created_at',
            render: (r) => dayjs(r.created_at).format('LL'),
            sortable: true,
          },
        ]}
        {...gridProps}
      />
    </AdminLayout>
  );
}
