import { DataTable } from 'mantine-datatable';
import type { PropsWithChildren } from 'react';

import * as Routes from '@/routes.js';
import { useTranslation } from 'react-i18next';

import AdminLayout from '@/pages/layouts/AdminLayout';
import AppAnchor from '@/components/AppAnchor';
import type { Lead, Grid } from '@/types/serializers';
import useDataTableProps from '@/hooks/useDataTableProps';
import { Modal, Button, Code } from '@mantine/core';
import { useState } from 'react';
import dayjs from 'dayjs';

type Props = PropsWithChildren & {
  leads: Lead[];
  grid: Grid;
};

function renderData(lead: Lead, col: string) {
  const [visible, setVisible] = useState(false);
  return (
    <>
      <Button variant="subtle" onClick={() => setVisible(true)}>
        data
      </Button>
      <Modal
        opened={visible}
        onClose={() => setVisible(false)}
        size="50vw"
      >
        <Code block>{JSON.stringify(lead[col as keyof Lead], null, 2)}</Code>
      </Modal>
    </>
  );
}

export default function Index({ grid, leads }: Props) {
  const { t } = useTranslation();
  const gridProps = useDataTableProps<Lead>(grid);

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
          { accessor: 'survey_answers_data', title: 'answers', render: (rec) => renderData(rec, 'survey_answers_data') },
          { accessor: 'courses_data', title: 'courses', render: (rec) => renderData(rec, 'courses_data') },
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

