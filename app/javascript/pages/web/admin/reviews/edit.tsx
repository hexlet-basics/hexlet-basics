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

export default function Edit({ reviewDto, courses }: Props) {
  const { t } = useTranslation();

  return (
    <AdminLayout
      header={t('admin.reviews.edit.header', { id: reviewDto.review.id })}
    >
      <Menu data={reviewDto} />
      <Form
        courses={courses}
        method="patch"
        data={reviewDto}
        url={Routes.admin_review_path(reviewDto.review.id)}
      />
    </AdminLayout>
  );
}
