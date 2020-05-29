class AddUploadToLanguagesModulesLessons < ActiveRecord::Migration[6.0]
  def change
    add_reference :language_module_lessons, :upload, index: true, foreign_key: true, null: false
    add_reference :language_modules, :upload, index: true, foreign_key: true, null: false
    add_reference :languages, :upload, index: true, foreign_key: true, null: false
  end
end
