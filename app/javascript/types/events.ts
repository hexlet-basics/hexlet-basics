type UserSignedInEvent = {
  type: "UserSignedInEvent";
  data: {
    id: number;
    email: string;
    first_name: string;
    last_name: string;
  };
};

type UserSignedUpEvent = {
  type: "UserSignedUpEvent";
  data: {
    id: number;
    email: string;
    first_name: string;
    last_name: string;
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

type LeadCreatedEvent = {
  type: "LeadCreatedEvent";
  data: {
    user_id: number;
    email: string;
    phone: string;
    first_name: string;
    last_name: string;
    // email: String,
    // phone: T.nilable(String),
    // telegram: T.nilable(String),
    // whatsapp: T.nilable(String),
    // survey_answers_data: T::Array[T::Hash[String, T.untyped]],
    // courses_data: T::Array[T::Hash[String, T.untyped]]
  };
};

export type BackendEvent =
  | UserSignedInEvent
  | UserSignedUpEvent
  | CourseStartedEvent
  | CourseFinishedEvent
  | LessonStartedEvent
  | LeadCreatedEvent
  | LessonFinishedEvent;
