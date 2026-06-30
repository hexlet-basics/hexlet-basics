class CreateStaffMembersAndRoles < ActiveRecord::Migration[8.1]
  def change
    create_table :staff_member_roles do |t|
      t.string :name, null: false
      t.text :description
      t.timestamps
    end

    create_table :staff_member_role_permissions do |t|
      t.references :role, null: false, foreign_key: { to_table: :staff_member_roles }
      t.string :resource, null: false
      t.boolean :can_index, null: false, default: false
      t.boolean :can_create, null: false, default: false
      t.boolean :can_update, null: false, default: false
      t.boolean :can_destroy, null: false, default: false
      t.timestamps
    end
    add_index :staff_member_role_permissions, %i[role_id resource], unique: true

    create_table :staff_members do |t|
      t.references :user, null: false, foreign_key: true, index: { unique: true }
      t.references :role, null: false, foreign_key: { to_table: :staff_member_roles }
      t.string :allowed_locales, array: true, null: false, default: [ "ru" ]
      t.timestamps
    end
  end
end
