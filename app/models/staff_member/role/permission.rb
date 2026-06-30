# typed: strict
# frozen_string_literal: true

class StaffMember::Role::Permission < ApplicationRecord
  self.table_name = "staff_member_role_permissions"

  class Resource < T::Enum
    enums do
      BlogPosts             = new("blog_posts")
      Banners               = new("banners")
      Reviews               = new("reviews")
      Leads                 = new("leads")
      Messages              = new("messages")
      LanguageCategories    = new("language_categories")
      LanguageLessons       = new("language_lessons")
      LanguageLessonReviews = new("language_lesson_reviews")
      LanguageLessonMembers = new("language_lesson_members")
      Languages             = new("languages")
      LanguageLandingPages  = new("language_landing_pages")
    end
  end

  belongs_to :role, class_name: "StaffMember::Role"

  typed_enum :resource, Resource

  validates :resource, presence: true
  validates :resource, uniqueness: { scope: :role_id }
end
