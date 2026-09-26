import { Anchor, Box, Button, Stack } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation } from "@tanstack/react-query";
import { useNavigate, useParams } from "@tanstack/react-router";
import { Trans, useTranslation } from "react-i18next";
import { createLeadMutation } from "@/client/@tanstack/react-query.gen";
import type { LeadInput } from "@/client/types.gen";
import { zLeadInput } from "@/client/zod.gen";
import { readFirstVisit } from "@/lib/first-visit";
import { useAppForm } from "@/lib/form";
import { metrikaClientId } from "@/lib/metrika";

type Props = {
  autoFocus?: boolean;
  // Where to go once the lead is sent. Legacy redirected every form to its
  // `from` parameter or, without one, to the home page.
  redirectHref?: string;
};

// The consultation request form, ported from legacy LeadFormBlock and shared by
// every page that embeds it (leads/new, the book page, course completion, the
// catalog pages). It leaves gating (signed in, ru only) to the embedding page,
// as legacy did. The Metrika client id and the first-visit attribution are
// read at submit time, not render time: both live only in the browser.
export default function LeadFormBlock({ autoFocus = false, redirectHref }: Props) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { locale } = useParams({ strict: false });
  const mutation = useMutation(createLeadMutation());

  const contactMethods = t(($) => $.models.attributes.user["contact_method/values"], {
    returnObjects: true,
  });
  const contactMethodData = Object.entries(contactMethods).map(([value, label]) => ({
    value,
    label,
  }));

  const communityLink = (
    <Anchor
      href={t(($) => $.common.community_url).trim()}
      target="_blank"
      rel="noopener noreferrer"
    />
  );

  const defaultValues: LeadInput = {
    contactMethod: "telegram",
    contactValue: "",
    ymClientId: null,
    firstVisit: null,
  };

  const form = useAppForm({
    defaultValues,
    validators: { onSubmit: zLeadInput },
    onSubmit: async ({ value }) => {
      try {
        await mutation.mutateAsync({
          body: { ...value, ymClientId: metrikaClientId(), firstVisit: readFirstVisit() },
        });
      } catch {
        notifications.show({ color: "red", message: t(($) => $.flash.leads.create.error) });
        return;
      }
      notifications.show({
        color: "green",
        message: (
          <Trans
            t={t}
            i18nKey={($) => $.flash.leads.create.success}
            components={{ a: communityLink }}
          />
        ),
      });
      if (redirectHref) {
        await navigate({ href: redirectHref });
      } else {
        await navigate({ to: "/{-$locale}", params: { locale } });
      }
    },
  });

  return (
    <Stack
      component="form"
      gap={0}
      onSubmit={(event) => {
        event.preventDefault();
        form.handleSubmit();
      }}
    >
      <form.AppField name="contactMethod">
        {(field) => (
          <field.SelectField
            aria-label={t(($) => $.models.attributes.lead.contact_method)}
            data={contactMethodData}
            allowDeselect={false}
            required
            mb="sm"
          />
        )}
      </form.AppField>
      <form.AppField name="contactValue">
        {(field) => (
          <field.TextField
            aria-label={t(($) => $.models.attributes.lead.contact_value)}
            required
            // Opt-in only: the course completion page puts the visitor straight
            // into the form, as legacy did.
            // oxlint-disable-next-line jsx-a11y/no-autofocus
            autoFocus={autoFocus}
            mb="sm"
          />
        )}
      </form.AppField>
      <Box fz="sm" my="lg">
        {t(($) => $.blocks.lead_form_block.description1)}{" "}
        <Trans
          t={t}
          i18nKey={($) => $.blocks.lead_form_block.description2}
          components={{ a: communityLink }}
        />
      </Box>
      <Button type="submit" fullWidth loading={mutation.isPending}>
        {t(($) => $.helpers.send)}
      </Button>
    </Stack>
  );
}
