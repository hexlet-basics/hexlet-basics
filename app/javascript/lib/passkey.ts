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
  typeof window !== "undefined" &&
  typeof window.PublicKeyCredential?.parseRequestOptionsFromJSON === "function";

export async function loginWithPasskey(): Promise<void> {
  const options = (await fetchOptions(
    Routes.new_passkey_session_path(),
  )) as PublicKeyCredentialRequestOptionsJSON;
  const publicKey = PublicKeyCredential.parseRequestOptionsFromJSON(options);
  const credential = await navigator.credentials.get({ publicKey });

  if (!(credential instanceof PublicKeyCredential)) {
    throw new Error("Passkey authentication was cancelled");
  }

  router.post(Routes.passkey_session_path(), {
    credential: JSON.stringify(credential.toJSON()),
  });
}

export async function registerPasskey(): Promise<void> {
  const options = (await fetchOptions(
    Routes.new_account_passkey_path(),
  )) as PublicKeyCredentialCreationOptionsJSON;
  const publicKey = PublicKeyCredential.parseCreationOptionsFromJSON(options);
  const credential = await navigator.credentials.create({ publicKey });

  if (!(credential instanceof PublicKeyCredential)) {
    throw new Error("Passkey registration was cancelled");
  }

  router.post(Routes.account_passkeys_path(), {
    credential: JSON.stringify(credential.toJSON()),
  });
}
