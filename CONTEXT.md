# Domain glossary

Shared vocabulary for `hexlet-basics`. Keep terms here in sync with the code as
the domain sharpens.

## Course (model: `Language`)

A learning course. The Active Record model is `Language` (historical), but the
**ubiquitous language is "Course"**: domain events (`CourseStartedEvent`,
`CourseFinishedEvent`), serialized props (`course`, `courseMember`,
`courseModules`), and analytics all say *course*. New domain code should prefer
the *Course* vocabulary even while operating on `Language` (e.g. the service is
`CourseProgressService`). A future rename `Language → Course` would close this
split.

## Course membership (`Language::Member`)

A user's enrollment in a course and the aggregate that holds their progress:
AASM `state` (`started` → `finished`), `finished_lessons_count`, the
`all_lessons_finished?` guard, and `next_lesson_info`. The `finish` transition is
guarded so a course finishes only once every lesson is done.

## Lesson membership (`Language::Lesson::Member`)

A user's progress on a single lesson within a course membership. AASM `state`
(`started` → `finished`).

## Course progress lifecycle (`CourseProgressService`)

The deep module that owns the lifecycle scattered across controllers before:
*start course → start lesson → solution checked → finish lesson → finish course*.
It drives the membership AASM transitions, publishes the resulting domain events
to RES (`EventSender.publish_event`), and returns those events so controllers can
selectively forward them to the frontend via `js_event` (used only when the
frontend needs an event, e.g. for analytics).
