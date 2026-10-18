import { Alert, Button, Card, Container, Group, List, Stack, Text, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import {
  approveOAuthConsentMutation,
  denyOAuthConsentMutation,
  getOAuthConsentOptions,
} from "@/client/@tanstack/react-query.gen";
import { requireAuth } from "@/lib/auth";
import { lenient, orNull } from "@/lib/lenient";

// The OAuth consent screen (ADR-0022): an assistant (an MCP client) was sent
// here by the authorization endpoint with the id of its request. A guest signs
// in first and comes back, because `requireAuth` carries this URL, request id
// included, through the sign-in page. The screen shows who asks, where the
// browser returns to and what is asked, and the decision is one of two contract
// operations whose answer is the next place for the browser to go.
const connectSearchSchema = z.object({
  id: lenient(z.string()),
});

export const Route = createFileRoute("/{-$locale}/connect")({
  validateSearch: connectSearchSchema,

  loaderDeps: ({ search }) => ({ id: search.id }),

  beforeLoad: ({ context, location }) => {
    requireAuth(context.user, location.href);
  },

  // An unknown, expired or already decided request is a normal outcome (a stale
  // tab), not a crash: the page says so instead of showing the error screen.
  loader: ({ context, deps }) =>
    deps.id
      ? orNull(
          context.queryClient.query({
            ...getOAuthConsentOptions({ path: { requestId: deps.id } }),
            staleTime: 0,
          }),
        )
      : null,

  head: ({ match }) => {
    const { t } = match.context.i18n;
    return {
      meta: [
        { title: t(($) => $.connect.title) },
        { content: t(($) => $.connect.meta.description), name: "description" },
        // A page with a one-time request id in its URL has nothing to index.
        { content: "noindex", name: "robots" },
      ],
    };
  },

  component: Connect,
});

// What each scope the server can offer means to a person. Scope names contain a
// colon, which i18next reads as a namespace separator, so the catalog uses
// underscored keys; a scope without an entry is shown as it is.
function useScopeLabel() {
  const { t } = useTranslation();
  const labels: Record<string, string> = {
    admin: t(($) => $.connect.scopes.admin),
    "catalog:read": t(($) => $.connect.scopes.catalog_read),
    offline_access: t(($) => $.connect.scopes.offline_access),
    openid: t(($) => $.connect.scopes.openid),
    "progress:read": t(($) => $.connect.scopes.progress_read),
  };
  return (scope: string) => labels[scope] ?? scope;
}

// A redirect that stays on this computer (RFC 8252 native apps, or a page of a
// developer) is legitimate, but it is also what a page that wants your approval
// to land on a program running locally looks like, so the screen says so. The
// host arrives with its port, and an IPv6 literal is bracketed.
function isLocalHost(host: string): boolean {
  const name = host
    .replace(/:\d+$/u, "")
    .replaceAll(/^\[|\]$/gu, "")
    .toLowerCase();
  return name === "localhost" || name === "::1" || /^127(?:\.\d{1,3}){3}$/u.test(name);
}

// The next place is outside the router: the authorization endpoint's callback
// on this origin, or the client's own redirect URI. Leaving takes a full
// navigation.
const leave = ({ redirectUrl }: { redirectUrl: string }) => globalThis.location.assign(redirectUrl);

function Connect() {
  const { t } = useTranslation();
  const { id } = Route.useSearch();
  const consent = Route.useLoaderData();
  const scopeLabel = useScopeLabel();

  const onError = () => notifications.show({ color: "red", message: t(($) => $.connect.error) });
  const approve = useMutation({ ...approveOAuthConsentMutation(), onError, onSuccess: leave });
  const deny = useMutation({ ...denyOAuthConsentMutation(), onError, onSuccess: leave });
  const deciding = approve.isPending || deny.isPending || approve.isSuccess || deny.isSuccess;

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.connect.title)}
        </Title>

        <Card withBorder p="xl" w={{ base: "100%", lg: "50%", md: "70%", sm: "80%" }}>
          {id && consent ? (
            <Stack>
              <Text>{t(($) => $.connect.intro, { client: consent.clientName })}</Text>
              <Text fw={500}>{t(($) => $.connect.scopes_title)}</Text>
              <List>
                {consent.scopes.map((scope) => (
                  <List.Item key={scope}>{scopeLabel(scope)}</List.Item>
                ))}
              </List>
              <Text size="sm" c="dimmed">
                {t(($) => $.connect.redirect, { host: consent.redirectHost })}
              </Text>
              {isLocalHost(consent.redirectHost) && (
                <Alert color="yellow" role="alert">
                  {t(($) => $.connect.localhost_warning, { host: consent.redirectHost })}
                </Alert>
              )}
              <Text size="sm" c="dimmed">
                {t(($) => $.connect.account, { account: consent.accountName })}
              </Text>
              <Group justify="flex-end">
                <Button
                  variant="default"
                  disabled={deciding}
                  onClick={() => deny.mutate({ path: { requestId: id } })}
                >
                  {t(($) => $.connect.deny)}
                </Button>
                <Button
                  loading={approve.isPending}
                  disabled={deciding}
                  onClick={() => approve.mutate({ path: { requestId: id } })}
                >
                  {t(($) => $.connect.approve)}
                </Button>
              </Group>
            </Stack>
          ) : (
            <Alert color="red">{t(($) => $.connect.not_found)}</Alert>
          )}
        </Card>
      </Stack>
    </Container>
  );
}
