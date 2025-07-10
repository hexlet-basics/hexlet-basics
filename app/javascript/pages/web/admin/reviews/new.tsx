import { useTranslation } from 'react-i18next';
import AdminLayout from '@/pages/layouts/AdminLayout';
import * as Routes from '@/routes.js';
import type { Language, Review } from '@/types/serializers';
import type ReviewCrud from '@/types/serializers/ReviewCrud';
import Form from './shared/form';
import { Menu } from './shared/menu';

type Props = {
  reviewDto: ReviewCrud;
  courses: Language[];
};

export default function New({ reviewDto, courses }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout header={t('admin.reviews.new.header')}>
      <Menu />
      <Form
        courses={courses}
        data={reviewDto}
        url={Routes.admin_reviews_path()}
      />
    </AdminLayout>
  );
}
