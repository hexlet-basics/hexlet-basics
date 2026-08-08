# Domain Glossary

## Language

### Accounts

**User sign-up**:
The first successful creation of a User account through public registration.
Later sign-ins, profile changes, and administrator-created users are not
sign-ups.

### Learning

**Member**:
Reserved for **Staff Member**, an employee carrying roles and permissions
(`StaffMember`, table `staff_members`). It used to do two more jobs, which are
now **Enrollment** and **Lesson Progress**. The tables still say `member`
(`language_members`, `language_lesson_members`) because storage names were kept
across the rename; nothing else should.

**Enrollment**:
A learner's record of where they are in one Course, and never created twice for
the same pair. It is created by a deliberate action — starting a Lesson, or
submitting a solution to one — and never by loading a page. It grants nothing:
every Course is readable without one, and without an account. An Enrollment
remembers progress, it does not permit access.
(`Enrollment`, table `language_members`)
_Avoid_: Course Membership, Language Member, subscription

**Lesson Progress**:
A learner's record of how far they have got with one Lesson of a Course they are
enrolled in. One row per learner per Lesson, carrying the state and nothing the
Lesson itself already says.
(`LessonProgress`, table `language_lesson_members`)
_Avoid_: Lesson Membership, Language Lesson Member, attempt, submission

**Staff Member**:
An employee, carrying roles and permissions. Nothing to do with learning; it
only shares a table-name prefix with Enrollment and Lesson Progress.

**Not started**:
The condition of a learner who has no Enrollment for a Course. It is the absence
of a record, not a state the record can be in.
_Avoid_: ready to start

**Started**:
The state of an Enrollment or Lesson Progress from creation until it is
finished. (`EnrollmentState`, carried by both.)

**Finished**:
The state of a Lesson Progress whose exercise the learner has passed, and of an
Enrollment all of whose Current Lessons are finished. Neither ever returns to
Started. Merging Guest progress is the one case where a Lesson Progress is
created Finished without the learner having passed that Lesson individually:
sequential progression means the whole prefix up to their furthest Lesson was
passed, so the rule implies the rest.
_Avoid_: completed, done, passed

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
_Avoid_: unlocked, permitted

**Completion**:
The share of a Course's Current Lessons a learner has finished, as a percentage.
Lessons finished under a retired Version do not count towards it. Carried as
`Enrollment.progress`; the word Progress on its own is ambiguous now that Lesson
Progress is a record, so prefer Completion for the number.

**Course Readiness**:
How finished the Course itself is — draft, in development, or completed
(`CourseReadiness`, column `progress`). A property of the product, not of any
learner. It was called progress until the rename, which is why the column still
is.
_Avoid_: course progress, course state

**Next Lesson**:
The first Lesson, in course order, among a Course's Current Lessons that the
learner has not finished — where they resume. It is not the gate: for a learner
with a gap, the Next Lesson is the gap, while their Available Lessons run all
the way to one past their furthest finished one.

**Guest progress**:
Lessons a visitor passed before having an account. It is one Lesson per Course —
the furthest one they finished — carried in a signed, httpOnly cookie, because
sequential progression makes a guest's gaps unrepresentable and that single
value therefore says everything. It is merged into an account on sign-up and on
sign-in, taking the further of the two positions, and becomes an ordinary
Enrollment plus Lesson Progress rows in the Finished state for the whole prefix
up to it.
_Avoid_: anonymous progress, trial progress
