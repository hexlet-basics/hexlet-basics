// Typelizer digest e0e1e0ba9a97742f0d52c53ad87d2369
//
// DO NOT MODIFY: This file was automatically generated by Typelizer.

type LanguageLessonMemberMessage = {
  id: number;
  language_id: number;
  language_lesson_id: number;
  language_lesson_member_id: number;
  role: "user" | "assistant" | null;
  body: string | null;
  created_at: string;
  language_slug: string;
  language_lesson_slug: string;
  user_id: number;
  language_lesson_name: string;
}

export default LanguageLessonMemberMessage;
