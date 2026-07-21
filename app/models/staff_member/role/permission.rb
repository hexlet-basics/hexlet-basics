# typed: strict
# frozen_string_literal: true

# == Schema Information
#
# Table name: staff_member_role_permissions
#
#  id          :bigint           not null, primary key
#  can_create  :boolean          default(FALSE), not null
#  can_destroy :boolean          default(FALSE), not null
#  can_index   :boolean          default(FALSE), not null
#  can_update  :boolean          default(FALSE), not null
#  resource    :string           not null
#  created_at  :datetime         not null
#  updated_at  :datetime         not null
#  role_id     :bigint           not null
#
# Indexes
#
#  index_staff_member_role_permissions_on_role_id_and_resource  (role_id,resource) UNIQUE
#
# Foreign Keys
#
#  fk_rails_...  (role_id => staff_member_roles.id)
#
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
