import { router } from "@inertiajs/react";
import type {
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
} from "@simplewebauthn/browser";
import {
  browserSupportsWebAuthn,
  startAuthentication,
  startRegistration,
} from "@simplewebauthn/browser";
import * as Routes from "@/routes.js";

async function fetchOptions(url: string): Promise<unknown> {
  const response = await fetch(url, {
    headers: { Accept: "application/json" },
    credentials: "same-origin",
  });

  return response.json();
}

export const passkeySupported = (): boolean => browserSupportsWebAuthn();

export async function loginWithPasskey(): Promise<void> {
  const optionsJSON = (await fetchOptions(
    Routes.new_passkey_session_path(),
  )) as PublicKeyCredentialRequestOptionsJSON;
  const credential = await startAuthentication({ optionsJSON });

  router.post(Routes.passkey_session_path(), {
    credential: JSON.stringify(credential),
  });
}

export async function registerPasskey(): Promise<void> {
  const optionsJSON = (await fetchOptions(
    Routes.new_account_passkey_path(),
  )) as PublicKeyCredentialCreationOptionsJSON;
  const credential = await startRegistration({ optionsJSON });

  router.post(Routes.account_passkeys_path(), {
    credential: JSON.stringify(credential),
  });
}
