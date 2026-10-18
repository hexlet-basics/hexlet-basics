import { useTranslation } from "react-i18next";
import {
  adminArchiveCourseLandingPageMutation,
  adminGetCourseLandingPageQueryKey,
  adminListCourseLandingPagesQueryKey,
  adminPublishCourseLandingPageMutation,
  adminUnpublishCourseLandingPageMutation,
} from "@/client/@tanstack/react-query.gen";
import type { CourseLandingPage } from "@/client/types.gen";
import { PublicationStatePanel } from "@/components/admin/PublicationStatePanel";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// A landing page's Publication State and the buttons that change it (see
// PublicationStatePanel); this one binds the page's three operations. A page
// with no state counts as a draft, as on the server.
export function CourseLandingPageStatePanel({ page }: { page: CourseLandingPage }) {
  const { t } = useTranslation();
  const options = {
    errorMessage: t(($) => $.admin.publicationState.changeError),
    invalidate: [
      adminListCourseLandingPagesQueryKey(),
      adminGetCourseLandingPageQueryKey({ path: { id: page.id } }),
    ],
    successMessage: t(($) => $.admin.publicationState.changed),
  };
  const publish = useResourceMutation({
    mutation: adminPublishCourseLandingPageMutation(),
    ...options,
  });
  const unpublish = useResourceMutation({
    mutation: adminUnpublishCourseLandingPageMutation(),
    ...options,
  });
  const archive = useResourceMutation({
    mutation: adminArchiveCourseLandingPageMutation(),
    ...options,
  });
  const state = page.state ?? "draft";

  return (
    <PublicationStatePanel
      id={page.id}
      title={t(($) => $.models.attributes.course_landing_page.state)}
      state={state}
      stateLabel={t(($) => $.models.attributes.course_landing_page["state/values"][state])}
      mutations={{ archive, publish, unpublish }}
    />
  );
}
