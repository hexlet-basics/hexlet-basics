import { router } from "@inertiajs/react";
import type {
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
} from "@simplewebauthn/browser";
import {
  browserSupportsWebAuthn,
  browserSupportsWebAuthnAutofill,
  startAuthentication,
  startRegistration,
  WebAuthnAbortService,
  WebAuthnError,
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

export const passkeyAutofillSupported = (): Promise<boolean> =>
  browserSupportsWebAuthnAutofill();

export const cancelPasskeyCeremony = (): void => {
  WebAuthnAbortService.cancelCeremony();
};

export function passkeyCancelled(error: unknown): boolean {
  if (!(error instanceof WebAuthnError)) {
    return false;
  }

  return (
    error.code === "ERROR_CEREMONY_ABORTED" ||
    (error.code === "ERROR_PASSTHROUGH_SEE_CAUSE_PROPERTY" &&
      error.cause instanceof DOMException &&
      error.cause.name === "NotAllowedError")
  );
}

export function passkeyPreviouslyRegistered(error: unknown): boolean {
  return (
    error instanceof WebAuthnError &&
    error.code === "ERROR_AUTHENTICATOR_PREVIOUSLY_REGISTERED"
  );
}

export async function loginWithPasskey(
  useBrowserAutofill = false,
): Promise<void> {
  const optionsJSON = (await fetchOptions(
    Routes.new_passkey_session_path(),
  )) as PublicKeyCredentialRequestOptionsJSON;
  const credential = await startAuthentication({
    optionsJSON,
    useBrowserAutofill,
  });

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
