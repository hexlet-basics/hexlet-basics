class AddCoverToLanguages < ActiveRecord::Migration[8.0]
  def change
    Language.find_each do |l|
      filepath = Rails.root.join("public/images/#{l.slug}.png")
      file = File.open(filepath)
      l.cover = file
      l.save!
    end
  end
end
