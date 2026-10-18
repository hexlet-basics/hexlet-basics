import { Badge, Button, Card, Group, Title } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { transitionsFrom } from "@/lib/publicationState";
import type { PublicationState, StateTransition } from "@/lib/publicationState";
import type { ReactElement } from "react";

// What the panel needs of one transition's mutation: the generated mutation of
// that operation, wired by the resource's own panel.
export interface StateMutation {
  isPending: boolean;
  mutate: (variables: { path: { id: number } }) => void;
}

// The Publication State of a record and the buttons that change it, shared by
// every resource that has one. A state is not an update: the form does not carry
// it, each button calls its own operation (its own Permission on the server).
// Only the moves legal for the current state are offered; the server stays the
// judge and refuses the rest with a 409.
export function PublicationStatePanel({
  id,
  title,
  state,
  stateLabel,
  mutations,
}: {
  id: number;
  title: string;
  state: PublicationState;
  stateLabel: string;
  mutations: Record<StateTransition, StateMutation>;
}): ReactElement {
  const { t } = useTranslation();
  const actions: Record<StateTransition, { label: string; color: string }> = {
    archive: { color: "red", label: t(($) => $.admin.publicationState.archive) },
    publish: { color: "green", label: t(($) => $.admin.publicationState.publish) },
    unpublish: { color: "gray", label: t(($) => $.admin.publicationState.unpublish) },
  };

  return (
    <Card withBorder p="xl" maw={720}>
      <Group justify="space-between">
        <Group>
          <Title order={4}>{title}</Title>
          <Badge variant="light">{stateLabel}</Badge>
        </Group>
        <Group>
          {transitionsFrom(state).map((name) => (
            <Button
              key={name}
              variant="light"
              color={actions[name].color}
              loading={mutations[name].isPending}
              onClick={() => {
                mutations[name].mutate({ path: { id } });
              }}
            >
              {actions[name].label}
            </Button>
          ))}
        </Group>
      </Group>
    </Card>
  );
}
