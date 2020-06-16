class AddReferencesToUpload < ActiveRecord::Migration[6.0]
  def change
    add_reference :language_versions, :language_upload, index: true, foreign_key: true
    add_reference :language_module_versions, :language_upload, index: true, foreign_key: true
    add_reference :language_module_lesson_versions, :language_upload, index: true, foreign_key: true
  end
end
