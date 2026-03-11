import { getResourceUrl } from "@/resources";

export type ReviewShowcaseId = "avdoshkin" | "tyurin" | "kleyman";

export const reviewShowcaseOrder: ReviewShowcaseId[] = [
  "avdoshkin",
  "tyurin",
  "kleyman",
];

export const reviewShowcaseAvatars: Record<ReviewShowcaseId, string> = {
  avdoshkin: getResourceUrl("avdoshkin.jpg"),
  tyurin: getResourceUrl("tyrin.jpg"),
  kleyman: getResourceUrl("user-avatar.png"),
};
