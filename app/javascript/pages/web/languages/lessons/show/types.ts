import type { Language, LanguageCategory, LanguageLesson, LanguageLessonMember } from "@/types/serializers";
import type { SharedProps } from "@/types";

export type Props = SharedProps & {
  courseCategory: LanguageCategory;
  course: Language;
  prevLesson?: LanguageLesson;
  nextLesson?: LanguageLesson;
  lessonMember?: LanguageLessonMember;
  lesson: LanguageLesson;
  lessons: LanguageLesson[];
};

export interface RootState {
  processState: "checked" | "unchecked" | "checking";
  currentTab: "editor" | "output" | "tests" | "solution";
  finished: boolean;
  result: 'error' | 'passed' | 'failed' | 'failed-infinity';
  output: string;
  passed: boolean;
  content: string;
  focusesCount: number;
  startTime: number;
  solutionState: "shown" | "canBeShown" | "notAllowedToBeShown";
}

export interface CheckingResponse {
  result: boolean;
  output: string;
  passed: boolean;
}

