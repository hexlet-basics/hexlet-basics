# Domain Glossary

## Naming rule

A term defined here is the name used **everywhere the name is ours to choose**:
Go and TypeScript identifiers, contract models and fields, ent schema and edge
names, i18n keys, documentation, tests, and open issues.

Three things keep their legacy names, because renaming them would move data or
break a link someone else holds:

- **table names** (`language_members`, `language_lesson_members`, …),
- **column names** (`members_count`, `progress`, `language_id`, …),
- **public URLs** (`/languages/:slug`, `/languages/:slug/lessons/:slug`).

An ent field or edge name is *not* a column name — `field.Int("course_id").StorageKey("language_id")`
renames the Go identifier and leaves storage untouched, so ent names follow the
glossary like everything else. The legacy identifier behind each term is listed
in [Storage mapping](#storage-mapping).

One exception is temporary rather than permanent: the **values** of
`PermissionResource` (`"language_lesson_members"`, `"languages"`, …) are rows in
`staff_member_role_permissions`, and the legacy Rails app reads the same rows in
production. They are renamed by a data migration after the hard cutover
(ADR-0002), not before.

Terms that reach a learner or an administrator carry their **ru** and **es**
wording, so the same concept is not translated two ways on two screens. Terms
that only ever appear in code carry none.

## Language

### Accounts

**User sign-up**:
The first successful creation of a User account through public registration.
Later sign-ins, profile changes, and administrator-created users are not
sign-ups.

### Learning

**Enrollment**:
A learner's record of where they are in one Course, and never created twice for
the same pair. It is created by a deliberate action — starting a Lesson, or
submitting a solution to one — and never by loading a page. It grants nothing:
every Course is readable without one, and without an account. An Enrollment
remembers progress, it does not permit access.
_ru_: запись на курс · _es_: inscripción al curso

**Lesson Progress**:
A learner's record of how far they have got with one Lesson of a Course they are
enrolled in. One row per learner per Lesson, carrying the state and nothing the
Lesson itself already says.
_ru_: прохождение урока · _es_: progreso de la lección

**Staff Member**:
An employee, carrying roles and permissions. Nothing to do with learning; it
only shares a table-name prefix with Enrollment and Lesson Progress.
_ru_: сотрудник · _es_: empleado

**Not started**:
The condition of a learner who has no Enrollment for a Course. It is the absence
of a record, not a state the record can be in.
_ru_: не начат · _es_: no iniciado

**Started**:
The state of an Enrollment or Lesson Progress from creation until it is
finished. (`EnrollmentState`, carried by both.)
_ru_: начат · _es_: iniciado

**Finished**:
The state of a Lesson Progress whose exercise the learner has passed, and of an
Enrollment all of whose Current Lessons are finished. Neither ever returns to
Started. Merging Guest progress is the one case where a Lesson Progress is
created Finished without the learner having passed that Lesson individually:
sequential progression means the whole prefix up to their furthest Lesson was
passed, so the rule implies the rest.
_ru_: завершён · _es_: finalizado

**Current Lessons**:
The Lessons belonging to a Course's current Version. A Course's Lesson set
changes when a new Version is promoted, so what counts as finishing the Course
is defined against the current Version and not against Lessons a learner may
have finished under an older one.

**Position**:
Where a Lesson sits in course order — a property of the Lesson's Version, not of
the Lesson, and therefore only meaningful relative to one Course Version. A
Lesson a learner finished that is absent from the current Version has no
position: it neither raises the gate nor counts towards Completion.

**Available Lesson**:
A Lesson a learner is allowed to take: one whose Position is at most one past
the highest Position they have finished. Availability is measured from the
highest finished Lesson, not from the first unfinished one, so a gap in a
learner's history does not block them. A learner with no history has exactly one
Available Lesson, the first.
_ru_: доступный урок · _es_: lección disponible

**Completion**:
The share of a Course's Current Lessons a learner has finished, as a percentage.
Lessons finished under a retired Version do not count towards it. Carried as
`Enrollment.completion`; the word Progress on its own is ambiguous now that
Lesson Progress is a record, so Completion is the name for the number.
_ru_: прогресс · _es_: progreso

**Course Readiness**:
How finished the Course itself is — draft, in development, or completed
(`CourseReadiness`). A property of the product, not of any learner.
_ru_: готовность курса · _es_: disponibilidad del curso

**Next Lesson**:
The first Lesson, in course order, among a Course's Current Lessons that the
learner has not finished — where they resume. It is not the gate: for a learner
with a gap, the Next Lesson is the gap, while their Available Lessons run all
the way to one past their furthest finished one.
_ru_: следующий урок · _es_: siguiente lección

**Guest progress**:
Lessons a visitor passed before having an account. It is one Lesson per Course —
the furthest one they finished — carried in a signed, httpOnly cookie, because
sequential progression makes a guest's gaps unrepresentable and that single
value therefore says everything. It is merged into an account on sign-up and on
sign-in, taking the further of the two positions, and becomes an ordinary
Enrollment plus Lesson Progress rows in the Finished state for the whole prefix
up to it.
_ru_: гостевой прогресс · _es_: progreso de invitado

## Storage mapping

Where a term's storage still carries the legacy word. This is the working half
of the naming rule: it says which legacy identifier a concept is allowed to keep,
and it is the only place that mapping is written down for concepts whose ent
schema does not exist yet.

| Term | Legacy table | Legacy columns worth knowing |
| --- | --- | --- |
| Enrollment | `language_members` | `finished_lessons_count`, `state`; Completion is computed, not stored |
| Lesson Progress | `language_lesson_members` | `language_member_id`, `lesson_id`, `state`, `UNIQUE (user_id, lesson_id)` |
| Course | `languages` | `members_count` (Enrollment count), `progress` (Course Readiness) |
| Course Lesson | `language_lessons` | `language_id`, `module_id` |
| Course Version | `language_versions` | `language_id` |
| Course Module | `language_modules` | `language_id` |
| Course Category | `language_categories` | — |
| Course Category Item | `language_category_items` | `language_category_id`, `language_landing_page_id`; no ent schema yet |
| Course Landing Page | `language_landing_pages` | `language_id` |
| Course Lesson Review | `language_lesson_reviews` | `language_lesson_version_id` |
| Staff Member | `staff_members` | — |
