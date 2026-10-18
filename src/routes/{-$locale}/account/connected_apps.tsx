import { Alert, Button, Card, Container, Group, List, Stack, Text, Title } from "@mantine/core";
import { modals } from "@mantine/modals";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  listConnectedAppsOptions,
  listConnectedAppsQueryKey,
  revokeConnectedAppMutation,
} from "@/client/@tanstack/react-query.gen";
import type { ConnectedApp } from "@/client/types.gen";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import { requireAuth } from "@/lib/auth";

// The apps the signed-in user has connected through the OAuth server (ADR-0022,
// #848), each with what it may do, and a way to disconnect one. The list is
// always fresh (staleTime 0): a revoke elsewhere or a lapse must not linger on
// screen. Access lapses 30 days after an app's last refresh, which the intro
// says in plain words.
export const Route = createFileRoute("/{-$locale}/account/connected_apps")({
  beforeLoad: ({ context, location }) => {
    requireAuth(context.user, location.href);
  },
  loader: ({ context }) =>
    context.queryClient.query({ ...listConnectedAppsOptions(), staleTime: 0 }),
  head: ({ match }) => {
    const { t } = match.context.i18n;
    return {
      meta: [
        { title: t(($) => $.account.connected_apps.title) },
        { name: "description", content: t(($) => $.account.connected_apps.meta.description) },
        { name: "robots", content: "noindex" },
      ],
    };
  },
  component: ConnectedApps,
});

// What each scope means to a person; same wording as the consent screen. A scope
// without an entry is shown as it is.
function useScopeLabel() {
  const { t } = useTranslation();
  const labels: Record<string, string> = {
    openid: t(($) => $.connect.scopes.openid),
    offline_access: t(($) => $.connect.scopes.offline_access),
    "catalog:read": t(($) => $.connect.scopes.catalog_read),
    "progress:read": t(($) => $.connect.scopes.progress_read),
    admin: t(($) => $.connect.scopes.admin),
  };
  return (scope: string) => labels[scope] ?? scope;
}

function ConnectedApps() {
  const { t, i18n } = useTranslation();
  // The loader filled the cache; the query keeps it current after a revoke.
  const { data } = useQuery({ ...listConnectedAppsOptions(), staleTime: 0 });
  const items = data?.items ?? [];
  const scopeLabel = useScopeLabel();

  const revoke = useResourceMutation({
    mutation: revokeConnectedAppMutation(),
    invalidate: [listConnectedAppsQueryKey()],
    successMessage: t(($) => $.account.connected_apps.revoked),
    errorMessage: t(($) => $.account.connected_apps.error),
  });

  // Disconnecting is immediate and cannot be undone by the user (the app has to
  // be connected again), so it sits behind the same yes/no confirmation as
  // deleting an account.
  const confirm = (app: ConnectedApp) =>
    modals.openConfirmModal({
      centered: true,
      title: t(($) => $.account.connected_apps.confirm_title, { client: app.name }),
      children: <Text size="sm">{t(($) => $.account.connected_apps.confirm_text)}</Text>,
      labels: {
        confirm: t(($) => $.account.connected_apps.confirm_label),
        cancel: t(($) => $.common.boolean.no),
      },
      confirmProps: { color: "red" },
      onConfirm: () => revoke.mutate({ path: { id: app.id } }),
    });

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.account.connected_apps.title)}
        </Title>

        <Stack w={{ base: "100%", sm: "80%", md: "70%" }}>
          <Text c="dimmed">{t(($) => $.account.connected_apps.intro)}</Text>

          {items.length === 0 ? (
            <Alert color="gray">{t(($) => $.account.connected_apps.empty)}</Alert>
          ) : (
            items.map((app) => (
              <Card key={app.id} withBorder p="lg">
                <Group justify="space-between" align="flex-start" wrap="nowrap">
                  <Stack gap="xs">
                    <Text fw={600}>{app.name}</Text>
                    {app.clientHost && (
                      <Text size="sm" c="dimmed">
                        {t(($) => $.account.connected_apps.host, { host: app.clientHost })}
                      </Text>
                    )}
                    <Text size="sm" c="dimmed">
                      {t(($) => $.account.connected_apps.connected_at, {
                        date: new Date(app.connectedAt).toLocaleDateString(i18n.language),
                      })}
                    </Text>
                    <Text size="sm" fw={500}>
                      {t(($) => $.account.connected_apps.permissions)}
                    </Text>
                    <List size="sm">
                      {app.scopes.map((scope) => (
                        <List.Item key={scope}>{scopeLabel(scope)}</List.Item>
                      ))}
                    </List>
                  </Stack>
                  <Button
                    variant="default"
                    color="red"
                    loading={revoke.isPending && revoke.variables?.path.id === app.id}
                    onClick={() => {
                      confirm(app);
                    }}
                  >
                    {t(($) => $.account.connected_apps.revoke)}
                  </Button>
                </Group>
              </Card>
            ))
          )}
        </Stack>
      </Stack>
    </Container>
  );
}
