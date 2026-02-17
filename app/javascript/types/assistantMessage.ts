export type AssistantMessage = {
  id?: string;
  role: "user" | "assistant" | "system";
  content: string;
};
