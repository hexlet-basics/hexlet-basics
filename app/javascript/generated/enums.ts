// This file is auto-generated. Do not edit manually.
// Run `bin/rails app:export_enums_to_ts` to update.

/* eslint-disable */
export const enums = {
  adminBlogPostFormState: ["draft","published","archived"],
  adminLanguageFormLearnAs: ["first_language","second_language"],
  adminLanguageFormProgress: ["completed","in_development","draft"],
  adminLanguageLandingPageFormState: ["draft","archived","published"],
  adminReviewFormLocale: ["ru","en"],
  adminReviewFormState: ["draft","published","archived"],
  adminSurveyScenarioFormState: ["active","archived"],
  adminUserFormContactMethod: ["telegram","phone","whatsapp"],
  blogPostState: ["draft","published","archived"],
  bookRequestState: ["requested","downloaded"],
  languageLearnAs: ["first_language","second_language"],
  languageProgress: ["completed","in_development","draft"],
  languageLandingPageState: ["draft","archived","published"],
  languageLessonState: ["created","active","archived"],
  languageLessonMemberMessageRole: ["user","assistant"],
  languageMemberState: ["started","finished"],
  reviewLocale: ["ru","en"],
  reviewState: ["draft","published","archived"],
  surveyItemState: ["active","archived"],
  surveyScenarioState: ["active","archived"],
  surveyScenarioMemberEventName: ["BookRequestedEvent","CourseFinishedEvent","CourseStartedEvent","EmailConfirmedEvent","LeadCreatedEvent","LessonFinishedEvent","LessonStartedEvent","SolutionCheckedEvent","SurveyAnsweredEvent","UserSignedInEvent","UserSignedUpEvent"],
  surveyScenarioMemberState: ["started","finished"],
  surveyScenarioTriggerEventName: ["BookRequestedEvent","CourseFinishedEvent","CourseStartedEvent","EmailConfirmedEvent","LeadCreatedEvent","LessonFinishedEvent","LessonStartedEvent","SolutionCheckedEvent","SurveyAnsweredEvent","UserSignedInEvent","UserSignedUpEvent"],
  userContactMethod: ["telegram","phone","whatsapp"],
  userPasswordFormContactMethod: ["telegram","phone","whatsapp"],
  userProfileFormContactMethod: ["telegram","phone","whatsapp"],
  userSignUpFormContactMethod: ["telegram","phone","whatsapp"],
  userSocialSignupFormContactMethod: ["telegram","phone","whatsapp"],
} as const;
