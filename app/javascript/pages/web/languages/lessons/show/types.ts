import type {
  Language,
  LanguageLesson,
  LanguageLessonMember,
  LessonCheckingResponse,
} from "@/types";

type CheckingResult = LessonCheckingResponse["result"] | "error" | null;

type TabName = "editor" | "output" | "tests" | "solution";
type SolutionState = "shown" | "canBeShown" | "notAllowedToBeShown";
type ProcessState = "checked" | "unchecked" | "checking";

export interface LessonProps {
  lessonMember?: LanguageLessonMember;
  lesson: LanguageLesson;
}

export interface LessonState {
  startTime: number;
  solutionState: SolutionState;
  finished: boolean;
  defaultCode: string;
  processState: ProcessState;
  currentTab: TabName;
  result: CheckingResult;
  resetsCount: number;
  output: string;
  passed: boolean;
  content: string;
  focusesCount: number;
  changeContent: (content: string) => void;
  resetContent: () => void;
  changeTab: (tab: TabName) => void;
  setStartTime: (startTime: number) => void;
  changeSolutionState: (solutionState: SolutionState) => void;
  runCheck: (params: {
    course: Language;
    lesson: LanguageLesson;
  }) => Promise<boolean>;
}
