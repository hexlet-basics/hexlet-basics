// This file is auto-generated. Do not edit manually.
// Run `bin/rails app:export_events_to_ts` to update.

export interface BookRequestedEvent {
  type: "BookRequestedEvent";
  name: "book_requested";
  event_id: string;
  metadata: { name: string };
  data: {
    locale: string;
  };
}

export interface CourseFinishedEvent {
  type: "CourseFinishedEvent";
  name: "course_finished";
  event_id: string;
  metadata: { name: string };
  data: {
    occurrence_count: number;
    slug: string;
    locale: string;
  };
}

export interface CourseStartedEvent {
  type: "CourseStartedEvent";
  name: "course_started";
  event_id: string;
  metadata: { name: string };
  data: {
    occurrence_count: number;
    slug: string;
    locale: string;
  };
}

export interface EmailConfirmedEvent {
  type: "EmailConfirmedEvent";
  name: "email_confirmed";
  event_id: string;
  metadata: { name: string };
  data: Record<string, never>;
}

export interface LeadCreatedEvent {
  type: "LeadCreatedEvent";
  name: "lead_created";
  event_id: string;
  metadata: { name: string };
  data: {
    lead_id: number;
    user_id: number;
    user_name: string;
    first_name: string | undefined;
    last_name: string | undefined;
    ym_client_id: string | undefined;
    utm_source: string | undefined;
    utm_medium: string | undefined;
    utm_campaign: string | undefined;
    utm_term: string | undefined;
    utm_content: string | undefined;
    email: string | undefined;
    phone: string | undefined;
    telegram: string | undefined;
    whatsapp: string | undefined;
  };
}

export interface LessonFinishedEvent {
  type: "LessonFinishedEvent";
  name: "lesson_finished";
  event_id: string;
  metadata: { name: string };
  data: {
    occurrence_count: number;
    lesson_slug: string;
    course_slug: string;
    locale: string;
  };
}

export interface LessonStartedEvent {
  type: "LessonStartedEvent";
  name: "lesson_started";
  event_id: string;
  metadata: { name: string };
  data: {
    occurrence_count: number;
    lesson_slug: string;
    course_slug: string;
    locale: string;
  };
}

export interface SolutionCheckedEvent {
  type: "SolutionCheckedEvent";
  name: "solution_checked";
  event_id: string;
  metadata: { name: string };
  data: {
    lesson_slug: string;
    course_slug: string;
    locale: string;
    passed: boolean;
  };
}

export interface UserSignedInEvent {
  type: "UserSignedInEvent";
  name: "user_signed_in";
  event_id: string;
  metadata: { name: string };
  data: {
    user_id: number;
    occurrence_count: number;
    email: string | undefined;
    locale: string;
  };
}

export interface UserSignedUpEvent {
  type: "UserSignedUpEvent";
  name: "user_signed_up";
  event_id: string;
  metadata: { name: string };
  data: {
    user_id: number;
    email: string | undefined;
    first_name: string | undefined;
    last_name: string | undefined;
    locale: string;
  };
}

export type ApplicationEvent =
  | BookRequestedEvent
  | CourseFinishedEvent
  | CourseStartedEvent
  | EmailConfirmedEvent
  | LeadCreatedEvent
  | LessonFinishedEvent
  | LessonStartedEvent
  | SolutionCheckedEvent
  | UserSignedInEvent
  | UserSignedUpEvent;
