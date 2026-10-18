import { notifications } from "@mantine/notifications";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { consumeMagicLink } from "@/client/sdk.gen";
import { isApiError } from "@/lib/api-error";
import { onIdentityChanged } from "@/lib/identity";
import { HTTP_NOT_FOUND } from "@/lib/http-status";

// The emailed Magic Link (legacy magic_links#show). Following it signs the
// visitor in and takes them home; a refused link — expired, forged, outlived by
// an email change, or issued by the Rails app — sends them back to the request
// form, which says the link is no longer valid.
//
// `ssr: false` keeps the call in the browser: the API answers with the session
// cookie, and a server-side fetch would keep that cookie in the Node process
// instead of handing it to the visitor.
export const Route = createFileRoute("/{-$locale}/magic_links/$token")({
  ssr: false,
  loader: async ({ context, params }) => {
    try {
      const { data: user } = await consumeMagicLink({
        path: { token: params.token },
        throwOnError: true,
      });
      // The response is the signed-in user. The redirect below is itself the
      // navigation that reads every page afresh, so there is no router to reload.
      await onIdentityChanged(context.queryClient, null, user);
      // Legacy had no message of its own for this one and fell back to the
      // application's generic success; the loader runs in the browser only.
      notifications.show({
        color: "green",
        message: context.i18n.t(($) => $.flash.application.base.success),
      });
    } catch (error) {
      if (isApiError(error, HTTP_NOT_FOUND)) {
        throw redirect({
          to: "/{-$locale}/magic_links/new",
          params: { locale: params.locale },
          search: { invalid: true },
        });
      }
      throw error;
    }
    throw redirect({ to: "/{-$locale}", params: { locale: params.locale } });
  },
});
