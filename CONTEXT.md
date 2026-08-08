# Domain Glossary

## Language

### Accounts

**User sign-up**:
The first successful creation of a User account through public registration.
Later sign-ins, profile changes, and administrator-created users are not
sign-ups.

### Learning

**Member**:
Never used unqualified — the word does three unrelated jobs in this codebase.
A **Course Membership** is a learner's enrollment in a Course (`CourseMember`,
table `language_members`); a **Lesson Membership** is their participation in one
Lesson (`CourseLessonMember`, table `language_lesson_members`); a **Staff
Member** is an employee with roles and permissions (`StaffMember`, table
`staff_members`). The first two carry `MemberState`; the third is not about
learning at all. Always say which.

**Course Membership**:
A learner's record of where they are in one Course, and never created twice for
the same pair. It is created by a deliberate action — starting a Lesson, or
submitting a solution to one — and never by loading a page. It grants nothing:
every Course is readable without one, and without an account. A Membership
remembers progress, it does not permit access.
_Also called_: course enrollment
_Avoid_: Language Member, subscription

**Lesson Membership**:
A learner's relationship with one Lesson of a Course they are a member of.
_Also called_: lesson participation
_Avoid_: Language Lesson Member, attempt, submission

**Staff Member**:
An employee, carrying roles and permissions. Unrelated to Course and Lesson
Membership despite the shared word.

**Not started**:
The condition of a learner who has no Course Membership for a Course. It is the
absence of a membership, not a state a membership can be in.
_Avoid_: ready to start

**Started**:
The state of a Course or Lesson Membership from creation until it is finished.

**Finished**:
The state of a Lesson Membership whose exercise the learner has passed, and of a
Course Membership all of whose current Lessons are finished. A membership never
returns to Started.
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
position: it neither raises the gate nor counts towards Progress.

**Available Lesson**:
A Lesson a learner is allowed to take: one whose Position is at most one past
the highest Position they have finished. Availability is measured from the
highest finished Lesson, not from the first unfinished one, so a gap in a
learner's history does not block them. A learner with no history has exactly one
Available Lesson, the first.
_Avoid_: unlocked, permitted

**Progress**:
The share of a Course's Current Lessons a learner has finished, as a percentage.
Lessons finished under a retired Version do not count towards it.

**Next Lesson**:
The first Lesson, in course order, among a Course's Current Lessons that the
learner has not finished — where they resume. Distinct from the Available
Lesson: for a learner with a gap the Next Lesson is the gap, while the Available
Lesson is one past their furthest finished one.

**Guest progress**:
Lessons a visitor passed before having an account. It is one Lesson per Course —
the furthest one they finished — carried in a signed, httpOnly cookie, because
sequential progression makes a guest's gaps unrepresentable and that single
value therefore says everything. It is merged into an account on sign-up and on
sign-in, taking the further of the two positions, and becomes ordinary Course
and Lesson Memberships in the Finished state for the whole prefix up to it.
_Avoid_: anonymous progress, trial progress
