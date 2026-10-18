import { Badge, Button, Card, Group, Title } from "@mantine/core";
import { useTranslation } from "react-i18next";
import {
  adminGetCourseQueryKey,
  adminListCoursesQueryKey,
  adminMarkCourseCompletedMutation,
  adminMarkCourseDraftMutation,
  adminMarkCourseInDevelopmentMutation,
} from "@/client/@tanstack/react-query.gen";
import type { Course } from "@/client/types.gen";
import { useResourceMutation } from "@/hooks/useResourceMutation";
import { transitionsFromReadiness } from "@/lib/courseReadiness";
import type { ReadinessTransition } from "@/lib/courseReadiness";

// A Course's Readiness and the buttons that change it. A state is not an update:
// The course form does not carry it, each button calls its own operation (its own
// Permission on the server). Every value but the current one is offered; the
// Server stays the judge and refuses the rest with a 409.
export function CourseReadinessPanel({ course }: { course: Course }) {
  const { t } = useTranslation();
  const options = {
    invalidate: [adminListCoursesQueryKey(), adminGetCourseQueryKey({ path: { id: course.id } })],
    successMessage: t(($) => $.admin.readiness.changed),
    errorMessage: t(($) => $.admin.readiness.changeError),
  };
  const markDraft = useResourceMutation({ mutation: adminMarkCourseDraftMutation(), ...options });
  const markInDevelopment = useResourceMutation({
    mutation: adminMarkCourseInDevelopmentMutation(),
    ...options,
  });
  const markCompleted = useResourceMutation({
    mutation: adminMarkCourseCompletedMutation(),
    ...options,
  });

  const actions: Record<
    ReadinessTransition,
    { label: string; color: string; isPending: boolean; run: () => void }
  > = {
    markDraft: {
      label: t(($) => $.admin.readiness.markDraft),
      color: "gray",
      isPending: markDraft.isPending,
      run: () => markDraft.mutate({ path: { id: course.id } }),
    },
    markInDevelopment: {
      label: t(($) => $.admin.readiness.markInDevelopment),
      color: "yellow",
      isPending: markInDevelopment.isPending,
      run: () => markInDevelopment.mutate({ path: { id: course.id } }),
    },
    markCompleted: {
      label: t(($) => $.admin.readiness.markCompleted),
      color: "green",
      isPending: markCompleted.isPending,
      run: () => markCompleted.mutate({ path: { id: course.id } }),
    },
  };

  return (
    <Card withBorder p="xl" maw={720}>
      <Group justify="space-between">
        <Group>
          <Title order={4}>{t(($) => $.models.attributes.course.readiness)}</Title>
          <Badge variant="light">
            {t(($) => $.models.attributes.course["readiness/values"][course.readiness ?? "draft"])}
          </Badge>
        </Group>
        <Group>
          {transitionsFromReadiness(course.readiness).map((name) => (
            <Button
              key={name}
              variant="light"
              color={actions[name].color}
              loading={actions[name].isPending}
              onClick={actions[name].run}
            >
              {actions[name].label}
            </Button>
          ))}
        </Group>
      </Group>
    </Card>
  );
}
