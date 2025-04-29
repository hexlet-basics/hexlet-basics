class PopulateUserInMessages < ActiveRecord::Migration[8.0]
  def change
    Language::Lesson::Member::Message.find_each do |m|
      m.user_id = m.language_lesson_member.user_id
      m.save!
    end
  end
end
