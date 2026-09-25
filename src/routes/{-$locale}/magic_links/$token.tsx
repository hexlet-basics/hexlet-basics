import { createFileRoute, redirect } from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { getCurrentUserQueryKey } from "@/client/@tanstack/react-query.gen";
import { consumeMagicLink } from "@/client/sdk.gen";

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
      // The response is the signed-in user, so the cache is set from it rather
      // than refetched: the root beforeLoad reads it through ensureQueryData,
      // which would otherwise keep serving the guest it resolved on arrival.
      context.queryClient.setQueryData(getCurrentUserQueryKey(), { user });
    } catch (error) {
      if (isAxiosError(error) && error.response?.status === 404) {
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
