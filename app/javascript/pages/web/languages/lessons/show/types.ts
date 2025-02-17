import type { SharedProps } from "@/types";
import type {
  Language,
  LanguageCategory,
  LanguageLesson,
  LanguageLessonMember,
  LessonCheckingResponse,
} from "@/types/serializers";
import { SetupOptions } from "node_modules/@inertiajs/react/types/createInertiaApp";

export type LessonSharedProps = SharedProps & {
  courseCategory: LanguageCategory;
  course: Language;
  prevLesson?: LanguageLesson;
  nextLesson?: LanguageLesson;
  lessonMember?: LanguageLessonMember;
  lesson: LanguageLesson;
  lessons: LanguageLesson[];
};

type CheckingResult = LessonCheckingResponse["result"] | "error" | null;

export interface RootState {
  processState: "checked" | "unchecked" | "checking";
  currentTab: "editor" | "output" | "tests" | "solution";
  finished: boolean;
  result: CheckingResult;
  resetsCount: number;
  defaultCode: string;
  output: string;
  passed: boolean;
  content: string;
  focusesCount: number;
  startTime: number;
  solutionState: "shown" | "canBeShown" | "notAllowedToBeShown";
}
