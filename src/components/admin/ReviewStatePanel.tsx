import { useTranslation } from "react-i18next";
import {
  adminArchiveReviewMutation,
  adminGetReviewQueryKey,
  adminListReviewsQueryKey,
  adminPublishReviewMutation,
  adminUnpublishReviewMutation,
} from "@/client/@tanstack/react-query.gen";
import type { Review } from "@/client/types.gen";
import { PublicationStatePanel } from "@/components/admin/PublicationStatePanel";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import type { ReactElement } from "react";

// A review's Publication State and the buttons that change it (see
// PublicationStatePanel); this one binds the review's three operations. A review
// with no state counts as a draft, as on the server.
export function ReviewStatePanel({ review }: { review: Review }): ReactElement {
  const { t } = useTranslation();
  const options = {
    errorMessage: t(($) => $.admin.publicationState.changeError),
    invalidate: [adminListReviewsQueryKey(), adminGetReviewQueryKey({ path: { id: review.id } })],
    successMessage: t(($) => $.admin.publicationState.changed),
  };
  const publish = useResourceMutation({ mutation: adminPublishReviewMutation(), ...options });
  const unpublish = useResourceMutation({ mutation: adminUnpublishReviewMutation(), ...options });
  const archive = useResourceMutation({ mutation: adminArchiveReviewMutation(), ...options });
  const state = review.state ?? "draft";

  return (
    <PublicationStatePanel
      id={review.id}
      title={t(($) => $.models.attributes.review.state)}
      state={state}
      stateLabel={t(($) => $.models.attributes.review["state/values"][state])}
      mutations={{ archive, publish, unpublish }}
    />
  );
}
