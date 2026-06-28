import {
  type CredentialCreationOptionsJSON,
  type CredentialRequestOptionsJSON,
  create,
  get,
} from "@github/webauthn-json";
import { router } from "@inertiajs/react";
import * as Routes from "@/routes.js";

async function fetchOptions(url: string): Promise<unknown> {
  const response = await fetch(url, {
    headers: { Accept: "application/json" },
    credentials: "same-origin",
  });

  return response.json();
}

export const passkeySupported = (): boolean =>
  typeof window !== "undefined" && Boolean(window.PublicKeyCredential);

export async function loginWithPasskey(): Promise<void> {
  const publicKey = (await fetchOptions(
    Routes.new_passkey_session_path(),
  )) as CredentialRequestOptionsJSON["publicKey"];
  const credential = await get({ publicKey });

  router.post(Routes.passkey_session_path(), {
    credential: JSON.stringify(credential),
  });
}

export async function registerPasskey(): Promise<void> {
  const publicKey = (await fetchOptions(
    Routes.new_account_passkey_path(),
  )) as CredentialCreationOptionsJSON["publicKey"];
  const credential = await create({ publicKey });

  router.post(Routes.account_passkeys_path(), {
    credential: JSON.stringify(credential),
  });
}
