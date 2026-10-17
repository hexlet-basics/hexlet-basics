# Hexlet Basics

Free interactive programming courses: a learner works through a Course's Lessons
in the browser, finishing each one by passing its Lesson Tests.

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

A fourth exception existed for a while and is gone: the **values** of
`PermissionResource` (`"language_lesson_members"`, `"languages"`, …) were rows in
the create/read/update/delete matrix. ADR-0017 replaced that matrix with
Permissions named with the terms (`courses.read`, `lesson_progress.read`), and
migrations carried the existing rows into them, so the legacy names are no
longer identifiers anywhere.

Terms that reach a learner or an administrator carry their **ru** and **es**
wording, so the same concept is not translated two ways on two screens. Terms
that only ever appear in code carry none.

## Language

### Accounts

**User sign-up**:
The first successful creation of a User account through public registration.
Later sign-ins, profile changes, and administrator-created users are not
sign-ups.

**Session**:
A signed-in User's hold on the site between requests: an opaque token in an
httpOnly cookie that names a server-side record. It lives 31 days and ends after
7 without a request, when the User signs out, when the account is removed, and
when the User's password changes. It proves who the User is and grants nothing
else: what they may do is read from the account on every request. Not the guest
progress cookie, which belongs to no account.
_ru_: сессия · _es_: sesión

**Magic Link**:
An emailed link that signs an existing User in without a password. It lives
fifteen minutes, is bound to the User's email so it dies if the email changes,
and may be followed more than once while it lives. It never creates an account,
and asking for one reveals nothing about whether the email is registered.
_ru_: ссылка для входа · _es_: enlace de acceso

**Password Reset**:
Replacing a forgotten password through an emailed link. The link lives fifteen
minutes and dies the moment the password changes; following it through signs
the User in. Asking for one reveals nothing about whether the email is
registered.
_ru_: восстановление пароля · _es_: restablecimiento de contraseña

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
An employee who holds one Staff Role and works in the locales their record
lists. Nothing to do with learning; it only shares a table-name prefix with
Enrollment and Lesson Progress.
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
Started. Lessons may be taken in any order, so finished is a property of each
Lesson on its own and a later one implies nothing about an earlier one.
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
position: it does not count towards Completion.

**Completion**:
The share of a Course's Current Lessons a learner has finished, as a percentage.
Lessons finished under a retired Version do not count towards it. The word
Progress on its own is ambiguous now that Lesson Progress is a record, so
Completion is the name for the number. Carried as `CourseProgress.completion`:
computed against the current Version, never stored on the Enrollment.
_ru_: прогресс · _es_: progreso

**Learner Lesson**:
A Lesson as the course page and the player list it: the Lesson with the
visitor's Position and whether they finished it. (`LearnerLesson`.)
_ru_: урок ученика · _es_: lección del alumno

**Course Progress**:
Where a visitor stands in one Course: its state, Completion and Next Lesson.
(`CourseProgress`.)
_ru_: прогресс по курсу · _es_: progreso del curso

**Course Readiness**:
How finished the Course itself is — draft, in development, or completed
(`CourseReadiness`). A property of the product, not of any learner.
_ru_: готовность курса · _es_: disponibilidad del curso

**Next Lesson**:
The first Lesson, in course order, among a Course's Current Lessons that the
learner has not finished — where they resume. For a learner who took Lessons out
of order, it is the first gap. Carried as `CourseProgress.nextLesson`. Not the
Lesson after the one being viewed: that is an Adjacent Lesson.
_ru_: следующий урок · _es_: siguiente lección

**Adjacent Lessons**:
The Lessons immediately before and after a given Lesson in course order, for
back/forward navigation in the player. Independent of what the learner has
finished, and so not the Next Lesson. Carried as `LessonView.adjacent`
(`previous`, `next`); either is absent at the edge of the Course.
_ru_: соседние уроки · _es_: lecciones adyacentes

**Starter Code**:
The code a Lesson's editor opens with, written by the Course author as the
learner's starting point, and what resetting the editor restores. Carried as
`CourseLesson.preparedCode`.
_ru_: заготовка кода · _es_: código inicial

**Reference Solution**:
The Course author's own answer to a Lesson, shown to a learner once they have
passed it or waited out the delay. Not "the teacher's solution": nobody teaches
here, an author wrote it. Carried as `CourseLesson.originalCode`.
_ru_: авторское решение · _es_: solución del autor

**Lesson Tests**:
The checks a submission has to satisfy to finish a Lesson. They are shown to
everyone, including visitors with no account — deliberately, since a learner who
reads the tests has read the specification of the exercise. Carried as
`CourseLesson.testCode`.
_ru_: тесты урока · _es_: pruebas de la lección

**Guest progress**:
Lessons a visitor passed before having an account. It is the set of finished
Lesson slugs per Course, carried in a signed, httpOnly cookie. It is merged into
an account on sign-up and on sign-in as a union, and becomes an ordinary
Enrollment plus Lesson Progress rows in the Finished state for exactly those
Lessons.
_ru_: гостевой прогресс · _es_: progreso de invitado

### Staff access and audit

**Staff Role**:
A named set of Permissions that a Staff Member holds.
_ru_: роль сотрудника · _es_: rol del empleado

**Permission**:
The right to perform one kind of staff operation, such as publishing blog posts,
or to use a staff capability, such as reading the Audit Record. Reading and
changing are separate Permissions.
_ru_: право · _es_: permiso

**Scope**:
The kind of access a User consented to give a Connected App, such as reading the
catalogue. It only narrows what the User's Permissions already allow and never
adds to them.
_Avoid_: Permission (a Permission is what a Staff Role holds, a Scope is what a
User gave away)
_ru_: область доступа · _es_: ámbito de acceso

**Connected App**:
An external application, typically an AI assistant, that a User has authorized
to read data on their behalf and, where their Permissions allow, to change it.
The User can revoke it at any time.
_ru_: подключённое приложение · _es_: aplicación conectada

**Audit Record**:
The permanent, unalterable note of one change a staff member made, or of one
read of personal data: who did it, through which Connected App or the website,
to what, and with what result.
_Avoid_: Audit event, audit log entry ("event" already means a domain event)
_ru_: запись аудита · _es_: registro de auditoría

## Storage mapping

Where a term's storage still carries the legacy word. This is the working half
of the naming rule: it says which legacy identifier a concept is allowed to keep,
and it is the only place that mapping is written down for concepts whose ent
schema does not exist yet.

| Term | Legacy table | Legacy columns worth knowing |
| --- | --- | --- |
| Enrollment | `language_members` | `finished_lessons_count`, `state`; Completion is computed, not stored |
| Lesson Progress | `language_lesson_members` | `language_member_id`, `lesson_id`, `state`, `UNIQUE (user_id, lesson_id)` |
| Course | `languages` | `members_count` (ent `enrollments_count`), `progress` (ent `readiness`) |
| Course Lesson | `language_lessons` | `language_id`, `module_id` |
| Course Version | `language_versions` | `language_id` |
| Course Module | `language_modules` | `language_id` |
| Course Category | `language_categories` | — |
| Course Category Item | `language_category_items` | `language_category_id`, `language_landing_page_id` |
| Course Landing Page | `language_landing_pages` | `language_id` |
| Course Lesson Review | `language_lesson_reviews` | `language_lesson_version_id` |
| Staff Member | `staff_members` | — |
| Staff Role | `staff_member_roles` | — |
| Permission | `staff_member_role_grants` | one row per Permission a Staff Role holds (ADR-0017); the former `staff_member_role_permissions` bit matrix is dropped |
