class ChangeForeignKeysToBigint < ActiveRecord::Migration[8.0]
  def up
    # language_landing_page_qna_items
    remove_foreign_key :language_landing_page_qna_items, :language_landing_pages
    change_column :language_landing_page_qna_items, :language_landing_page_id, :bigint
    add_foreign_key :language_landing_page_qna_items, :language_landing_pages

    # language_landing_pages - language_id
    remove_foreign_key :language_landing_pages, :languages
    change_column :language_landing_pages, :language_id, :bigint
    add_foreign_key :language_landing_pages, :languages

    # language_landing_pages - language_category_id
    remove_foreign_key :language_landing_pages, :language_categories
    change_column :language_landing_pages, :language_category_id, :bigint
    add_foreign_key :language_landing_pages, :language_categories

    # language_lesson_member_messages - language_id
    remove_foreign_key :language_lesson_member_messages, column: :language_id
    change_column :language_lesson_member_messages, :language_id, :bigint
    add_foreign_key :language_lesson_member_messages, :languages

    # language_lesson_member_messages - language_lesson_id
    remove_foreign_key :language_lesson_member_messages, column: :language_lesson_id
    change_column :language_lesson_member_messages, :language_lesson_id, :bigint
    add_foreign_key :language_lesson_member_messages, :language_lessons

    # language_lesson_member_messages - language_lesson_member_id
    remove_foreign_key :language_lesson_member_messages, column: :language_lesson_member_id
    change_column :language_lesson_member_messages, :language_lesson_member_id, :bigint
    add_foreign_key :language_lesson_member_messages, :language_lesson_members

    # language_lesson_members - language_member_id
    remove_foreign_key :language_lesson_members, column: :language_member_id
    change_column :language_lesson_members, :language_member_id, :bigint
    add_foreign_key :language_lesson_members, :language_members

    # survey_answers - survey_id
    remove_foreign_key :survey_answers, :surveys
    change_column :survey_answers, :survey_id, :bigint
    add_foreign_key :survey_answers, :surveys

    # survey_answers - survey_item_id
    remove_foreign_key :survey_answers, :survey_items
    change_column :survey_answers, :survey_item_id, :bigint
    add_foreign_key :survey_answers, :survey_items

    # survey_answers - user_id
    remove_foreign_key :survey_answers, :users
    change_column :survey_answers, :user_id, :bigint
    add_foreign_key :survey_answers, :users

    # survey_items - survey_id
    remove_foreign_key :survey_items, :surveys
    change_column :survey_items, :survey_id, :bigint
    add_foreign_key :survey_items, :surveys
  end

  def down
    # language_landing_page_qna_items
    remove_foreign_key :language_landing_page_qna_items, :language_landing_pages
    change_column :language_landing_page_qna_items, :language_landing_page_id, :integer
    add_foreign_key :language_landing_page_qna_items, :language_landing_pages

    # language_landing_pages - language_id
    remove_foreign_key :language_landing_pages, :languages
    change_column :language_landing_pages, :language_id, :integer
    add_foreign_key :language_landing_pages, :languages

    # language_landing_pages - language_category_id
    remove_foreign_key :language_landing_pages, :language_categories
    change_column :language_landing_pages, :language_category_id, :integer
    add_foreign_key :language_landing_pages, :language_categories

    # language_lesson_member_messages - language_id
    remove_foreign_key :language_lesson_member_messages, column: :language_id
    change_column :language_lesson_member_messages, :language_id, :integer
    add_foreign_key :language_lesson_member_messages, :languages

    # language_lesson_member_messages - language_lesson_id
    remove_foreign_key :language_lesson_member_messages, column: :language_lesson_id
    change_column :language_lesson_member_messages, :language_lesson_id, :integer
    add_foreign_key :language_lesson_member_messages, :language_lessons

    # language_lesson_member_messages - language_lesson_member_id
    remove_foreign_key :language_lesson_member_messages, column: :language_lesson_member_id
    change_column :language_lesson_member_messages, :language_lesson_member_id, :integer
    add_foreign_key :language_lesson_member_messages, :language_lesson_members

    # language_lesson_members - language_member_id
    remove_foreign_key :language_lesson_members, column: :language_member_id
    change_column :language_lesson_members, :language_member_id, :integer
    add_foreign_key :language_lesson_members, :language_members

    # survey_answers - survey_id
    remove_foreign_key :survey_answers, :surveys
    change_column :survey_answers, :survey_id, :integer
    add_foreign_key :survey_answers, :surveys

    # survey_answers - survey_item_id
    remove_foreign_key :survey_answers, :survey_items
    change_column :survey_answers, :survey_item_id, :integer
    add_foreign_key :survey_answers, :survey_items

    # survey_answers - user_id
    remove_foreign_key :survey_answers, :users
    change_column :survey_answers, :user_id, :integer
    add_foreign_key :survey_answers, :users

    # survey_items - survey_id
    remove_foreign_key :survey_items, :surveys
    change_column :survey_items, :survey_id, :integer
    add_foreign_key :survey_items, :surveys
  end
end
