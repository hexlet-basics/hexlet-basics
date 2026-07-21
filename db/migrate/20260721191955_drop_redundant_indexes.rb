# typed: false
# frozen_string_literal: true

class DropRedundantIndexes < ActiveRecord::Migration[8.1]
  # Each single-column index below is fully covered by a composite index whose
  # leading column is the same, so the B-tree already serves these lookups.
  # Dropping them removes write-time overhead without losing any query support
  # (the covering indexes are unique, so no constraint is lost either).
  #
  # Concurrent drop does not block reads/writes; it cannot run in a transaction.
  disable_ddl_transaction!

  def change
    remove_index :ai_chats, column: :user_id,
                 name: "index_ai_chats_on_user_id",
                 algorithm: :concurrently, if_exists: true

    remove_index :ai_models, column: :provider,
                 name: "index_ai_models_on_provider",
                 algorithm: :concurrently, if_exists: true

    remove_index :language_lesson_members, column: :user_id,
                 name: "user_finished_lessons_user_id_index",
                 algorithm: :concurrently, if_exists: true

    remove_index :language_lessons, column: :language_id,
                 name: "language_module_lessons_language_id_index",
                 algorithm: :concurrently, if_exists: true

    remove_index :staff_member_role_permissions, column: :role_id,
                 name: "index_staff_member_role_permissions_on_role_id",
                 algorithm: :concurrently, if_exists: true

    remove_index :survey_answers, column: :survey_id,
                 name: "index_survey_answers_on_survey_id",
                 algorithm: :concurrently, if_exists: true

    remove_index :survey_scenario_items, column: :survey_id,
                 name: "index_survey_scenario_items_on_survey_id",
                 algorithm: :concurrently, if_exists: true
  end
end
