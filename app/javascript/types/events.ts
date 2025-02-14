type UserSignedInEvent = {
  type: "UserSignedInEvent";
  data: {
    id: string;
    email: string;
  };
};

type UserSignedUpEvent = {
  type: "UserSignedUpEvent";
  data: {
    id: string;
    email: string;
  };
};

type CourseStartedEvent = {
  type: "CourseStartedEvent";
  data: {
    slug: string;
    locale: string;
  };
};

type CourseFinishedEvent = {
  type: "CourseFinishedEvent";
};

type LessonStartedEvent = {
  type: "LessonStartedEvent";
  data: {
    course_slug: string;
    lesson_slug: string;
    locale: string;
  };
};

type LessonFinishedEvent = {
  type: "LessonFinishedEvent";
};

export type BackendEvent =
  | UserSignedInEvent
  | UserSignedUpEvent
  | CourseStartedEvent
  | CourseFinishedEvent
  | LessonStartedEvent
  | LessonFinishedEvent;
