import * as Routes from "@/routes.js";
import type { User } from "@/types";
import { Link } from "@inertiajs/react";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { Submit } from "use-inertia-form";
import XssContent from "./XssContent";
import { XForm, XInput } from "./forms";

type Props = PropsWithChildren & {
  user: User;
  autoFocus?: boolean
};

export default function SignUpFormBlock({ user, autoFocus = false }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <XForm model="user_sign_up_form" data={user} to={Routes.users_path()}>
      <XInput name="first_name" autoComplete="name" autoFocus={autoFocus} />
      <XInput name="email" autoComplete="email" />
      <XInput name="password" type="password" autoComplete="current-password" />
      <div className="text-end text-muted small mb-4">
        {t("users.new.have_account")}{" "}
        <Link href={Routes.new_session_path()} className="text-decoration-none">
          {t("users.new.sign_in")}
        </Link>
      </div>
      <Submit className="btn w-100 btn-lg btn-primary mb-3">
        {tHelpers("submit.user_sign_up_form.create")}
      </Submit>
      <XssContent className="small text-muted">
        {t("users.new.confirmation_html", {
          url: Routes.page_path("tos"),
        })}
      </XssContent>
    </XForm>
  );
}
