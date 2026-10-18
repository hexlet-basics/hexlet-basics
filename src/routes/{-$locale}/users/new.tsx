import { Card, Container, Stack, Title } from "@mantine/core";
import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import SignUpFormBlock from "@/components/SignUpFormBlock";
import { safeRedirectPath } from "@/lib/safe-redirect";
import { lenient } from "@/lib/lenient";

// Registration page, ported from legacy users/new. The form itself is the
// shared SignUpFormBlock.
// Where to go once the account exists. The lesson player sends a guest here with
// the lesson they were on, so signing up does not lose their place. Only a path
// on this site is honoured — anything else is dropped rather than followed, so
// the page cannot be used to bounce a new account off to another origin.
const signUpSearchSchema = z.object({
  redirect: lenient(z.string().transform(safeRedirectPath)),
});

export const Route = createFileRoute("/{-$locale}/users/new")({
  component: New,
  validateSearch: signUpSearchSchema,
});

function New() {
  const { t } = useTranslation();
  const { redirect } = Route.useSearch();

  return (
    <Container my="xl">
      <Stack align="center">
        <Title order={1} ta="center">
          {t(($) => $.users.new.title)}
        </Title>

        <Card withBorder p="xl" w={{ base: "100%", lg: "50%", md: "70%", sm: "80%" }}>
          <SignUpFormBlock redirect={redirect} />
        </Card>
      </Stack>
    </Container>
  );
}
