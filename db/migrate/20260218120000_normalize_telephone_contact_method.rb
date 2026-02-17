class NormalizeTelephoneContactMethod < ActiveRecord::Migration[8.0]
  def up
    execute <<~SQL.squish
      UPDATE users
      SET contact_method = 'phone'
      WHERE contact_method = 'telephone'
    SQL
  end

  def down
    raise ActiveRecord::IrreversibleMigration
  end
end
