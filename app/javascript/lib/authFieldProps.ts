export const authInputModeProps = {
  autoCapitalize: "none" as const,
  inputMode: "email" as const,
  spellCheck: false,
  type: "email" as const,
};

export const firstNameInputProps = {
  autoComplete: "given-name",
  enterKeyHint: "next" as const,
};

export const loginEmailInputProps = {
  ...authInputModeProps,
  autoComplete: "username",
  enterKeyHint: "next" as const,
};

export const emailInputProps = {
  ...authInputModeProps,
  autoComplete: "email",
};

export const registrationEmailInputProps = {
  ...emailInputProps,
  enterKeyHint: "next" as const,
};

export const currentPasswordInputProps = {
  autoCapitalize: "none" as const,
  autoComplete: "current-password",
  enterKeyHint: "go" as const,
  spellCheck: false,
  type: "password" as const,
};

export const newPasswordInputProps = {
  autoCapitalize: "none" as const,
  autoComplete: "new-password",
  enterKeyHint: "done" as const,
  spellCheck: false,
  type: "password" as const,
};
