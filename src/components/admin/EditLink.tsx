import { Button, type ButtonProps } from "@mantine/core";
import { createLink, type LinkComponent } from "@tanstack/react-router";
import { forwardRef } from "react";
import { useTranslation } from "react-i18next";

// The "Edit" button of a list row. It is a typed router link, so each resource
// passes its own `to` and `params` (route paths are literals the router checks;
// nothing here hides them) while the look and the label are shared.
const EditLinkBase = forwardRef<HTMLAnchorElement, Omit<ButtonProps, "href" | "children">>(
  (props, ref) => {
    const { t } = useTranslation();
    return (
      <Button ref={ref} component="a" size="xs" variant="light" {...props}>
        {t(($) => $.admin.crud.edit)}
      </Button>
    );
  },
);

const EditLinkCreated = createLink(EditLinkBase);

export const EditLink: LinkComponent<typeof EditLinkBase> = (props) => (
  <EditLinkCreated preload="intent" {...props} />
);
