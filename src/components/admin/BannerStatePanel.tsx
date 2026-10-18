import { useTranslation } from "react-i18next";
import {
  adminArchiveBannerMutation,
  adminGetBannerQueryKey,
  adminListBannersQueryKey,
  adminPublishBannerMutation,
  adminUnpublishBannerMutation,
} from "@/client/@tanstack/react-query.gen";
import type { Banner } from "@/client/types.gen";
import { PublicationStatePanel } from "@/components/admin/PublicationStatePanel";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// A banner's Publication State and the buttons that change it (see
// PublicationStatePanel); this one binds the banner's three operations.
export function BannerStatePanel({ banner }: { banner: Banner }) {
  const { t } = useTranslation();
  const options = {
    errorMessage: t(($) => $.admin.publicationState.changeError),
    invalidate: [adminListBannersQueryKey(), adminGetBannerQueryKey({ path: { id: banner.id } })],
    successMessage: t(($) => $.admin.publicationState.changed),
  };
  const publish = useResourceMutation({ mutation: adminPublishBannerMutation(), ...options });
  const unpublish = useResourceMutation({ mutation: adminUnpublishBannerMutation(), ...options });
  const archive = useResourceMutation({ mutation: adminArchiveBannerMutation(), ...options });

  return (
    <PublicationStatePanel
      id={banner.id}
      title={t(($) => $.models.attributes.banner.state)}
      state={banner.state}
      stateLabel={t(($) => $.models.attributes.banner["state/values"][banner.state])}
      mutations={{ archive, publish, unpublish }}
    />
  );
}
