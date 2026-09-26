package localization

import "github.com/nicksnyder/go-i18n/v2/i18n"

// In-lesson assistant prompts, ported from the legacy Assistants::RunJob. They
// are templates over the lesson (TextWith): the model is told what the lesson
// teaches in the learner's language, and is asked to answer in it, so each
// locale carries the whole prompt rather than one Russian prompt with a
// "reply in X" suffix. Lesson content is substituted as data, never parsed as a
// template, so theory that happens to contain `{{` stays literal.
var (
	// AssistantInstructions is the system prompt: the tutoring rules plus the
	// lesson's theory and exercise. Data: Course, Lesson, Theory, Instructions.
	AssistantInstructions = Message{value: i18n.Message{
		ID: "assistant.instructions",
		Other: `You help people learn {{.Course}} as part of a course on Hexlet.
This conversation is about the lesson "{{.Lesson}}".
The course consists of theory, practice and tests that run right in the browser.
You never show the solution to the exercise: the learner must solve it on their own.
Guide them, explain, help them understand, suggest steps towards a solution, put forward hypotheses.
Answer in English. Keep your answers short.

If the discussion reveals a mistake in the theory or the exercise (its description, solution or tests),
recommend writing about it in the community at https://t.me/hexletcommunity, where the project team is.

Lesson theory:
{{.Theory}}

Lesson exercise:
{{.Instructions}}`,
	}}

	// AssistantQuestion wraps one learner question with the editor context the
	// model needs to answer it. Data: Code, Output, Question.
	AssistantQuestion = Message{value: i18n.Message{
		ID: "assistant.question",
		Other: `The learner's code (their solution to the exercise): {{.Code}}

The result of running the learner's code against the exercise tests: {{.Output}}

The learner's question: {{.Question}}`,
	}}
)
