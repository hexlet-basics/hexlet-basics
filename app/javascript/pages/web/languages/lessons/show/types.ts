import type { Language, LanguageCategory, LanguageLesson, LanguageLessonMember } from "@/types/serializers";
import type { SharedProps } from "@/types/types";

export type Props = SharedProps & {
  courseCategory: LanguageCategory;
  course: Language;
  prevLesson?: LanguageLesson;
  nextLesson?: LanguageLesson;
  lessonMember?: LanguageLessonMember;
  lesson: LanguageLesson;
  lessons: LanguageLesson[];
};
